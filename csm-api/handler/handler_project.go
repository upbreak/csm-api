package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
	"encoding/json"
	"github.com/go-chi/chi/v5"
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

// struct, func: 공사관리시스템 등록 프로젝트 전체 조회
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
	search.JobNo = utils.ParseNullString(jobNo)
	search.CompName = utils.ParseNullString(compName)
	search.OrderCompName = utils.ParseNullString(orderCompName)
	search.JobName = utils.ParseNullString(jobName)
	search.JobPmName = utils.ParseNullString(jobPmName)
	search.JobSd = utils.ParseNullString(jobSd)
	search.JobEd = utils.ParseNullString(jobEd)
	search.CdNm = utils.ParseNullString(cdNm)

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

// struct, func: 프로젝트 전체 조회
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
	search.JobNo = utils.ParseNullString(jobNo)
	search.CompName = utils.ParseNullString(compName)
	search.OrderCompName = utils.ParseNullString(orderCompName)
	search.JobName = utils.ParseNullString(jobName)
	search.JobPmName = utils.ParseNullString(jobPmName)
	search.JobSd = utils.ParseNullString(jobSd)
	search.JobEd = utils.ParseNullString(jobEd)
	search.CdNm = utils.ParseNullString(cdNm)

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

// struct, func: 본인이 속한 프로젝트 조회
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
	search.JobNo = utils.ParseNullString(jobNo)
	search.CompName = utils.ParseNullString(compName)
	search.OrderCompName = utils.ParseNullString(orderCompName)
	search.JobName = utils.ParseNullString(jobName)
	search.JobPmName = utils.ParseNullString(jobPmName)
	search.JobSd = utils.ParseNullString(jobSd)
	search.JobEd = utils.ParseNullString(jobEd)
	search.CdNm = utils.ParseNullString(cdNm)

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

// struct, func: 조직도 조회
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
			Client entity.OrganizationPartitions `json:"client"`
			Hitech entity.OrganizationPartitions `json:"hitech"`
		}{Client: *client, Hitech: *hitech},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)

}

// 본인이 속한 프로젝트 이름 목록
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

// struct, func: 현장근태 사용되지 않은 프로젝트
type HandlerNonUsedProject struct {
	Service service.ProjectService
}

func (h *HandlerNonUsedProject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page := entity.Page{}
	search := entity.NonUsedProject{}

	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	order := r.URL.Query().Get("order")
	rnumOrder := r.URL.Query().Get("rnum_order")
	retrySearch := r.URL.Query().Get("retry_search")
	jno := r.URL.Query().Get("jno")
	jobNo := r.URL.Query().Get("job_no")
	JobName := r.URL.Query().Get("job_name")
	JobYear := r.URL.Query().Get("job_year")
	JobSd := r.URL.Query().Get("job_sd")
	JobEd := r.URL.Query().Get("job_ed")
	UserName := r.URL.Query().Get("job_pm_nm")

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
	page.RnumOrder = rnumOrder

	search.Jno = utils.ParseNullInt(jno)
	search.JobNo = utils.ParseNullString(jobNo)
	search.JobName = utils.ParseNullString(JobName)
	search.JobYear = utils.ParseNullInt(JobYear)
	search.JobSd = utils.ParseNullString(JobSd)
	search.JobEd = utils.ParseNullString(JobEd)
	search.JobPmNm = utils.ParseNullString(UserName)

	list, err := h.Service.GetNonUsedProjectList(ctx, page, search, retrySearch)
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

	// 개수 조회
	count, err := h.Service.GetNonUsedProjectCount(ctx, search, retrySearch)
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
			List  entity.NonUsedProjects `json:"list"`
			Count int                    `json:"count"`
		}{List: *list, Count: count},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct, func: 현장 프로젝트 추가
type HandlerAddProject struct {
	Service service.ProjectService
}

func (h *HandlerAddProject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	project := entity.ReqProject{}

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
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

	err := h.Service.AddProject(ctx, project)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        DataAddFailed,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	rsp := Response{
		Result: Success,
	}
	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct, func: 현장 기본 프로젝트 변경
type HandlerModifyDefaultProject struct {
	Service service.ProjectService
}

func (h *HandlerModifyDefaultProject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project := entity.ReqProject{}
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
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

	err := h.Service.ModifyDefaultProject(ctx, project)
	if err != nil {
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

	rsp := Response{
		Result: Success,
	}
	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct, func: 현장 프로젝트 사용여부 변경
type HandlerModifyUseProject struct {
	Service service.ProjectService
}

func (h *HandlerModifyUseProject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project := entity.ReqProject{}
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
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

	err := h.Service.ModifyUseProject(ctx, project)
	if err != nil {
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

	rsp := Response{
		Result: Success,
	}
	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct, func: 현장 프로젝트 삭제
type HandlerRemoveProject struct {
	Service service.ProjectService
}

func (h *HandlerRemoveProject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	snoString := chi.URLParam(r, "sno")
	jnoString := chi.URLParam(r, "jno")
	if snoString == "" || jnoString == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}
	sno, _ := strconv.ParseInt(snoString, 10, 64)
	jno, _ := strconv.ParseInt(jnoString, 10, 64)

	err := h.Service.RemoveProject(ctx, sno, jno)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        DataRemoveFailed,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	rsp := Response{
		Result: Success,
	}
	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
