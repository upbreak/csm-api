package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
	"strconv"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-20
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct, func: 프로젝트별 근로자 수 조회
type HandlerProjectWorkerCount struct {
	Service service.ProjectService
}

// func:
// @param
//   - get parameter
//     targetDate: 현재시간
func (h *HandlerProjectWorkerCount) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	targetDateString := r.URL.Query().Get("targetDate")
	if targetDateString == "" {
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
	targetDate, err := time.Parse("2006-01-02", targetDateString)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}

	list, err := h.Service.GetProjectWorkerCountList(ctx, targetDate)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			List entity.ProjectInfos `json:"list"`
		}{List: *list},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct, func: 프로젝트 이름 조회
type HandlerProjectNm struct {
	Service service.ProjectService
}

func (h *HandlerProjectNm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := h.Service.GetProjectNmList(ctx)
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
			List entity.ProjectInfos `json:"list"`
		}{List: *list},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct, func: 프로젝트 전체 조회
type HandlerUsedProject struct {
	Service service.ProjectService
}

func (h *HandlerUsedProject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// http get paramter를 저장할 구조체 생성 및 파싱
	page := entity.Page{}
	search := entity.JobInfo{}

	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	order := r.URL.Query().Get("order")
	jobNo := r.URL.Query().Get("job_no")
	compName := r.URL.Query().Get("comp_name")
	orderCompName := r.URL.Query().Get("order_comp_name")
	jobName := r.URL.Query().Get("job_name")
	jobPmName := r.URL.Query().Get("job_pm_name")
	jobSd := r.URL.Query().Get("job_sd")
	jobEd := r.URL.Query().Get("job_ed")
	cdNm := r.URL.Query().Get("cd_nm")
	if pageNum == "" || rowSize == "" {
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
	page.PageNum, _ = strconv.Atoi(pageNum)
	page.RowSize, _ = strconv.Atoi(rowSize)
	page.Order = order
	search.JobNo = jobNo
	search.CompName = compName
	search.OrderCompName = orderCompName
	search.JobName = jobName
	search.JobPmName = jobPmName
	search.JobSd = jobSd
	search.JobEd = jobEd
	search.CdNm = cdNm

	// 프로젝트 조회
	list, err := h.Service.GetUsedProjectList(ctx, page, search)
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

	// 개수
	count, err := h.Service.GetUsedProjectCount(ctx, search)
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
			List  entity.JobInfos `json:"list"`
			Count int             `json:"count"`
		}{List: *list, Count: count},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

type HandlerAllProject struct {
	Service service.ProjectService
}

func (h *HandlerAllProject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page := entity.Page{}
	search := entity.JobInfo{}

	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	order := r.URL.Query().Get("order")
	jobNo := r.URL.Query().Get("job_no")
	compName := r.URL.Query().Get("comp_name")
	orderCompName := r.URL.Query().Get("order_comp_name")
	jobName := r.URL.Query().Get("job_name")
	jobPmName := r.URL.Query().Get("job_pm_name")
	jobSd := r.URL.Query().Get("job_sd")
	jobEd := r.URL.Query().Get("job_ed")
	cdNm := r.URL.Query().Get("cd_nm")

	if pageNum == "" || rowSize == "" {
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

	page.PageNum, _ = strconv.Atoi(pageNum)
	page.RowSize, _ = strconv.Atoi(rowSize)
	page.Order = order
	search.JobNo = jobNo
	search.CompName = compName
	search.OrderCompName = orderCompName
	search.JobName = jobName
	search.JobPmName = jobPmName
	search.JobSd = jobSd
	search.JobEd = jobEd
	search.CdNm = cdNm

	list, err := h.Service.GetAllProjectList(ctx, page, search)

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

	count, err := h.Service.GetAllProjectCount(ctx, search)

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
			List  entity.JobInfos `json:"list"`
			Count int             `json:"count"`
		}{List: *list, Count: count},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)

}

type HandlerStaffProject struct {
	Service service.ProjectService
}

func (h *HandlerStaffProject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uno := r.PathValue("uno")

	int64UNO, err := strconv.ParseInt(uno, 10, 64)
	if err != nil {
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
	}

	page := entity.Page{}
	search := entity.JobInfo{}

	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	order := r.URL.Query().Get("order")
	jobNo := r.URL.Query().Get("job_no")
	compName := r.URL.Query().Get("comp_name")
	orderCompName := r.URL.Query().Get("order_comp_name")
	jobName := r.URL.Query().Get("job_name")
	jobPmName := r.URL.Query().Get("job_pm_name")
	jobSd := r.URL.Query().Get("job_sd")
	jobEd := r.URL.Query().Get("job_ed")
	cdNm := r.URL.Query().Get("cd_nm")

	if pageNum == "" || rowSize == "" {
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
	}

	page.PageNum, _ = strconv.Atoi(pageNum)
	page.RowSize, _ = strconv.Atoi(rowSize)
	page.Order = order
	search.JobNo = jobNo
	search.CompName = compName
	search.OrderCompName = orderCompName
	search.JobName = jobName
	search.JobPmName = jobPmName
	search.JobSd = jobSd
	search.JobEd = jobEd
	search.CdNm = cdNm

	list, err := h.Service.GetStaffProjectList(ctx, page, search, int64UNO)

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
	}

	count, err := h.Service.GetStaffProjectCount(ctx, search, int64UNO)

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
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			List  entity.JobInfos `json:"list"`
			Count int             `json:"count"`
		}{List: *list, Count: count},
	}
	RespondJSON(ctx, w, &rsp, http.StatusOK)

}

type HandlerOrganization struct {
	Service service.ProjectService
}

func (h *HandlerOrganization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jno := r.PathValue("jno")
	int64JNO, err := strconv.ParseInt(jno, 10, 64)

	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "get parameter is missing jno",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
	}

	// 고객사 조직도 조회
	client, err := h.Service.GetClientOrganization(ctx, int64JNO)
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
	}

	// 계약사 조직도 조회
	hitech, err := h.Service.GetHitechOrganization(ctx, int64JNO)
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
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			Client entity.OrganizationPartition  `json:"client"`
			Hitech entity.OrganizationPartitions `json:"hitech"`
		}{Client: *client, Hitech: *hitech},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)

}

type HandlerProjectNmUno struct {
	Service service.ProjectService
}

func (h *HandlerProjectNmUno) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uno := r.PathValue("uno")
	role := r.URL.Query().Get("role")

	int64UNO, err := strconv.ParseInt(uno, 10, 64)

	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "get parameter is missing uno",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
	}

	projectNm, err := h.Service.GetProjectNmUnoList(ctx, int64UNO, role)
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
			ProjectNm entity.ProjectInfos `json:"project_nm"`
		}{ProjectNm: *projectNm},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
