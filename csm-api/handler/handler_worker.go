package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
	"strconv"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-17
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

type HandlerWorkerTotalList struct {
	Service service.WorkerService
}

// func: 전체 근로자 조회
// @param
// - response: http get paramter
func (h *HandlerWorkerTotalList) ServeHttp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// http get paramter를 저장할 구조체 생성 및 파싱
	page := entity.Page{}
	search := entity.Worker{}

	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	order := r.URL.Query().Get("order")
	rnumOrder := r.URL.Query().Get("rnum_order")
	jobName := r.URL.Query().Get("job_name")
	userId := r.URL.Query().Get("user_id")
	userNm := r.URL.Query().Get("user_nm")
	department := r.URL.Query().Get("department")
	phone := r.URL.Query().Get("phone")
	workerType := r.URL.Query().Get("worker_type")

	retrySearch := r.URL.Query().Get("retry_search")

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
	search.JobName = jobName
	search.UserId = userId
	search.UserNm = userNm
	search.Department = department
	search.Phone = phone
	search.WorkerType = workerType

	// 조회
	list, err := h.Service.GetWorkerTotalList(ctx, page, search, retrySearch)
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
	count, err := h.Service.GetWorkerTotalCount(ctx, search, retrySearch)
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
			List  entity.Workers `json:"list"`
			Count int            `json:"count"`
		}{List: *list, Count: count},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

type HandlerWorkerSiteBaseList struct {
	Service service.WorkerService
}

// struct, func: 근로자 검색(현장근로자 추가시 사용)
type HandlerWorkerByUserId struct {
	Service service.WorkerService
}

func (h *HandlerWorkerByUserId) ServeHttp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page := entity.Page{}
	search := entity.WorkerDaily{}
	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	retrySearch := r.URL.Query().Get("retry_search")
	searchStartTime := r.URL.Query().Get("search_start_time")
	jno := r.URL.Query().Get("jno")

	if pageNum == "" || rowSize == "" || jno == "" {
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
	search.Jno, _ = strconv.ParseInt(jno, 10, 64)
	search.SearchStartTime = searchStartTime

	list, err := h.Service.GetWorkerListByUserId(ctx, page, search, retrySearch)
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

	count, err := h.Service.GetWorkerCountByUserId(ctx, search, retrySearch)
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
			List  entity.Workers `json:"list"`
			Count int            `json:"count"`
		}{List: *list, Count: count},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct, func: 근로자 추가
type HandlerWorkerAdd struct {
	Service service.WorkerService
}

// @param
// - http method: post
// - param: entity.Worker - json(raw)
func (h *HandlerWorkerAdd) ServeHttp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//데이터 파싱
	worker := entity.Worker{}
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
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

	err := h.Service.AddWorker(ctx, worker)
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

// struct, func: 근로자 수정
type HandlerWorkerMod struct {
	Service service.WorkerService
}

// @param
// - http method: put
// - param: entity.Worker - json(raw)
func (h *HandlerWorkerMod) ServeHttp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//데이터 파싱
	worker := entity.Worker{}
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
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

	err := h.Service.ModifyWorker(ctx, worker)
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

// func: 현장 근로자 조회
// @param
// - response: http get paramter
func (h *HandlerWorkerSiteBaseList) ServeHttp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// http get paramter를 저장할 구조체 생성 및 파싱
	page := entity.Page{}
	search := entity.WorkerDaily{}

	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	order := r.URL.Query().Get("order")
	rnumOrder := r.URL.Query().Get("rnum_order")
	retrySearch := r.URL.Query().Get("retry_search")
	jno := r.URL.Query().Get("jno")
	userId := r.URL.Query().Get("user_id")
	userNm := r.URL.Query().Get("user_nm")
	department := r.URL.Query().Get("department")
	searchStartTime := r.URL.Query().Get("search_start_time")
	searchEndTime := r.URL.Query().Get("search_end_time")

	if pageNum == "" || rowSize == "" || searchStartTime == "" || searchEndTime == "" || jno == "" {
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
	search.Jno, _ = strconv.ParseInt(jno, 10, 64)
	search.UserId = userId
	search.UserNm = userNm
	search.Department = department
	search.SearchStartTime = searchStartTime
	search.SearchEndTime = searchEndTime

	// 조회
	list, err := h.Service.GetWorkerSiteBaseList(ctx, page, search, retrySearch)
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
	count, err := h.Service.GetWorkerSiteBaseCount(ctx, search, retrySearch)
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
			List  entity.WorkerDailys `json:"list"`
			Count int                 `json:"count"`
		}{List: *list, Count: count},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct, func: 현장근로자 추가/수정
type HandlerSiteBaseMerge struct {
	Service service.WorkerService
}

// @param
// - http method: post
// - param: entity.WorkerDailys - json(raw)
func (h *HandlerSiteBaseMerge) ServeHttp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//데이터 파싱
	workers := entity.WorkerDailys{}
	if err := json.NewDecoder(r.Body).Decode(&workers); err != nil {
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

	err := h.Service.MergeSiteBaseWorker(ctx, workers)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        DataMergeFailed,
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

// struct, func: 근로자 일괄마감
type HandlerWorkerDeadline struct {
	Service service.WorkerService
}

// @param
// - http method: post
// - param: entity.WorkerDailys - json(raw)
func (h *HandlerWorkerDeadline) ServeHttp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//데이터 파싱
	workers := entity.WorkerDailys{}
	if err := json.NewDecoder(r.Body).Decode(&workers); err != nil {
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

	err := h.Service.ModifyWorkerDeadline(ctx, workers)
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

// struct, func: 현장 근로자 프로젝트 변경
type HandlerWorkerProject struct {
	Service service.WorkerService
}

// @param
// - http method: post
// - param: entity.WorkerDailys - json(raw)
func (h *HandlerWorkerProject) ServeHttp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//데이터 파싱
	workers := entity.WorkerDailys{}
	if err := json.NewDecoder(r.Body).Decode(&workers); err != nil {
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

	err := h.Service.ModifyWorkerProject(ctx, workers)
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
