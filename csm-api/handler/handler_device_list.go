package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"fmt"
	"net/http"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct: 근태인식기 조회
type DeviceListHandler struct {
	Service service.DeviceService
}

// func: 근태인식기 조회
// @param
// - response: entity.Page
func (d *DeviceListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// body의 raw data를 저장할 구조체 생성
	page := entity.Page{}

	// body 데이터 파싱
	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        BodyDataParseError,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	// 근태인식기 목록
	list, err := d.Service.GetDeviceList(ctx, page)
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

	// 근태인식기 전체 개수
	count, err := d.Service.GetDeviceListCount(ctx)
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
			List  entity.Devices `json:"list"`
			Count int            `json:"count"`
		}{List: *list, Count: count},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// struct: 근태인식기 추가
type DeviceAddHandler struct {
	Service service.DeviceService
}

// func: 근태인식기 추가
// @param
// - response: entity.Device
func (d *DeviceAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//response 데이터 파싱
	device := entity.Device{}
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        BodyDataParseError,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	// 근태인식기 추가
	err := d.Service.AddDevice(ctx, device)
	if err != nil {
		fmt.Println("handler add_task.go ServeHTTP Validator error")
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        DataAddFailed,
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
