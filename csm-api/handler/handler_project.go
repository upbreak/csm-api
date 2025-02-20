package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
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
