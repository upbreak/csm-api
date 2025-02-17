package handler

import "net/http"

type LogoutHandler struct{}

func (l *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 쿠키 삭제
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,                      // 즉시 만료
		HttpOnly: true,                    // JavaScript 접근 방지
		Secure:   false,                   // HTTPS/HTTP 설정 (필요시 true로 변경)
		SameSite: http.SameSiteStrictMode, // 동일 출처에서만 쿠키 전송
	}

	http.SetCookie(w, cookie)

	// 로그아웃 성공 응답
	RespondJSON(ctx, w, &Response{
		Result: Success,
	}, http.StatusOK)

}
