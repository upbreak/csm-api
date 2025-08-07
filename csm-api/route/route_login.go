package route

import (
	"csm-api/auth"
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func LoginRoute(jwt *auth.JWTUtils, safeDB *sqlx.DB, timesheetDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	loginHandler := &handler.LoginHandler{
		Service: &service.UserValid{
			DB:    safeDB,
			Store: r,
			UserService: &service.ServiceUser{
				SafeDB:      safeDB,
				TimeSheetDB: timesheetDB,
				Store:       r,
			},
		},
		Jwt: jwt,
	}

	router.Post("/", loginHandler.ServeHTTP)

	return router
}

func LogoutRoute() chi.Router {
	router := chi.NewRouter()

	logoutHandler := &handler.LogoutHandler{}
	router.Post("/", logoutHandler.ServeHTTP)

	return router
}
