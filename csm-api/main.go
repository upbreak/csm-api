package main

import (
	"context"
	"csm-api/config"
	"fmt"
	"net"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error {
	// port 환경변수
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config.NewConfig: %w", err)
	}

	// port 설정
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	fmt.Printf("Listening at %s\n", url)

	// db 환경변수
	dbCfg, err := config.NewDBConfig()
	if err != nil {
		return fmt.Errorf("config.NewDBConfig: %w", err)
	}

	// 라우팅 설정
	mux, cleanup, err := newMux(ctx, dbCfg)
	defer func() {
		for _, clean := range cleanup {
			clean()
		}
	}()
	if err != nil {
		return fmt.Errorf("newMux: %w", err)
	}

	// http 서버 생성 및 실행
	server := NewServer(l, mux)
	return server.Run(ctx)
}
