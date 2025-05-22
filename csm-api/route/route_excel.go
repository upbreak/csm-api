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
		Service: &service.ServiceExcel{},
		FileService: &service.ServiceUploadFile{
			DB:    safeDB,
			TDB:   safeDB,
			Store: r,
		},
	}

	router.Post("/export/daily-deduction", excelHandler.ExportDailyDeduction) // 일별퇴직공제 export
	router.Post("/import/deduction", excelHandler.ImportDeduction)            // 퇴직공제 import
	router.Post("/import", excelHandler.ImportExcel)                          // excel import
	router.Post("/export", excelHandler.ExportExcel)                          // excel export

	return router
}
