package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
	"net/http"
	"strconv"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-03-05
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct&func: 초단기 예보 api 호출
type HandlerWhetherSrtNcst struct {
	Service service.WhetherApiService
}

func (h *HandlerWhetherSrtNcst) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get parameter
	latStr := r.URL.Query().Get("lat") // 위도
	lonStr := r.URL.Query().Get("lon") // 경도
	if latStr == "" || lonStr == "" {
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

	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)
	// Lambert Conformal Conic(LCC) 투영법을 사용해야 함.
	nx, ny := utils.LatLonToXY(lat, lon)

	// 현재날짜, 시간
	now := time.Now()
	baseDate := now.Format("20060102")
	baseTime := now.Format("1504")

	res, err := h.Service.GetWhetherSrtNcst(baseDate, baseTime, nx, ny)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        CallApiFailed,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			List entity.WhetherSrtEntityRes `json:"list"`
		}{List: res},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}

// 기상청 기상특보통보문조회
type HandlerWhetherWrnMsg struct {
	Service service.WhetherApiService
}

func (h *HandlerWhetherWrnMsg) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.Service.GetWhetherWrnMsg()
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        CallApiFailed,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			List entity.WhetherWrnMsgList `json:"list"`
		}{List: res},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)

}
