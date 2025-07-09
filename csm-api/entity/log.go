package entity

import (
	"csm-api/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ItemLogEntry struct {
	Time     string                 `json:"time"`
	Type     string                 `json:"type"`
	Menu     string                 `json:"menu"`
	UserName string                 `json:"user_name"`
	UserUno  int64                  `json:"user_uno"`
	Record   map[string]interface{} `json:"record"`
	Item     map[string]interface{} `json:"item"`
}

func DecodeItem[T any](r *http.Request, model T) (*ItemLogEntry, T, error) {
	var itemLog ItemLogEntry
	var result T

	// 1. 전체 요청 파싱
	if err := json.NewDecoder(r.Body).Decode(&itemLog); err != nil {
		return nil, result, err
	}

	// 2. itemLog.Item → JSON
	b, err := json.Marshal(itemLog.Item)
	if err != nil {
		return &itemLog, result, err
	}

	// 3. JSON → 원하는 타입(T)
	if err := json.Unmarshal(b, &result); err != nil {
		return &itemLog, result, err
	}

	return &itemLog, result, nil
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
	defer f.Close()

	if _, err := f.WriteString(string(logJSON) + "\n"); err != nil {
		log.Printf("로그 쓰기 실패 (%s): %v\n", logFilePath, err)
		return
	}
}
