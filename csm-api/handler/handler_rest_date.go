package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
	"strconv"
	"time"
)

type HandlerRestDate struct {
	Service service.RestDateApiService
}

func (h *HandlerRestDate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")

	if year == "" {
		now := time.Now()
		year = strconv.Itoa(now.Year())
	}

	list, err := h.Service.GetRestDelDates(year, month)
	if err != nil {
		FailResponse(ctx, w, err)
	}

	values := struct {
		List entity.RestDates `json:"list"`
	}{List: list}

	SuccessValuesResponse(ctx, w, values)
}
