package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func CompareRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	compareHandler := &handler.HandlerCompare{
		Service: &service.ServiceCompare{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
	}

	router.Get("/", compareHandler.List)         // 일일근로자비교 리스트
	router.Put("/", compareHandler.CompareState) // 일일근로자비교 반영

	return router
}
