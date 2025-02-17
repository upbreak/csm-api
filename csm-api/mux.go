package main

import (
	"context"
	"csm-api/auth"
	"csm-api/clock"
	"csm-api/config"
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// func:
// @param
// - cfg *config.DBConfigs: db 접속 정보
// chi패키지를 이용하여 http method에 따른 여러 요청을 라우팅 할 수 있음 함수 구현
func newMux(ctx context.Context, cfg *config.DBConfigs) (http.Handler, []func(), error) {
	mux := chi.NewRouter()

	// CORS 미들웨어 설정
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3002", "http://10.10.103.241"}, // 허용할 도메인
		AllowCredentials: true,                                                      // 쿠키 허용
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},       // 허용할 메서드
		AllowedHeaders:   []string{"Content-Type", "Authorization"},                 // 허용할 헤더
	})

	// db연결 설정
	var cleanup []func()
	safeDb, safeCleanup, err := store.New(ctx, cfg.Safe)
	cleanup = append(cleanup, func() { safeCleanup() })
	if err != nil {
		return nil, cleanup, err
	}

	r := store.Repository{Clocker: clock.RealClock{}}

	// jwt struct 생성
	jwt, err := auth.JwtNew(clock.RealClock{})
	if err != nil {
		return nil, cleanup, err
	}

	// 라우팅:: begin
	// 로그인
	loginHandler := &handler.LoginHandler{
		Service: &service.UserValid{
			DB:    safeDb,
			Store: &r,
		},
		Jwt: jwt,
	}
	mux.Post("/login", loginHandler.ServeHTTP)

	// jwt 유효성 검사
	jwtVaildHandler := &handler.JwtValidHandler{
		Jwt: jwt,
	}
	mux.Get("/jwt-validation", jwtVaildHandler.ServeHTTP)

	// 현장관리 리스트 조회
	siteListHandler := &handler.SiteListHandler{
		Service: &service.ServiceSite{
			DB:    safeDb,
			Store: &r,
			ProjectService: &service.ServiceProject{
				DB:    safeDb,
				Store: &r,
				UserService: &service.ServiceUser{
					DB:    safeDb,
					Store: &r,
				},
			},
			ProjectDailyService: &service.ServiceProjectDaily{
				DB:    safeDb,
				Store: &r,
			},
			SitePosService: &service.ServiceSitePos{
				DB:    safeDb,
				Store: &r,
			},
			SiteDateService: &service.ServiceSiteDate{
				DB:    safeDb,
				Store: &r,
			},
		},
		CodeService: &service.ServiceCode{
			DB:    safeDb,
			Store: &r,
		},
		Jwt: jwt,
	}
	mux.Get("/site", siteListHandler.ServeHTTP)

	//현장 데이터 조회
	siteNmListHandler := &handler.SiteNmListHandler{
		Service: &service.ServiceSite{
			DB:    safeDb,
			Store: &r,
		},
	}
	mux.Get("/site-nm", siteNmListHandler.ServeHTTP)

	// 근태인식기
	// 조회
	deviceListHandler := &handler.DeviceListHandler{
		Service: &service.ServiceDevice{
			DB:    safeDb,
			Store: &r,
		},
	}
	// 추가
	deviceAddHandler := &handler.DeviceAddHandler{
		Service: &service.ServiceDevice{
			TDB:   safeDb,
			Store: &r,
		},
	}
	// 수정
	deviceModifyHandler := &handler.DeviceModifyHandler{
		Service: &service.ServiceDevice{
			TDB:   safeDb,
			Store: &r,
		},
	}
	// 삭제
	deviceRemoveHandler := &handler.DeviceRemoveHandler{
		Service: &service.ServiceDevice{
			TDB:   safeDb,
			Store: &r,
		},
	}
	mux.Get("/device", deviceListHandler.ServeHTTP)
	mux.Post("/device", deviceAddHandler.ServeHTTP)
	mux.Put("/device", deviceModifyHandler.ServeHTTP)
	mux.Delete("/device/{id}", deviceRemoveHandler.ServeHTTP)

	// 미들웨어를 사용하여 토큰 검사 후 ServeHTTP 실행
	mux.Route("/test", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		// 테스트용 라우팅
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		})
	})

	// 미들웨어 사용하여 토큰 검사 후 ServeHTTP 실행
	mux.Route("/notice", func(router chi.Router) {
		// router.Use(handler.AuthMiddleware(jwt)) // 이거 넣으면 새로고침 할때마다 인증해야함

		// 공지사항 추가
		addNoticeHandler := &handler.NoticeAddHandler{
			Service: &service.ServiceNotice{
				TDB:   safeDb,
				Store: &r,
			},
		}
		router.Post("/", addNoticeHandler.ServeHTTP)

		// 전체 공지사항 조회
		listNoticeHandler := &handler.NoticeListHandler{
			Service: &service.ServiceNotice{
				DB:    safeDb,
				Store: &r,
			},
		}
		router.Get("/", listNoticeHandler.ServeHTTP)

		// 공지사항 수정

		// 공지사항 삭제

	})

	// 라우팅:: end

	handlerMux := c.Handler(mux)

	return handlerMux, cleanup, nil
}
