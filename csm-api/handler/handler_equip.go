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

func (h *HandlerEquip) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := h.Service.GetEquipList(ctx)
	if err != nil {
		FailResponse(ctx, w, err)
	}

	values := struct {
		List entity.EquipTemps `json:"list"`
	}{List: list}
	SuccessValuesResponse(ctx, w, values)
}

func (h *HandlerEquip) Merge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	equips := entity.EquipTemps{}

	if err := json.NewDecoder(r.Body).Decode(&equips); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	if err := h.Service.MergeEquipCnt(ctx, equips); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}
