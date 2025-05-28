package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
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

type HandlerWorker struct {
	Service service.WorkerService
}

// func: 전체 근로자 조회
// @param
// - response: http get paramter
func (h *HandlerWorker) TotalList(w http.ResponseWriter, r *http.Request) {
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
	discName := r.URL.Query().Get("disc_name")

	retrySearch := r.URL.Query().Get("retry_search")

	if pageNum == "" || rowSize == "" {
		BadRequestResponse(ctx, w)
		return
	}

	page.PageNum, _ = strconv.Atoi(pageNum)
	page.RowSize, _ = strconv.Atoi(rowSize)
	page.Order = order
	page.RnumOrder = rnumOrder
	search.JobName = utils.ParseNullString(jobName)
	search.UserId = utils.ParseNullString(userId)
	search.UserNm = utils.ParseNullString(userNm)
	search.Department = utils.ParseNullString(department)
	search.Phone = utils.ParseNullString(phone)
	search.WorkerType = utils.ParseNullString(workerType)
	search.DiscName = utils.ParseNullString(discName)

	// 조회
	list, err := h.Service.GetWorkerTotalList(ctx, page, search, retrySearch)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 개수 조회
	count, err := h.Service.GetWorkerTotalCount(ctx, search, retrySearch)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List  entity.Workers `json:"list"`
		Count int            `json:"count"`
	}{List: *list, Count: count}
	SuccessValuesResponse(ctx, w, values)
}

// func: 출근 안한 근로자 검색
// @param
// -
func (h *HandlerWorker) AbsentList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page := entity.Page{}
	search := entity.WorkerDaily{}
	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	retrySearch := r.URL.Query().Get("retry_search")
	searchStartTime := r.URL.Query().Get("search_start_time")
	jno := r.URL.Query().Get("jno")

	if pageNum == "" || rowSize == "" || jno == "" {
		BadRequestResponse(ctx, w)
		return
	}
	page.PageNum, _ = strconv.Atoi(pageNum)
	page.RowSize, _ = strconv.Atoi(rowSize)
	search.Jno = utils.ParseNullInt(jno)
	search.SearchStartTime = utils.ParseNullString(searchStartTime)

	list, err := h.Service.GetAbsentWorkerList(ctx, page, search, retrySearch)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	count, err := h.Service.GetAbsentWorkerCount(ctx, search, retrySearch)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List  entity.Workers `json:"list"`
		Count int            `json:"count"`
	}{List: *list, Count: count}
	SuccessValuesResponse(ctx, w, values)
}

// 프로젝트에 참여한 회사명 리스트
func (h *HandlerWorker) DepartList(w http.ResponseWriter, r *http.Request) {
	jnoString := r.URL.Query().Get("jno")
	if jnoString == "" {
		BadRequestResponse(r.Context(), w)
		return
	}

	jno, _ := strconv.ParseInt(jnoString, 10, 64)
	list, err := h.Service.GetWorkerDepartList(r.Context(), jno)
	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}
	SuccessValuesResponse(r.Context(), w, list)
}

// func: 근로자 추가
// @param
// - http method: post
// - param: entity.Worker - json(raw)
func (h *HandlerWorker) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//데이터 파싱
	worker := entity.Worker{}
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	err := h.Service.AddWorker(ctx, worker)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}

// func: 근로자 수정
// @param
// - http method: put
// - param: entity.Worker - json(raw)
func (h *HandlerWorker) Modify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//데이터 파싱
	worker := entity.Worker{}
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	err := h.Service.ModifyWorker(ctx, worker)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}

// func: 현장 근로자 조회
// @param
// - response: http get paramter
func (h *HandlerWorker) SiteBaseList(w http.ResponseWriter, r *http.Request) {
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
		BadRequestResponse(ctx, w)
		return
	}

	page.PageNum, _ = strconv.Atoi(pageNum)
	page.RowSize, _ = strconv.Atoi(rowSize)
	page.Order = order
	page.RnumOrder = rnumOrder
	search.Jno = utils.ParseNullInt(jno)
	search.UserId = utils.ParseNullString(userId)
	search.UserNm = utils.ParseNullString(userNm)
	search.Department = utils.ParseNullString(department)
	search.SearchStartTime = utils.ParseNullString(searchStartTime)
	search.SearchEndTime = utils.ParseNullString(searchEndTime)

	// 조회
	list, err := h.Service.GetWorkerSiteBaseList(ctx, page, search, retrySearch)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 개수 조회
	count, err := h.Service.GetWorkerSiteBaseCount(ctx, search, retrySearch)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List  entity.WorkerDailys `json:"list"`
		Count int                 `json:"count"`
	}{List: *list, Count: count}
	SuccessValuesResponse(ctx, w, values)
}

// func: 현장근로자 추가/수정
// @param
// - http method: post
// - param: entity.WorkerDailys - json(raw)
func (h *HandlerWorker) Merge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//데이터 파싱
	workers := entity.WorkerDailys{}
	if err := json.NewDecoder(r.Body).Decode(&workers); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	err := h.Service.MergeSiteBaseWorker(ctx, workers)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}
	SuccessResponse(ctx, w)
}

// func: 근로자 일괄마감
// @param
// - http method: post
// - param: entity.WorkerDailys - json(raw)
func (h *HandlerWorker) ModifyDeadline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//데이터 파싱
	workers := entity.WorkerDailys{}
	if err := json.NewDecoder(r.Body).Decode(&workers); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	err := h.Service.ModifyWorkerDeadline(ctx, workers)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}

// func: 현장 근로자 프로젝트 변경
// @param
// - http method: post
// - param: entity.WorkerDailys - json(raw)
func (h *HandlerWorker) ModifyProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//데이터 파싱
	workers := entity.WorkerDailys{}
	if err := json.NewDecoder(r.Body).Decode(&workers); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	err := h.Service.ModifyWorkerProject(ctx, workers)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}
