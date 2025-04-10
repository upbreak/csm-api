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
type HandlerCompany struct {
	Service service.CompanyService
}

// func: job 정보 조회
// @param
// - http get paramter
func (h *HandlerCompany) JobInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		BadRequestResponse(ctx, w)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	data, err := h.Service.GetJobInfo(ctx, jnoInt)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		Data entity.JobInfo `json:"data"`
	}{Data: *data}

	SuccessValuesResponse(ctx, w, values)
}

// func: 현장소장 조회
// @param
// - http get paramter
func (h *HandlerCompany) SiteManager(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		BadRequestResponse(ctx, w)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	list, err := h.Service.GetSiteManagerList(ctx, jnoInt)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List entity.Managers `json:"list"`
	}{List: *list}

	SuccessValuesResponse(ctx, w, values)
}

// func: 안전관리자 조회
// @param
// - http get paramter
func (h *HandlerCompany) SafeManager(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		BadRequestResponse(ctx, w)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	list, err := h.Service.GetSafeManagerList(ctx, jnoInt)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List entity.Managers `json:"list"`
	}{List: *list}

	SuccessValuesResponse(ctx, w, values)
}

// func: 관리감독자 조회
// @param
// - http get paramter
func (h *HandlerCompany) Supervisor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		BadRequestResponse(ctx, w)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	list, err := h.Service.GetSupervisorList(ctx, jnoInt)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List entity.Supervisors `json:"list"`
	}{List: *list}

	SuccessValuesResponse(ctx, w, values)
}

// func: 공종 정보 조회
// @param
func (h *HandlerCompany) WorkInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := h.Service.GetWorkInfoList(ctx)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List entity.WorkInfos `json:"list"`
	}{List: *list}

	SuccessValuesResponse(ctx, w, values)
}

// func: 협력업체 정보
// @param
// - http get paramter
func (h *HandlerCompany) CompanyInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	jno := r.URL.Query().Get("jno")
	if jno == "" {
		BadRequestResponse(ctx, w)
		return
	}
	jnoInt, _ := strconv.ParseInt(jno, 10, 64)

	list, err := h.Service.GetCompanyInfoList(ctx, jnoInt)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List entity.CompanyInfoResList `json:"list"`
	}{List: *list}

	SuccessValuesResponse(ctx, w, values)
}
