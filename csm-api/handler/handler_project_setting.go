package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type HandlerProjectSetting struct {
	Service service.ProjectSettingService
}

// func: 프로젝트 기본 설정 정보 확인
// @param
// - jno: 프로젝트pk
func (h *HandlerProjectSetting) ProjectSettingList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jnoString := r.PathValue("jno")
	if jnoString == "" {
		BadRequestResponse(ctx, w)
		return
	}
	jno, _ := strconv.ParseInt(jnoString, 10, 64)

	setting, err := h.Service.GetProjectSetting(ctx, jno)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}
	values := struct {
		Project entity.ProjectSettings `json:"project"`
	}{Project: *setting}

	SuccessValuesResponse(ctx, w, values)
}

// func: 프로젝트 기본 설정 추가 및 수정
// @param
// - projectSetting
func (h *HandlerProjectSetting) MergeProjectSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	setting := entity.ProjectSetting{}

	if err := json.NewDecoder(r.Body).Decode(&setting); err != nil {
		BadRequestResponse(ctx, w)
		return
	}

	err := h.Service.MergeProjectSetting(ctx, setting)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}

// func: 프로젝트 공수 정보 확인
// @param
// - jno: 프로젝트pk
func (h *HandlerProjectSetting) ManHourList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jnoString := r.PathValue("jno")
	if jnoString == "" {
		BadRequestResponse(ctx, w)
		return
	}
	jno, _ := strconv.ParseInt(jnoString, 10, 64)

	manhours, err := h.Service.GetManHourList(ctx, jno)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}
	values := struct {
		ManHours entity.ManHours `json:"man_hours"`
	}{ManHours: *manhours}

	SuccessValuesResponse(ctx, w, values)
}

// func: 공수 설정 추가 및 수정
// @param
// - mamHours
func (h *HandlerProjectSetting) MergeManHours(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	manhours := entity.ManHours{}
	
	if err := json.NewDecoder(r.Body).Decode(&manhours); err != nil {
		BadRequestResponse(ctx, w)
		return
	}

	err := h.Service.MergeManHours(ctx, &manhours)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}

// func: 공수 삭제
// @param
// - mhno: 공수pk
func (h *HandlerProjectSetting) DeleteManHour(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	mhnoString := r.PathValue("mhno")
	if mhnoString == "" {
		BadRequestResponse(ctx, w)
		return
	}
	mhno, _ := strconv.ParseInt(mhnoString, 10, 64)

	err := h.Service.DeleteManHour(ctx, mhno)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}
