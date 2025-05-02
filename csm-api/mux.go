package main

import (
	"context"
	"csm-api/auth"
	"csm-api/clock"
	"csm-api/config"
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"net/http"
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
func newMux(ctx context.Context, safeDb *sqlx.DB, timesheetDb *sqlx.DB) (http.Handler, error) {
	mux := chi.NewRouter()

	// CORS 미들웨어 설정
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3002", "http://10.10.103.241"}, // 허용할 도메인
		AllowCredentials: true,                                                      // 쿠키 허용
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},       // 허용할 메서드
		AllowedHeaders:   []string{"Content-Type", "Authorization"},                 // 허용할 헤더
	})
	r := store.Repository{Clocker: clock.RealClock{}}

	// jwt struct 생성
	jwt, err := auth.JwtNew(clock.RealClock{})
	if err != nil {
		return nil, err
	}

	// api config 생성
	apiCfg, err := config.GetApiConfig()
	if err != nil {
		return nil, err
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

	// Begin :: 코드관리
	codeHandler := &handler.HandlerCode{
		Service: service.ServiceCode{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
	}

	mux.Route("/code", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/", codeHandler.ListByPCode)          // code 조회
		r.Get("/tree", codeHandler.ListCodeTree)     // codeTree 조회
		r.Get("/check", codeHandler.DuplicateByCode) // code 조회
		r.Post("/", codeHandler.Merge)               // 코드 추가 및 수정
		r.Delete("/{idx}", codeHandler.Remove)       // 코드 삭제
		r.Post("/sort", codeHandler.SortNoModify)    // 코드순서 수정
	})
	// End :: 코드관리

	// Begin::현장관리
	siteHandler := &handler.HandlerSite{
		Service: &service.ServiceSite{
			SafeDB:            safeDb,
			SafeTDB:           safeDb,
			Store:             &r,
			ProjectStore:      &r,
			ProjectDailyStore: &r,
			SitePosStore:      &r,
			SiteDateStore:     &r,
			ProjectService: &service.ServiceProject{
				SafeDB:    safeDb,
				Store:     &r,
				UserStore: &r,
			},
			WhetherApiService: &service.ServiceWhether{
				ApiKey: apiCfg,
			},
			AddressSearchAPIService: &service.ServiceAddressSearch{
				ApiKey: apiCfg,
			},
		},
		CodeService: &service.ServiceCode{
			SafeDB: safeDb,
			Store:  &r,
		},
		Jwt: jwt,
	}
	mux.Route("/site", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/", siteHandler.List)                // 현장관리 조회
		r.Get("/nm", siteHandler.SiteNameList)      // 현장명 조회
		r.Get("/stats", siteHandler.StatsList)      // 현장상태조회
		r.Post("/", siteHandler.Add)                // 현장 생성
		r.Put("/", siteHandler.Modify)              // 수정
		r.Put("/non-use", siteHandler.ModifyNonUse) // 현장 사용안함
	})
	// End::현장관리

	// Begin::지도좌표
	roadAddressHandler := &handler.HandlerRoadAddress{
		Service: &service.ServiceAddressSearch{
			ApiKey: apiCfg,
		},
	}
	mux.Route("/map", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/point", roadAddressHandler.AddressPoint)
	})
	// End::지도좌표

	// Begin:: api 호출
	// 기상청 초단기 실황
	handlerWhetherSrt := &handler.HandlerWhetherSrtNcst{
		Service: &service.ServiceWhether{
			ApiKey: apiCfg,
		},
		SitePosService: &service.ServiceSitePos{
			DB:    safeDb,
			Store: &r,
		},
	}

	// 기상청 기상특보통보문 조회
	handlerWhetherWrn := &handler.HandlerWhetherWrnMsg{
		Service: &service.ServiceWhether{
			ApiKey: apiCfg,
		},
	}

	// 공휴일 조회
	restDatehandler := &handler.HandlerRestDate{
		Service: &service.ServiceRestDate{
			ApiKey: apiCfg,
		},
	}

	mux.Route("/api", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/whether/srt", handlerWhetherSrt.ServeHTTP)
		r.Get("/whether/wrn", handlerWhetherWrn.ServeHTTP)
		r.Get("/rest-date", restDatehandler.ServeHTTP)
	})
	// End:: api 호출

	// Begin::프로젝트
	projectHandler := &handler.HandlerProject{
		Service: &service.ServiceProject{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
	}
	mux.Route("/project", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/", projectHandler.RegList)                        // 공사관리시스템 등록 프로젝트 전체 조회
		r.Get("/worker-count", projectHandler.WorkerCountList)    // 프로젝트별 근로자 수 조회
		r.Get("/enterprise", projectHandler.EnterpriseList)       // 프로젝트 전체 조회
		r.Get("/job_name", projectHandler.JobNameList)            // 프로젝트 이름 조회
		r.Get("/my-org/{uno}", projectHandler.MyOrgList)          // 본인이 속한 조직도의 프로젝트 조회
		r.Get("/my-job_name/{uno}", projectHandler.MyJobNameList) // 본인이 속한 프로젝트 이름 목록
		r.Get("/non-reg", projectHandler.NonRegList)              // 현장근태 사용되지 않은 프로젝트
		r.Post("/", projectHandler.Add)                           // 추가
		r.Put("/default", projectHandler.ModifyDefault)           // 현장 기본 프로젝트 변경
		r.Put("/use", projectHandler.ModifyIsUse)                 // 현장 프로젝트 사용여부 변경
		r.Delete("/{sno}/{jno}", projectHandler.Remove)           // 현장 프로젝트 삭제
	})
	// End::프로젝트 조회

	// Begin::조직도
	organizationHandler := &handler.HandlerOrganization{
		Service: &service.ServiceOrganization{
			SafeDB: timesheetDb,
			Store:  &r,
		},
	}
	mux.Route("/organization", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/{jno}", organizationHandler.ByProjectList) // 선택한 프로젝트의 조직도 조회
	})
	// End::조직도

	// Begin::근태인식기
	deviceHandler := &handler.DeviceHandler{
		Service: &service.ServiceDevice{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
	}
	mux.Route("/device", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/", deviceHandler.List)                            // 조회
		r.Post("/", deviceHandler.Add)                            // 추가
		r.Put("/", deviceHandler.Modify)                          // 수정
		r.Delete("/{id}", deviceHandler.Remove)                   // 삭제
		r.Get("/check-registered", deviceHandler.CheckRegistered) // 장치 등록 확인
	})
	// End::근태인식기

	// Begin::근로자
	workerHandler := handler.HandlerWorker{
		Service: &service.ServiceWorker{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
	}
	mux.Route("/worker", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/total", workerHandler.TotalList)                    // 전체근로자 조회
		r.Get("/total/simple", workerHandler.AbsentList)            // 근로자 검색(현장근로자 추가시 사용)
		r.Post("/total", workerHandler.Add)                         // 추가
		r.Put("/total", workerHandler.Modify)                       // 수정
		r.Get("/site-base", workerHandler.SiteBaseList)             // 현장근로자 조회
		r.Post("/site-base", workerHandler.Merge)                   // 현장근로자 추가&수정
		r.Post("/site-base/deadline", workerHandler.ModifyDeadline) // 현장근로자 마감처리
		r.Post("/site-base/project", workerHandler.ModifyProject)   // 현장근로자 프로젝트 이동
	})
	// End::근로자

	// Begin::협력업체
	companyHandler := handler.HandlerCompany{
		Service: &service.ServiceCompany{
			SafeDB:      safeDb,
			TimeSheetDB: timesheetDb,
			Store:       &r,
		},
	}
	mux.Route("/company", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/job-info", companyHandler.JobInfo)         // job 정보 조회
		r.Get("/site-manager", companyHandler.SiteManager) // 현장소장 조회
		r.Get("/safe-manager", companyHandler.SafeManager) // 안전관리자 조회
		r.Get("/supervisor", companyHandler.Supervisor)    // 관리감독자 조회
		r.Get("/work-info", companyHandler.WorkInfo)       // 공종 정보 조회
		r.Get("/company-info", companyHandler.CompanyInfo) // 협력업체 정보
	})
	// End::협력업체

	// Begin::공지사항
	noticeHandler := &handler.NoticeHandler{
		Service: &service.ServiceNotice{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
	}
	mux.Route("/notice", func(router chi.Router) {
		router.Use(handler.AuthMiddleware(jwt))
		router.Get("/{uno}", noticeHandler.List)      // 조회
		router.Post("/", noticeHandler.Add)           // 추가
		router.Put("/", noticeHandler.Modify)         // 수정
		router.Delete("/{idx}", noticeHandler.Remove) // 삭제
	})
	// End::공지사항

	// Begin::장비
	// 장비 핸들러
	equipHandler := &handler.HandlerEquip{
		Service: &service.ServiceEquip{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
	}
	mux.Route("/equip", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/", equipHandler.List)   // 장비 조회 (임시)
		r.Post("/", equipHandler.Merge) // 장비 입력 (임시)
	})
	// End::장비

	// Begin::일정관리
	// 휴무일
	restScheduleHandler := &handler.HandlerRestSchedule{
		Service: &service.ServiceSchedule{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
	}
	// 작업내용
	dailyJobHandler := &handler.HandlerProjectDaily{
		Service: &service.ServiceProjectDaily{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
	}
	mux.Route("/schedule", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/rest", restScheduleHandler.RestList)            // 휴무일 조회
		r.Post("/rest", restScheduleHandler.RestAdd)            // 휴무일 추가
		r.Put("/rest", restScheduleHandler.RestModify)          // 휴무일 수정
		r.Delete("/rest/{cno}", restScheduleHandler.RestRemove) // 휴무일 삭제
		r.Get("/daily-job", dailyJobHandler.List)               // 작업내용 조회
		r.Post("/daily-job", dailyJobHandler.Add)               // 작업내용 추가
		r.Put("/daily-job", dailyJobHandler.Modify)             // 작업내용 수정
		r.Delete("/daily-job/{cno}", dailyJobHandler.Remove)    // 작업내용 삭제
	})
	// End::일정관리

	// Begin::엑셀
	excelHandler := &handler.HandlerExcel{}
	mux.Route("/excel", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Post("/daily-deduction", excelHandler.DailyDeduction) // 일별퇴직공제
	})
	// End::엑셀

	// 라우팅:: end

	handlerMux := c.Handler(mux)

	return handlerMux, nil
}
