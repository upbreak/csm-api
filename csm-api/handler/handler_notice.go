package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
	"strconv"
)

type NoticeListHandler struct {
	Service service.NoticeService
}

// func : 공지사항 전체조회
func (n *NoticeListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page := entity.Page{}

	pageNum := r.URL.Query().Get(entity.PageNumKey)
	rowSize := r.URL.Query().Get(entity.RowSizeKey)

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

	notices, err := n.Service.GetNoticeList(ctx, page)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	count, err := n.Service.GetNoticeListCount(ctx)
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

type NoticeAddHandler struct {
	Service service.NoticeService
}
