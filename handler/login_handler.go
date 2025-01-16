package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
)

type LoginHandler struct {
	Service service.GetUserValidService
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
				Result:  "failed",
				Message: err.Error(),
			},
			http.StatusInternalServerError)
		return
	}

	// 유저 유효성 검사
	userId, err := l.Service.GetUserValid(ctx, login.UserId, login.UserPwd)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:  "failed",
				Message: err.Error(),
			},
			http.StatusInternalServerError)
		return
	}

	// 응답 json 구조 정의
	rsp := struct {
		Result string `json:"result"`
		Data   any    `json:"data"`
	}{Result: "success", Data: struct {
		UserId string `json:"user_id"`
	}{UserId: userId}}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
