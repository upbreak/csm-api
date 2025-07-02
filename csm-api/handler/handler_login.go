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

	var (
		user entity.User
		err  error
	)
	if login.IsCompany {
		// 협력업체 유효성 검사
		user, err = l.Service.GetCompanyUserValid(ctx, login.UserId, login.UserPwd)
	} else {
		// 직원 유효성 검사
		user, err = l.Service.GetUserValid(ctx, login.UserId, login.UserPwd)
	}
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        InvalidUser,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	//jwt 생성
	tokenString, err := l.Jwt.GenerateToken(&auth.JWTClaims{Uno: user.Uno, UserId: user.UserId, UserName: user.UserName, IsSaved: login.IsSaved, Role: auth.JWTRole(user.RoleCode)})
	if err != nil {
		RespondJSON(
			ctx,
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

	// 쿠키 설정
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   auth.GetCookieMaxAge(login.IsSaved),
		HttpOnly: true,                    // JavaScript로 접근 불가
		Secure:   false,                   // true:HTTPS, false:HTTPS/HTTP
		SameSite: http.SameSiteStrictMode, // 동일 출처에서만 쿠키 전송
	}
	http.SetCookie(w, cookie)

	rsp := Response{
		Result: Success,
		Values: struct {
			LoginUserId string `json:"login_user_id"`
		}{LoginUserId: login.UserId},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
