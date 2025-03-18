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
	search := entity.Worker{}

	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	order := r.URL.Query().Get("order")
	sno := r.URL.Query().Get("sno")
	jobName := r.URL.Query().Get("job_name")
	userNm := r.URL.Query().Get("user_nm")
	department := r.URL.Query().Get("department")
	searchStartTime := r.URL.Query().Get("search_start_time")
	searchEndTime := r.URL.Query().Get("search_end_time")

	if pageNum == "" || rowSize == "" || searchStartTime == "" || searchEndTime == "" || sno == "" {
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
	search.Sno, _ = strconv.ParseInt(sno, 10, 64)
	search.JobName = jobName
	search.UserNm = userNm
	search.Department = department

	// 조회
	list, err := h.Service.GetWorkerSiteBaseList(ctx, page, search)
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
	count, err := h.Service.GetWorkerSiteBaseCount(ctx, search)
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
