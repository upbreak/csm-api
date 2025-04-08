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
	role := utils.ParseNullString(r.URL.Query().Get("role"))

	intUNO := utils.ParseNullInt(uno)

	if uno == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "uno parameter is missing",
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

	search.Jno = utils.ParseNullInt(r.URL.Query().Get("jno"))
	search.JobLocName = utils.ParseNullString(r.URL.Query().Get("job_loc_name"))
	search.JobName = utils.ParseNullString(r.URL.Query().Get("job_name"))
	search.Title = utils.ParseNullString(r.URL.Query().Get("title"))
	search.UserInfo = utils.ParseNullString(r.URL.Query().Get("user_info"))

	notices, err := n.Service.GetNoticeList(ctx, intUNO, role, page, search)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	count, err := n.Service.GetNoticeListCount(ctx, intUNO, role, search)
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

	intIdx := utils.ParseNullInt(r.PathValue("idx"))

	if intIdx.Valid == false {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "idx parameter is missing",
				Details:        ParsingError,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)

		return
	}

	if err := n.Service.RemoveNotice(ctx, intIdx); err != nil {
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
