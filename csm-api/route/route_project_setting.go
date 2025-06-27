package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func ProjectSettingRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	projectSettingHandler := &handler.HandlerProjectSetting{
		Service: &service.ServiceProjectSetting{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
	}

	router.Get("/{jno}", projectSettingHandler.ProjectSettingList)        // 프로젝트 기본 설정 정보 조회
	router.Post("/", projectSettingHandler.MergeProjectSetting)           // 프로젝트 기본 정보 추가 및 수정
	router.Get("/man-hours/{jno}", projectSettingHandler.ManHourList)     // 프로젝트 공수 정보 조회
	router.Post("/man-hours", projectSettingHandler.MergeManHours)        // 프로젝트 공수 정보 추가 및 수정
	router.Post("/man-hours/{mhno}", projectSettingHandler.DeleteManHour) // 프로젝트 공수정보 삭제
	return router
}
