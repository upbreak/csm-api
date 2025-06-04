package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func ProjectRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	projectHandler := &handler.HandlerProject{
		Service: &service.ServiceProject{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
	}

	router.Get("/", projectHandler.RegList)                        // 공사관리시스템 등록 프로젝트 전체 조회
	router.Get("/worker-count", projectHandler.WorkerCountList)    // 프로젝트별 근로자 수 조회
	router.Get("/enterprise", projectHandler.EnterpriseList)       // 프로젝트 전체 조회
	router.Get("/job_name", projectHandler.JobNameList)            // 프로젝트 이름 조회
	router.Get("/my-org/{uno}", projectHandler.MyOrgList)          // 본인이 속한 조직도의 프로젝트 조회
	router.Get("/my-job_name/{uno}", projectHandler.MyJobNameList) // 본인이 속한 프로젝트 이름 목록
	router.Get("/non-reg", projectHandler.NonRegList)              // 현장근태 사용되지 않은 프로젝트
	router.Get("/project-by-site", projectHandler.ProjectBySite)   // 현장별 프로젝트 조회
	router.Post("/", projectHandler.Add)                           // 추가
	router.Put("/default", projectHandler.ModifyDefault)           // 현장 기본 프로젝트 변경
	router.Put("/use", projectHandler.ModifyIsUse)                 // 현장 프로젝트 사용여부 변경
	router.Delete("/{sno}/{jno}", projectHandler.Remove)           // 현장 프로젝트 삭제

	return router
}
