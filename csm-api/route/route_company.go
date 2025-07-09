package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func CompanyRoute(safeDB *sqlx.DB, timeSheetDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	companyHandler := handler.HandlerCompany{
		Service: &service.ServiceCompany{
			SafeDB:        safeDB,
			TimeSheetDB:   timeSheetDB,
			Store:         r,
			UserRoleStore: r,
		},
	}

	router.Get("/job-info", companyHandler.JobInfo)         // job 정보 조회
	router.Get("/site-manager", companyHandler.SiteManager) // 현장소장 조회
	router.Get("/safe-manager", companyHandler.SafeManager) // 안전관리자 조회
	router.Get("/supervisor", companyHandler.Supervisor)    // 관리감독자 조회
	router.Get("/work-info", companyHandler.WorkInfo)       // 공종 정보 조회
	router.Get("/company-info", companyHandler.CompanyInfo) // 협력업체 정보

	return router
}
