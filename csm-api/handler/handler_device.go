package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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
// - response: http get paramter
func (d *DeviceListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// http get paramter를 저장할 구조체 생성
	page := entity.Page{}

	pageNum := r.URL.Query().Get("page_num")
	rowSize := r.URL.Query().Get("row_size")
	if pageNum == "" || rowSize == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "get parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}
	page.PageNum, _ = strconv.Atoi(pageNum)
	page.RowSize, _ = strconv.Atoi(rowSize)

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
// - response: entity.Device - json(raw)
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

// struct: 근태인식기 수정
type DeviceModifyHandler struct {
	Service service.DeviceService
}

// func: 근태인식기 수정
// @param
// - response: entity.Device - json(raw)
func (d *DeviceModifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// 근태인식기 수정
	err := d.Service.ModifyDevice(ctx, device)
	if err != nil {
		fmt.Println("handler add_task.go ServeHTTP Validator error")
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        DataModifyFailed,
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

// struct: 근태인식기 삭제
type DeviceRemoveHandler struct {
	Service service.DeviceService
}

// func: 근태인식기 삭제
// @param
// - response: url parameter
func (d *DeviceRemoveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// URL에서 {id} 값을 가져오기
	id := chi.URLParam(r, "id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        ParsingError,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}

	// 서비스 호출하여 삭제 처리
	err = d.Service.RemoveDevice(ctx, int64ID)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        DataRemoveFailed,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	// 성공 응답
	rsp := Response{
		Result: Success,
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
