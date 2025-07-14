package route

import (
	"csm-api/auth"
	"csm-api/handler"
	"github.com/go-chi/chi/v5"
)

func JwtVaildRoute(jwt *auth.JWTUtils) chi.Router {
	router := chi.NewRouter()

	jwtVaildHandler := &handler.JwtValidHandler{
		Jwt: jwt,
	}
	router.Get("/", jwtVaildHandler.ServeHTTP)

	return router
}
