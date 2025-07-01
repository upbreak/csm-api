package handler

import (
	"csm-api/service"
	"net/http"
	"strings"
)

type HandlerMenu struct {
	Service service.MenuService
}

// 권한별 메뉴 리스트
func (h *HandlerMenu) List(w http.ResponseWriter, r *http.Request) {
	rolesPipe := r.URL.Query().Get("roles")
	if rolesPipe == "" {
		BadRequestResponse(r.Context(), w)
		return
	}

	var roles []string
	rolesArr := strings.Split(rolesPipe, "|")
	for _, role := range rolesArr {
		roles = append(roles, strings.TrimSpace(role))
	}

	list, err := h.Service.GetMenu(r.Context(), roles)
	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	SuccessValuesResponse(r.Context(), w, list)
}
