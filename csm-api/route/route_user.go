package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func UserRoute(safeDB *sqlx.DB, timeSheetDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	userHandler := &handler.HandlerUser{
		Service: &service.ServiceUser{
			SafeDB:      safeDB,
			TimeSheetDB: timeSheetDB,
			Store:       r,
		},
	}

	router.Get("/role", userHandler.UserRole) // 사용자 권한 조회(프로젝트 선택시 사용)

	return router
}
