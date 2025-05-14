package handler

import (
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
		BadRequestResponse(ctx, w)
		return
	}

	mapPoint, err := s.Service.GetAPISiteMapPoint(roadAddress)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessValuesResponse(ctx, w, mapPoint)
}
