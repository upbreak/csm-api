package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
)

type DeviceListHandler struct {
	Service service.DeviceService
}

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
