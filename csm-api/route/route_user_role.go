package route

import (
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func UserRoleRoute(safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	userRoleHandler := &handler.HandlerUserRole{
		Service: &service.ServiceUserRole{
			SafeDB:  safeDB,
			SafeTDB: safeDB,
			Store:   r,
		},
	}

	router.Get("/uno", userRoleHandler.GetUserRoleListByUno)     // 사용자 권한 조회
	router.Post("/add", userRoleHandler.AddUserRole)             // 사용자 권한 추가
	router.Post("/remove", userRoleHandler.RemoveUserRole)       // 사용자 권한 삭제
	router.Get("/menu-valid", userRoleHandler.UserMenuRoleCheck) // 사용자 메뉴 접근 권한 체크

	return router
}
