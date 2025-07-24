package main

import (
	"context"
	"csm-api/auth"
	"csm-api/config"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 로그파일 경로 세팅
	// 운영서버 sudo vi /etc/logrotate.d/csm 에 로그 정책 설정
	logDir := os.Getenv("CONSOLE_LOG_PATH")
	if logDir == "" {
		logDir = "logs/console"
	}
	logFilePath := logDir + "/csm.log"
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "log file open error: %v\n", err)
		log.SetOutput(os.Stderr) // 파일 Writer 없이 stderr만!
	} else {
		defer func(logFile *os.File) {
			err = logFile.Close()
			if err != nil {
				// 로그 남기지 않아도 됨
			}
		}(logFile)
		mw := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(mw)
	}

	// 리커버
	defer func() {
		if r := recover(); r != nil {
			_ = entity.WriteErrorLog(context.Background(), utils.CustomMessageErrorf("panic recovered", fmt.Errorf("%v", r)))
		}
	}()

	// 서버 실행
	ctx := context.Background()
	ctx = auth.SetContext(ctx, auth.UserId{}, "SYSTEM_MAIN")
	ctx = auth.SetContext(ctx, auth.Uno{}, "0")

	if err := run(ctx); err != nil {
		if !entity.IsLoggedError(err) {
			_ = entity.WriteErrorLog(ctx, utils.CustomMessageErrorf("main() run 실패", err))
		}
	}
}

func run(ctx context.Context) error {
	// 시스템 종료 신호 받을 수 있게 context 세팅
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// config 설정
	cfg, err := config.NewConfig()
	if err != nil {
		return utils.CustomMessageErrorf("config.NewConfig", err)
	}

	// DB config 설정
	dbCfg, err := config.NewDBConfig()
	if err != nil {
		return utils.CustomMessageErrorf("config.NewDBConfig", err)
	}

	// DB connect
	var cleanup []func()
	safeDb, safeCleanup, err := store.New(ctx, dbCfg.Safe)
	if err != nil {
		return utils.CustomMessageErrorf("store.New", err)
	}
	cleanup = append(cleanup, func() { safeCleanup() })

	timesheetDb, timesheetCleanup, err := store.New(ctx, dbCfg.TimeSheet)
	if err != nil {
		return utils.CustomMessageErrorf("store.New", err)
	}
	cleanup = append(cleanup, func() { timesheetCleanup() })

	defer func() {
		for _, clean := range cleanup {
			clean()
		}
	}()

	env := os.Getenv("ENV")
	log.Printf("start go:build env:%s", env)
	role := os.Getenv("ROLE")
	log.Printf("start go:build role:%s", role)

	// 초기화 (Init 객체 생성)
	if env == "development" || env == "local" {
		init, err := NewInit(safeDb)
		if err != nil {
			return utils.CustomMessageErrorf("NewInit fail", err)
		}
		// 초기화 실행
		err = init.RunInitializations(ctx)
		if err != nil {
			return utils.CustomMessageErrorf("RunInitializations", err)
		}
	}

	switch role {
	case "web":
		return runWeb(ctx, cfg, safeDb, timesheetDb)
	case "schedule":
		return runSchedule(ctx, safeDb, timesheetDb)
	default:
		return runWeb(ctx, cfg, safeDb, timesheetDb)
	}

}
