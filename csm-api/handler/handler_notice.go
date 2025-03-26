package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
	"strconv"
)

/**
 * @author 작성자: 정지영
 * @created 작성일: 2025-02-17
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct : 공지사항 조회
type NoticeListHandler struct {
	Service service.NoticeService
}

// func : 공지사항 전체조회
// @param
// - response: hhtp get parameter
func (n *NoticeListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
				Message:        err.Error(),
				Details:        ParsingError,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	page := entity.Page{}
	search := entity.Notice{}

	pageNum := r.URL.Query().Get(entity.PageNumKey)
	rowSize := r.URL.Query().Get(entity.RowSizeKey)
	order := r.URL.Query().Get(entity.OrderKey)

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

	search.Jno, _ = strconv.ParseInt(r.URL.Query().Get("jno"), 10, 64)
	search.JobLocName = r.URL.Query().Get("loc_name")
	search.JobName = r.URL.Query().Get("site_nm")
	search.Title = r.URL.Query().Get("title")
	search.UserInfo = r.URL.Query().Get("user_info")

	notices, err := n.Service.GetNoticeList(ctx, int64UNO, role, page, search)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	count, err := n.Service.GetNoticeListCount(ctx, int64UNO, role, search)
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
			Notices entity.Notices `json:"notices"`
			Count   int            `json:"count"`
		}{Notices: *notices, Count: count},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)

}

// struct: 공지사항 추가
type NoticeAddHandler struct {
	Service service.NoticeService
}

// func: 공지사항 추가
// @param
// - request: entity.Notice - json(raw)
func (n *NoticeAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	notice := entity.Notice{}
	if err := json.NewDecoder(r.Body).Decode(&notice); err != nil {
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

	err := n.Service.AddNotice(ctx, notice)
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

	RespondJSON(ctx, w, &Response{Result: Success}, http.StatusOK)

}

// 공지사항 수정
type NoticeModifyHandler struct {
	Service service.NoticeService
}

// func: 공지사항 수정
// @param
// - request: entity.Notice - json(raw)
func (n *NoticeModifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// request 데이터 파싱
	notice := entity.Notice{}
	if err := json.NewDecoder(r.Body).Decode(&notice); err != nil {
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

	if err := n.Service.ModifyNotice(ctx, notice); err != nil {
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

// 공지사항 삭제
type NoticeDeleteHandler struct {
	Service service.NoticeService
}

// func: 공지사항 삭제
// @param
// - idx : 공지사항 인덱스
func (n *NoticeDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idx := r.PathValue("idx")

	int64IDX, err := strconv.ParseInt(idx, 10, 64)

	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        ParsingError,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)

		return
	}

	if err := n.Service.RemoveNotice(ctx, int64IDX); err != nil {
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

	RespondJSON(ctx, w, &Response{Result: Success}, http.StatusOK)

}

// 공지기간 조회
type NoticePeriodHandler struct {
	Service service.NoticeService
}

// func: 공지기간 조회
// @params
// -
func (n *NoticePeriodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	noticePeriods, err := n.Service.GetNoticePeriod(ctx)
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
			Periods entity.NoticePeriods `json:"periods"`
		}{Periods: *noticePeriods},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
