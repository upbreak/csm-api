package handler

import (
	"csm-api/auth"
	"net/http"
)

// api호출시 jwt를 확인하는 미들웨어
func AuthMiddleware(jwt *auth.JWTUtils) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, err := jwt.FillContext(r)
			if err != nil {
				RespondJSON(
					r.Context(),
					w,
					ErrResponse{
						Result:         Failure,
						Message:        err.Error(),
						Details:        InvalidToken,
						HttpStatusCode: http.StatusUnauthorized,
					},
					http.StatusOK,
				)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}
