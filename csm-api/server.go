package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"sync"
)

// http 서버를 실행하기 위한 struct
type Server struct {
	srv          *http.Server
	l            net.Listener
	shutdownOnce sync.Once
}

// Server struct에 담긴 값들을 바탕으로 실제 http서버를 고루틴을 이용하여 실행
func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		log.Println("HTTP server starting...")
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return fmt.Errorf("server error: %w", err)
		}
		log.Println("HTTP server closed normally.")
		return nil
	})

	return eg.Wait()
}

// Server struct 생성
func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

// 서버 종료
func (s *Server) Shutdown(ctx context.Context) error {
	var shutdownErr error
	s.shutdownOnce.Do(func() {
		log.Println("HTTP server shutdown initiated...")
		if err := s.srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown error: %+v", err)
			shutdownErr = err
		} else {
			log.Println("HTTP server shutdown completed gracefully.")
		}
	})
	return shutdownErr
}
