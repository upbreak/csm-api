package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func OrganiztionRoute(timeSheetDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	organizationHandler := &handler.HandlerOrganization{
		Service: &service.ServiceOrganization{
			TimeSheetDB: timeSheetDB,
			Store:       r,
		},
	}

	router.Get("/{jno}", organizationHandler.ByProjectList) // 선택한 프로젝트의 조직도 조회

	return router
}
