package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
)

type HandlerEquip struct {
	Service service.EquipService
}

func (h *HandlerEquip) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	equips := entity.EquipTemps{}

	if err := json.NewDecoder(r.Body).Decode(&equips); err != nil {
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

	if err := h.Service.MergeEquipCnt(ctx, equips); err != nil {
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
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
