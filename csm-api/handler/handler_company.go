package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
	"strconv"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-18
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct: job 정보 조회
type HandlerJobInfoCompany struct {
	Service service.CompanyService
}

// func: job 정보 조회
// @param
// - http get paramter
func (h *HandlerJobInfoCompany) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "get parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	data, err := h.Service.GetJobInfo(ctx, jnoInt)
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
			Data entity.JobInfo `json:"data"`
		}{Data: *data},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct: 현장소장
type HandlerSiteManagerCompany struct {
	Service service.CompanyService
}

// func: 현장소장 조회
// @param
// - http get paramter
func (h *HandlerSiteManagerCompany) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "get parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	list, err := h.Service.GetSiteManagerList(ctx, jnoInt)
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
			List entity.Managers `json:"list"`
		}{List: *list},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct: 안전관리자
type HandlerSafeManagerCompany struct {
	Service service.CompanyService
}

// func: 안전관리자 조회
// @param
// - http get paramter
func (h *HandlerSafeManagerCompany) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "get parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	list, err := h.Service.GetSafeManagerList(ctx, jnoInt)
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
			List entity.Managers `json:"list"`
		}{List: *list},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct: 관리감독자
type HandlerSupervisorCompany struct {
	Service service.CompanyService
}

// func: 관리감독자 조회
// @param
// - http get paramter
func (h *HandlerSupervisorCompany) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "get parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	list, err := h.Service.GetSupervisorList(ctx, jnoInt)
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
			List entity.Supervisors `json:"list"`
		}{List: *list},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct: 협력업체 정보
type HandlerCompanyInfoCompany struct {
	Service service.CompanyService
}

// func: 협력업체 정보
// @param
// - http get paramter
func (h *HandlerCompanyInfoCompany) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "get parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	list, err := h.Service.GetCompanyInfoList(ctx, jnoInt)
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
			List entity.CompanyInfoResList `json:"list"`
		}{List: *list},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
