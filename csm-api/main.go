package main

import (
	"context"
	"csm-api/auth"
	"csm-api/config"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			_ = entity.WriteErrorLog(context.Background(), utils.CustomMessageErrorf("panic recovered", fmt.Errorf("%v", r)))
		}
	}()

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

	// domain, port 설정
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Domain, cfg.Port))
	if err != nil {
		return utils.CustomMessageErrorf("net.Listen", err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	fmt.Printf("Listening at %s\n", url)

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

	// api config 생성
	apiCfg, err := config.GetApiConfig()
	if err != nil {
		return utils.CustomMessageErrorf("config.ApiConfig", err)
	}

	defer func() {
		for _, clean := range cleanup {
			clean()
		}
	}()

	// 초기화 (Init 객체 생성)
	init, err := NewInit(safeDb)
	if err != nil {
		return utils.CustomMessageErrorf("NewInit fail", err)
	}
	// 초기화 실행
	err = init.RunInitializations(ctx)
	if err != nil {
		return utils.CustomMessageErrorf("RunInitializations", err)
	}

	// 라우팅 설정
	mux, err := newMux(ctx, safeDb, timesheetDb)
	if err != nil {
		return utils.CustomMessageErrorf("newMux", err)
	}

	// HTTP server 생성
	server := NewServer(l, mux)

	// scheduler 생성
	scheduler, err := NewScheduler(safeDb, apiCfg, timesheetDb)
	if err != nil {
		return utils.CustomMessageErrorf("NewScheduler", err)
	}

	// 서버와 스케줄러 동시에 실행
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return server.Run(ctx)
	})

	eg.Go(func() error {
		return scheduler.Run(ctx)
	})

	// 종료 신호 대기
	select {
	case <-ctx.Done():
		fmt.Println("\nShutdown signal received")

		// 서버 shutdown (5초 안에 처리)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = server.Shutdown(shutdownCtx); err != nil {
			return utils.CustomMessageErrorf("server graceful shutdown", err)
		}
		log.Println("Server exited normally.")
	}

	return nil
}
