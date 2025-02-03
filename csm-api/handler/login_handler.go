package handler

import (
	"csm-api/auth"
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
)

type LoginHandler struct {
	Service service.GetUserValidService
	Jwt     *auth.JWTUtils
}

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// body의 raw data를 저장할 구조체 생성
	login := entity.User{}

	// body 데이터 파싱
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        BodyDataParseError,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	// 유저 유효성 검사
	user, err := l.Service.GetUserValid(ctx, login.UserId, login.UserPwd)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         "failed",
				Message:        err.Error(),
				Details:        InvalidUser,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	//jwt 생성
	tokenString, err := l.Jwt.GenerateToken(&auth.JWTClaims{UserId: user.UserId})
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         "jwt created failed",
				Message:        err.Error(),
				Details:        TokenCreatedFail,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	// 쿠키 설정
	cookie := &http.Cookie{
		Name:  "jwt",
		Value: tokenString,
		Path:  "/",
		//Expires:  time.Now().Add(24 * time.Hour), // 1일 동안 유효
		HttpOnly: true,                    // JavaScript로 접근 불가
		Secure:   false,                   // true:HTTPS, false:HTTPS/HTTP
		SameSite: http.SameSiteStrictMode, // 동일 출처에서만 쿠키 전송
	}
	http.SetCookie(w, cookie)

	// 응답 json 구조 정의
	//rsp := struct {
	//	Result string `json:"result"`
	//	Data   any    `json:"data"`
	//}{
	//	Result: "success",
	//	Data: struct {
	//		LoginUserId string `json:"login_user_id"`
	//	}{LoginUserId: login.UserId},
	//}
	rsp := Response{
		Result: Success,
		Values: struct {
			LoginUserId string `json:"login_user_id"`
		}{LoginUserId: login.UserId},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
