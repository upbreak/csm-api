package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func DeadlineRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	deadlineHandler := &handler.HandlerDeadline{
		UploadService: &service.ServiceUploadFile{
			DB:    safeDB,
			TDB:   safeDB,
			Store: r,
		},
	}

	router.Get("/", deadlineHandler.UploadFileList) // 일일마감 엑셀 자료 정보

	return router
}
