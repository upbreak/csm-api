package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func WorkerRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	workerHandler := handler.HandlerWorker{
		Service: &service.ServiceWorker{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
	}

	router.Get("/total", workerHandler.TotalList)                                   // 전체근로자 조회
	router.Get("/total/simple", workerHandler.AbsentList)                           // 근로자 검색(현장근로자 추가시 사용)
	router.Get("/total/depart", workerHandler.DepartList)                           // 프로젝트 회사명 조회
	router.Post("/total", workerHandler.Add)                                        // 추가
	router.Put("/total", workerHandler.Modify)                                      // 수정
	router.Get("/site-base", workerHandler.SiteBaseList)                            // 현장근로자 조회
	router.Post("/site-base", workerHandler.Merge)                                  // 현장근로자 추가&수정
	router.Post("/site-base/deadline", workerHandler.ModifyDeadline)                // 현장근로자 마감처리
	router.Post("/site-base/project", workerHandler.ModifyProject)                  // 현장근로자 프로젝트 이동
	router.Post("/site-base/delete", workerHandler.SiteBaseRemove)                  // 현장근로자 삭제
	router.Post("/site-base/deadline-cancel", workerHandler.SiteBaseDeadlineCancel) // 마감 취소

	return router
}
