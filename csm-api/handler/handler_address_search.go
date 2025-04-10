package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
)

// 지도 x, y좌표 조회
type HandlerRoadAddress struct {
	Service service.AddressSearchAPIService
}

func (s *HandlerRoadAddress) AddressPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roadAddress := r.URL.Query().Get("roadAddress")

	if roadAddress == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "roadAddress parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}

	mapPoint, err := s.Service.GetAPISiteMapPoint(roadAddress)
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
			MapPoint entity.MapPoint `json:"point"`
		}{MapPoint: *mapPoint},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
