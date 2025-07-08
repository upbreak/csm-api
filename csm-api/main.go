package main

import (
	"context"
	"csm-api/config"
	"csm-api/store"
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
	if err := run(context.Background()); err != nil {
		fmt.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error {
	// 시스템 종료 신호 받을 수 있게 context 세팅
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// config 설정
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config.NewConfig: %w", err)
	}

	// domain, port 설정
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Domain, cfg.Port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	fmt.Printf("Listening at %s\n", url)

	// DB config 설정
	dbCfg, err := config.NewDBConfig()
	if err != nil {
		return fmt.Errorf("config.NewDBConfig: %w", err)
	}

	// DB connect
	var cleanup []func()
	safeDb, safeCleanup, err := store.New(ctx, dbCfg.Safe)
	if err != nil {
		return fmt.Errorf("store.New (safeDb): %w", err)
	}
	cleanup = append(cleanup, func() { safeCleanup() })

	timesheetDb, timesheetCleanup, err := store.New(ctx, dbCfg.TimeSheet)
	if err != nil {
		return fmt.Errorf("store.New (timesheetDb): %w", err)
	}
	cleanup = append(cleanup, func() { timesheetCleanup() })

	// api config 생성
	apiCfg, err := config.GetApiConfig()
	if err != nil {
		return fmt.Errorf("config.ApiConfig: %w", err)
	}

	defer func() {
		for _, clean := range cleanup {
			clean()
		}
	}()

	// 초기화 (Init 객체 생성)
	init, err := NewInit(safeDb)
	if err != nil {
		return fmt.Errorf("NewInit fail: %w", err)
	}
	// 초기화 실행
	err = init.RunInitializations(ctx)
	if err != nil {
		return fmt.Errorf("RunInitializations fail: %w", err)
	}

	// 라우팅 설정
	mux, err := newMux(ctx, safeDb, timesheetDb)
	if err != nil {
		return fmt.Errorf("newMux: %w", err)
	}

	// HTTP server 생성
	server := NewServer(l, mux)

	// scheduler 생성
	scheduler, err := NewScheduler(safeDb, apiCfg)
	if err != nil {
		return fmt.Errorf("NewScheduler fail: %w", err)
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
			return fmt.Errorf("server graceful shutdown failed: %w", err)
		}
		log.Println("Server exited normally.")
	}

	return nil
}
