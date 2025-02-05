package handler

import (
	"csm-api/auth"
	"csm-api/service"
	"net/http"
	"time"
)

type SiteListHandler struct {
	Service service.SiteService
	Jwt     *auth.JWTUtils
}

func (s *SiteListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// GET 요청에서 파라미터 값 읽기
	targetDateString := r.URL.Query().Get("targetDate")
	if targetDateString == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "targetDate parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}
	targetDate, err := time.Parse("2006-01-02", targetDateString)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "Error parsing targetDate",
				Details:        ParsingError,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	// 현장 관리 리스트 조회
	sites, err := s.Service.GetSiteList(ctx, targetDate)
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

	rsp := Response{
		Result: Success,
		Values: sites,
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)

}
