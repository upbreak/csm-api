package entity

import (
	"context"
	"csm-api/auth"
	"csm-api/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 추가/수정/삭제 정상 기록 로그
type ItemLogEntry struct {
	Time     string                   `json:"time"`
	Type     string                   `json:"type"`
	Menu     string                   `json:"menu"`
	UserName string                   `json:"user_name"`
	UserUno  int64                    `json:"user_uno"`
	Record   map[string]interface{}   `json:"record"`
	Item     map[string]interface{}   `json:"item"`
	Items    []map[string]interface{} `json:"items"`
}

// 에러 메세지 기록 로그
type ItemErrLogEntry struct {
	Time       string `json:"time"`
	UserId     string `json:"user_id"`
	UserUno    int64  `json:"user_uno"`
	ErrMessage string `json:"err_message"`
}

type LoggedError struct {
	Err error
}

func (e *LoggedError) Error() string {
	return e.Err.Error()
}

func (e *LoggedError) Unwrap() error {
	return e.Err
}

// Helper: 이미 로그 찍은 에러로 감싸기
func MarkAsLogged(err error) error {
	return &LoggedError{Err: err}
}

// Helper: 이 에러가 이미 로깅된 것인지 확인
func IsLoggedError(err error) bool {
	_, ok := err.(*LoggedError)
	return ok
}

// 정상 기록 구조체 파싱
func DecodeItem[T any](r *http.Request, model T) (*ItemLogEntry, T, error) {
	var itemLog ItemLogEntry
	var result T

	if err := json.NewDecoder(r.Body).Decode(&itemLog); err != nil {
		return nil, result, err
	}

	// 1. 우선 items 배열로 처리 시도
	if itemLog.Items != nil {
		b, err := json.Marshal(itemLog.Items)
		if err != nil {
			return &itemLog, result, err
		}
		if err := json.Unmarshal(b, &result); err != nil {
			return &itemLog, result, err
		}
		return &itemLog, result, nil
	}

	// 2. item 단일 객체 처리
	if itemLog.Item != nil {
		b, err := json.Marshal(itemLog.Item)
		if err != nil {
			return &itemLog, result, err
		}
		if err := json.Unmarshal(b, &result); err != nil {
			return &itemLog, result, err
		}
		return &itemLog, result, nil
	}

	return &itemLog, result, fmt.Errorf("no item or items field found")
}

// 성공한 요청만 로그 파일에 기록 (에러는 콘솔 출력)
func WriteLog(itemLog *ItemLogEntry) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Printf("config.NewConfig() 실패: %v\n", err)
		return
	}

	now := time.Now()
	year := now.Format("2006") // ex: "2025"
	month := now.Format("01")  // ex: "07"
	day := now.Format("20060102")

	// 로그 디렉토리: logs/2025/07
	logDir := filepath.Join(cfg.LogPath, year, month)

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Printf("로그 디렉토리 생성 실패 (%s): %v\n", logDir, err)
		return
	}

	// 로그 파일 경로: logs/2025/07/csm_20250708.log
	logFileName := fmt.Sprintf("csm_%s.log", day)
	logFilePath := filepath.Join(logDir, logFileName)

	// 로그 내용 구성 (item 제외)
	logContent := map[string]interface{}{
		"time":      itemLog.Time,
		"type":      itemLog.Type,
		"user_name": itemLog.UserName,
		"user_uno":  itemLog.UserUno,
		"menu":      itemLog.Menu,
		"record":    itemLog.Record,
	}

	logJSON, err := json.MarshalIndent(logContent, "", "\t")
	if err != nil {
		log.Printf("로그 JSON 직렬화 실패: %v\n", err)
		return
	}

	// 로그 파일에 이어쓰기
	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("로그 파일 열기 실패 (%s): %v\n", logFilePath, err)
		return
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Printf("로그 파일 닫기 실패: %v\n", err)
		}
	}(f)

	if _, err := f.WriteString(string(logJSON) + "\n"); err != nil {
		log.Printf("로그 쓰기 실패 (%s): %v\n", logFilePath, err)
		return
	}
}

// 에러메세지 에러로그파일에 저장 후 에러 반환
func WriteErrorLog(ctx context.Context, err error) error {
	cfg, cfgErr := config.NewConfig()
	if cfgErr != nil {
		log.Printf("config.NewConfig() 실패: %v\n", cfgErr)
		return err // config 실패해도 원래 에러는 그대로 리턴
	}
	
	// 현재 시간
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("20060102")

	// 로그 디렉토리 생성: ex) /tmp/data/csm/error/2025/07
	logDir := filepath.Join(cfg.ErrLogPath, year, month)
	if mkErr := os.MkdirAll(logDir, os.ModePerm); mkErr != nil {
		log.Printf("에러 로그 디렉토리 생성 실패 (%s): %v\n", logDir, mkErr)
		return err
	}

	// 파일명: csm_error_YYYYMMDD.log
	logFileName := fmt.Sprintf("csm_error_%s.log", day)
	logFilePath := filepath.Join(logDir, logFileName)

	// context에서 사용자 정보 추출
	userId, ok := auth.GetContext(ctx, auth.UserId{})
	if !ok {
		userId = "unknown"
	}
	unoStr, ok := auth.GetContext(ctx, auth.Uno{})
	if !ok {
		unoStr = "0"
	}
	userUno, parseErr := strconv.ParseInt(unoStr, 10, 64)
	if parseErr != nil {
		userUno = 0
	}

	// 로그 데이터 생성
	item := &ItemErrLogEntry{
		Time:       now.Format("2006-01-02 15:04:05"),
		UserId:     userId,
		UserUno:    userUno,
		ErrMessage: err.Error(),
	}

	// JSON 직렬화
	logJSON, marshalErr := json.MarshalIndent(item, "", "\t")
	if marshalErr != nil {
		log.Printf("에러 로그 직렬화 실패: %v\n", marshalErr)
		return err
	}

	// 파일에 이어쓰기
	f, openErr := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if openErr != nil {
		log.Printf("에러 로그 파일 열기 실패 (%s): %v\n", logFilePath, openErr)
		return err
	}
	defer func(f *os.File) {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("에러 로그 파일 닫기 실패: %v\n", closeErr)
		}
	}(f)

	if _, writeErr := f.WriteString(string(logJSON) + "\n"); writeErr != nil {
		log.Printf("에러 로그 쓰기 실패 (%s): %v\n", logFilePath, writeErr)
		return err
	}

	log.Printf("userId: %s, userUno: %d, err: %v\n", userId, userUno, err)

	return MarkAsLogged(err)
}
