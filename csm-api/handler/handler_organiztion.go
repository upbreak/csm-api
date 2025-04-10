package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
	"strconv"
)

type HandlerOrganization struct {
	Service service.OrganizationService
}

func (h *HandlerOrganization) ByProjectList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jno := r.PathValue("jno")
	int64JNO, err := strconv.ParseInt(jno, 10, 64)

	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 고객사 조직도 조회
	client, err := h.Service.GetOrganizationClientList(ctx, int64JNO)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 계약사 조직도 조회
	hitech, err := h.Service.GetOrganizationHtencList(ctx, int64JNO)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		Client entity.OrganizationPartitions `json:"client"`
		Hitech entity.OrganizationPartitions `json:"hitech"`
	}{Client: *client, Hitech: *hitech}
	SuccessValuesResponse(ctx, w, values)

}
