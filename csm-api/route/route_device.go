package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func DeviceRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	deviceHandler := &handler.DeviceHandler{
		Service: &service.ServiceDevice{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
	}

	router.Get("/", deviceHandler.List)                            // 조회
	router.Post("/", deviceHandler.Add)                            // 추가
	router.Put("/", deviceHandler.Modify)                          // 수정
	router.Delete("/{id}", deviceHandler.Remove)                   // 삭제
	router.Get("/check-registered", deviceHandler.CheckRegistered) // 장치 등록 확인

	return router
}
