package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func CodeRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	codeHandler := &handler.HandlerCode{
		Service: service.ServiceCode{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
	}

	router.Get("/", codeHandler.ListByPCode)          // code 조회
	router.Get("/tree", codeHandler.ListCodeTree)     // codeTree 조회
	router.Get("/check", codeHandler.DuplicateByCode) // code 조회
	router.Post("/", codeHandler.Merge)               // 코드 추가 및 수정
	router.Delete("/{idx}", codeHandler.Remove)       // 코드 삭제
	router.Post("/sort", codeHandler.SortNoModify)    // 코드순서 수정

	return router
}
