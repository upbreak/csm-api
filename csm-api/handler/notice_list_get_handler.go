package handler

import (
	"csm-api/service"
	"net/http"
)

type ListNotice struct {
	Service service.NoticeService
}

func (ln *ListNotice) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rsp, err := ln.Service.GetNoticeList(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
