package handler

import (
	"csm-api/auth"
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct: 현장 관리 리스트 핸들러 구조체
type SiteListHandler struct {
	Service     service.SiteService
	CodeService service.CodeService
	Jwt         *auth.JWTUtils
}

// func: 현장 관리 리스트 핸들러 함수
// @param
// - response: targetDate(현재날짜), pCode(부모코드) - url parameter
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

	pCode := r.URL.Query().Get("pCode")
	if pCode == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "pCode parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
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
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	// 현장 관리 코드 조회
	codes, err := s.CodeService.GetCodeList(ctx, pCode)
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
		Values: entity.SiteRes{
			Site: *sites,
			Code: *codes,
		},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)

}

// struct: 현장 관리 수정 핸들러 구조체
type SiteModifyHandler struct {
	Service service.SiteService
}

func (s *SiteModifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	site := entity.Site{}
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
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

	if err := s.Service.ModifySite(ctx, site); err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        DataModifyFailed,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	RespondJSON(ctx, w, &Response{Result: Success}, http.StatusOK)
}

// struct: 현장 데이터 조회 핸들러 구조체
type SiteNmListHandler struct {
	Service service.SiteService
}

// func: 현장 데이터 조회 핸들러 함수
// @param
// -
func (s *SiteNmListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := s.Service.GetSiteNmList(ctx)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			List entity.Sites `json:"list"`
		}{List: *list},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct | func: 현장 상태 조회
type SiteStatsHandler struct {
	Service service.SiteService
}

func (s *SiteStatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	list, err := s.Service.GetSiteStatsList(ctx, targetDate)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			List entity.Sites `json:"list"`
		}{List: *list},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// 지도 x, y좌표 조회
type SiteRoadAddressHandler struct {
	Service *service.ServiceAddressSearching
}

func (s *SiteRoadAddressHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roadAddress := r.URL.Query().Get("roadAddress")

	if roadAddress == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "roadAddress parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}

	mapPoint, err := s.Service.GetAPISiteMapPoint(roadAddress)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			MapPoint entity.MapPoint `json:"point"`
		}{MapPoint: *mapPoint},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
