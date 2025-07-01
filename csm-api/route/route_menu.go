package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func MenuRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	menuHandler := &handler.HandlerMenu{
		Service: &service.ServiceMenu{
			SafeDB: safeDB,
			Store:  r,
		},
	}

	router.Get("/", menuHandler.List) // 권한별 메뉴 리스트

	return router
}
