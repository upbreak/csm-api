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
 * @modified 최종 수정일: 2025-02-21
 * @modifiedBy 최종 수정자: 정지영
 * @modified description
 * - 공지사항 기능 추가
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
	// safe 스키마 연결
	safeDb, safeCleanup, err := store.New(ctx, cfg.Safe)
	cleanup = append(cleanup, func() { safeCleanup() })
	if err != nil {
		return nil, cleanup, err
	}
	// timesheet 스키마 연결
	timesheetDb, timeSheetCleanup, err := store.New(ctx, cfg.TimeSheet)
	cleanup = append(cleanup, func() { timeSheetCleanup() })
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
	// 로그아웃
	logoutHandler := &handler.LogoutHandler{}
	mux.Post("/logout", logoutHandler.ServeHTTP)

	// jwt 유효성 검사
	jwtVaildHandler := &handler.JwtValidHandler{
		Jwt: jwt,
	}
	mux.Get("/jwt-validation", jwtVaildHandler.ServeHTTP)

	// Begin::현장관리
	// 현장관리 조회
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
	// 현장 이름 리스트 조회
	siteNmListHandler := &handler.SiteNmListHandler{
		Service: &service.ServiceSite{
			DB:    safeDb,
			Store: &r,
		},
	}
	// 현장 상태 리스트 조회
	siteStatsHandler := &handler.SiteStatsHandler{
		Service: &service.ServiceSite{
			DB:    safeDb,
			Store: &r,
		},
	}
	mux.Route("/site", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/", siteListHandler.ServeHTTP)
		r.Get("/nm", siteNmListHandler.ServeHTTP)
		r.Get("/stats", siteStatsHandler.ServeHTTP)
	})
	// End::현장관리

	// Begin::프로젝트 조회
	// 프로젝트 이름 데이터 조회
	projectNmHandler := &handler.HandlerProjectNm{
		Service: &service.ServiceProject{
			DB:    safeDb,
			Store: &r,
		},
	}
	// 프로젝트 전체 조회
	usedProjectHandler := &handler.HandlerUsedProject{
		Service: &service.ServiceProject{
			DB:    safeDb,
			Store: &r,
		},
	}

	// 진행중인 프로젝트 전체 조회
	allProjectHandler := &handler.HandlerAllProject{
		Service: &service.ServiceProject{
			DB:    safeDb,
			Store: &r,
		},
	}

	// 조직도 프로젝트 정보
	staffProjectHandler := &handler.HandlerStaffProject{
		Service: &service.ServiceProject{
			DB:    safeDb,
			Store: &r,
		},
	}

	// 조직도 정보
	organizationHandler := &handler.HandlerOrganization{
		Service: &service.ServiceProject{
			DB:    timesheetDb,
			Store: &r,
		},
	}

	mux.Route("/project", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/used", usedProjectHandler.ServeHTTP)
		r.Get("/all", allProjectHandler.ServeHTTP)
		r.Get("/nm", projectNmHandler.ServeHTTP)
		r.Get("/staff/{uno}", staffProjectHandler.ServeHTTP)
		r.Get("/organization/{jno}", organizationHandler.ServeHTTP)
	})
	// End::프로젝트 조회

	// Begin::근태인식기
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
	mux.Route("/device", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/", deviceListHandler.ServeHTTP)
		r.Post("/", deviceAddHandler.ServeHTTP)
		r.Put("/", deviceModifyHandler.ServeHTTP)
		r.Delete("/{id}", deviceRemoveHandler.ServeHTTP)
	})
	// End::근태인식기

	// Begin::근로자
	// 전체 근로자 조회
	workerTotalListHandler := handler.HandlerWorkerTotalList{
		Service: &service.ServiceWorker{
			DB:    safeDb,
			Store: &r,
		},
	}
	// 현장 근로자 조회
	workerSiteBaseListHandler := handler.HandlerWorkerSiteBaseList{
		Service: &service.ServiceWorker{
			DB:    safeDb,
			Store: &r,
		},
	}
	mux.Route("/worker", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/total", workerTotalListHandler.ServeHttp)
		r.Get("/site-base", workerSiteBaseListHandler.ServeHttp)
	})
	// End::근로자

	// Begin::협력업체
	// job(프로젝트) 정보
	companyJobInfoHandler := handler.HandlerJobInfoCompany{
		Service: &service.ServiceCompany{
			SafeDB: safeDb,
			Store:  &r,
		},
	}
	// 현장소장
	companySiteManagerHandler := handler.HandlerSiteManagerCompany{
		Service: &service.ServiceCompany{
			TimeSheetDB: timesheetDb,
			Store:       &r,
		},
	}
	// 안전관리자
	companySafeManagerHandler := handler.HandlerSafeManagerCompany{
		Service: &service.ServiceCompany{
			SafeDB: safeDb,
			Store:  &r,
		},
	}
	// 관리감독자 정보
	companySupervisorHandler := handler.HandlerSupervisorCompany{
		Service: &service.ServiceCompany{
			SafeDB: safeDb,
			Store:  &r,
		},
	}
	// 공종 정보
	companyWorkInfoHandler := handler.HandlerWorkInfoCompany{
		Service: &service.ServiceCompany{
			SafeDB: safeDb,
			Store:  &r,
		},
	}
	// 협력업체 정보
	companyCompanyInfoHandler := handler.HandlerCompanyInfoCompany{
		Service: &service.ServiceCompany{
			SafeDB: safeDb,
			Store:  &r,
		},
	}
	mux.Route("/company", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/job-info", companyJobInfoHandler.ServeHTTP)
		r.Get("/site-manager", companySiteManagerHandler.ServeHTTP)
		r.Get("/safe-manager", companySafeManagerHandler.ServeHTTP)
		r.Get("/supervisor", companySupervisorHandler.ServeHTTP)
		r.Get("/work-info", companyWorkInfoHandler.ServeHTTP)
		r.Get("/company-info", companyCompanyInfoHandler.ServeHTTP)
	})
	// End::협력업체

	// Begin::공지사항
	mux.Route("/notice", func(router chi.Router) {
		router.Use(handler.AuthMiddleware(jwt))

		// 공지사항 추가
		noticeAddHandler := &handler.NoticeAddHandler{
			Service: &service.ServiceNotice{
				TDB:   safeDb,
				Store: &r,
			},
		}

		// 전체 공지사항 조회
		noticeListHandler := &handler.NoticeListHandler{
			Service: &service.ServiceNotice{
				DB:    safeDb,
				Store: &r,
			},
		}

		// 공지사항 수정
		noticeModifyHandler := &handler.NoticeModifyHandler{
			Service: &service.ServiceNotice{
				TDB:   safeDb,
				Store: &r,
			},
		}
		// 공지사항 삭제
		noticeDeleteHandler := &handler.NoticeDeleteHandler{
			Service: &service.ServiceNotice{
				TDB:   safeDb,
				Store: &r,
			},
		}

		router.Post("/", noticeAddHandler.ServeHTTP)
		router.Get("/", noticeListHandler.ServeHTTP)
		router.Put("/", noticeModifyHandler.ServeHTTP)
		router.Delete("/{idx}", noticeDeleteHandler.ServeHTTP)

	})
	// End::공지사항

	// 라우팅:: end

	handlerMux := c.Handler(mux)

	return handlerMux, cleanup, nil
}
