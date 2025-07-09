package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
	"strconv"
)

type HandlerUserRole struct {
	Service service.UserRoleService
}

// 사용자 권한 조회
// param HTTP GET uno (사용자번호)
func (h *HandlerUserRole) GetUserRoleListByUno(w http.ResponseWriter, r *http.Request) {
	unoString := r.URL.Query().Get("uno")
	if unoString == "" {
		BadRequestResponse(r.Context(), w)
		return
	}
	uno, _ := strconv.ParseInt(unoString, 10, 64)
	userRoles, err := h.Service.GetUserRoleListByUno(r.Context(), uno)
	if err != nil {
		FailResponse(r.Context(), w, err)
	}
	SuccessValuesResponse(r.Context(), w, userRoles)
}

// 사용자 권한 추가
func (h *HandlerUserRole) AddUserRole(w http.ResponseWriter, r *http.Request) {
	itemLog, userRoles, err := entity.DecodeItem(r, []entity.UserRoleMap{})
	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	err = h.Service.AddUserRole(r.Context(), userRoles)
	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	entity.WriteLog(itemLog)
	SuccessResponse(r.Context(), w)
}

// 사용자 권한 삭제
func (h *HandlerUserRole) RemoveUserRole(w http.ResponseWriter, r *http.Request) {
	itemLog, userRoles, err := entity.DecodeItem(r, []entity.UserRoleMap{})
	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	err = h.Service.RemoveUserRole(r.Context(), userRoles)
	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	entity.WriteLog(itemLog)
	SuccessResponse(r.Context(), w)
}
