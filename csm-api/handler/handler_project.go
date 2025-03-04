package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
	"strconv"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-20
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */
type HandlerProjectNm struct {
	Service service.ProjectService
}

// func: 프로젝트 이름 조회
// @param
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
