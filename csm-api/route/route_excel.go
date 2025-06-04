package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func ExcelRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	excelHandler := &handler.HandlerExcel{
		Service: &service.ServiceExcel{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
		FileService: &service.ServiceUploadFile{
			DB:    safeDB,
			TDB:   safeDB,
			Store: r,
		},
		DB: safeDB,
	}

	router.Post("/import", excelHandler.ImportExcel) // excel import
	router.Get("/export", excelHandler.ExportExcel)  // excel export

	return router
}
