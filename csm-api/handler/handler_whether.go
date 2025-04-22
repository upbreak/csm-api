package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
	"net/http"
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
	Service        service.WhetherApiService
	SitePosService service.SitePosService
}

func (h *HandlerWhetherSrtNcst) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := h.SitePosService.GetSitePosList(ctx)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	var whetherList entity.WhetherSrtRes
	for _, item := range list {
		now := time.Now()
		baseDate := now.Format("20060102")
		baseTime := now.Add(time.Minute * -30).Format("1504") // 기상청에서 30분 단위로 발표하기 때문에 30분 전의 데이터 요청
		nx, ny := utils.LatLonToXY(item.Latitude.Float64, item.Longitude.Float64)

		res, err := h.Service.GetWhetherSrtNcst(baseDate, baseTime, nx, ny)
		if err != nil {
			FailResponse(ctx, w, err)
			return
		}

		whether := entity.WhetherSrt{}
		whether.Whether = res
		whether.Sno = item.Sno.Int64
		whetherList = append(whetherList, whether)
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			List entity.WhetherSrtRes `json:"list"`
		}{List: whetherList},
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
