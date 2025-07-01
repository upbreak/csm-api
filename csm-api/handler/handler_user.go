package handler

import (
	"csm-api/service"
	"net/http"
	"strconv"
)

type HandlerUser struct {
	Service service.UserService
}

// 사용자 권한 조회(프로젝트 선택시 사용)
func (h *HandlerUser) UserRole(w http.ResponseWriter, r *http.Request) {
	jnoString := r.URL.Query().Get("jno")
	unoString := r.URL.Query().Get("uno")
	if jnoString == "" || unoString == "" {
		BadRequestResponse(r.Context(), w)
		return
	}
	jno, _ := strconv.ParseInt(jnoString, 10, 64)
	uno, _ := strconv.ParseInt(unoString, 10, 64)
	role, err := h.Service.GetUserRole(r.Context(), jno, uno)

	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	SuccessValuesResponse(r.Context(), w, role)
}
