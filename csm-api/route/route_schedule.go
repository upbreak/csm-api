package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func ScheduleRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	// 휴무일
	restScheduleHandler := &handler.HandlerRestSchedule{
		Service: &service.ServiceSchedule{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
	}
	// 작업내용
	dailyJobHandler := &handler.HandlerProjectDaily{
		Service: &service.ServiceProjectDaily{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
	}

	router.Get("/rest", restScheduleHandler.RestList)            // 휴무일 조회
	router.Post("/rest", restScheduleHandler.RestAdd)            // 휴무일 추가
	router.Put("/rest", restScheduleHandler.RestModify)          // 휴무일 수정
	router.Delete("/rest/{cno}", restScheduleHandler.RestRemove) // 휴무일 삭제
	router.Get("/daily-job", dailyJobHandler.List)               // 작업내용 조회
	router.Post("/daily-job", dailyJobHandler.Add)               // 작업내용 추가
	router.Put("/daily-job", dailyJobHandler.Modify)             // 작업내용 수정
	router.Delete("/daily-job/{idx}", dailyJobHandler.Remove)    // 작업내용 삭제

	return router
}
