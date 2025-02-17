package handler

import (
	"csm-api/auth"
	"net/http"
)

// api호출시 jwt를 확인하는 미들웨어
func AuthMiddleware(jwt *auth.JWTUtils) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// jwt 유효성 검사
			req, claims, err := jwt.FillContext(r)
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
					http.StatusUnauthorized,
				)
				return
			}

			// 아이디 저장을 안 했을 경우
			if !claims.IsSaved {
				//jwt 생성
				tokenString, err := jwt.GenerateToken(claims)
				if err != nil {
					RespondJSON(
						r.Context(),
						w,
						&ErrResponse{
							Result:         Failure,
							Message:        err.Error(),
							Details:        TokenCreatedFail,
							HttpStatusCode: http.StatusInternalServerError,
						},
						http.StatusOK)
					return
				}

				// 쿠키 설정(토큰 재발급)
				cookie := &http.Cookie{
					Name:     "jwt",
					Value:    tokenString,
					Path:     "/",
					MaxAge:   auth.GetCookieMaxAge(claims.IsSaved),
					HttpOnly: true,                    // JavaScript로 접근 불가
					Secure:   false,                   // true:HTTPS, false:HTTPS/HTTP
					SameSite: http.SameSiteStrictMode, // 동일 출처에서만 쿠키 전송
				}
				http.SetCookie(w, cookie)
			}

			next.ServeHTTP(w, req)
		})
	}
}
