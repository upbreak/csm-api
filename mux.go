package main

import (
	"context"
	"csm-api/config"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"net/http"
)

// chi패키지를 이용하여 http method에 따른 여러 요청을 라우팅 할 수 있음 함수 구현
func newMux(ctx context.Context, config *config.DBConfigs) (http.Handler, []func(), error) {
	mux := chi.NewRouter()

	// 테스트용 라우팅
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	// CORS 미들웨어 설정
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://example.com", "http://localhost:3000"}, // 허용할 도메인
		AllowCredentials: true,                                                    // 쿠키 허용
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},     // 허용할 메서드
		AllowedHeaders:   []string{"Content-Type", "Authorization"},               // 허용할 헤더
	})

	// db연결 설정
	var cleanup []func()

	// jwt struct 생성

	// 라우팅 begin
	// 라우팅 end

	handlerMux := c.Handler(mux)

	return handlerMux, cleanup, nil
}
