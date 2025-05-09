package route

import (
	"csm-api/handler"
	"csm-api/service"
	"github.com/go-chi/chi/v5"
)

func ExcelRoute() chi.Router {
	router := chi.NewRouter()

	excelHandler := &handler.HandlerExcel{
		Service: &service.ServiceExcel{},
	}

	router.Post("/export/daily-deduction", excelHandler.ExportDailyDeduction) // 일별퇴직공제 export
	router.Post("/import/deduction", excelHandler.ImportDeduction)            // 퇴직공제 import

	return router
}
