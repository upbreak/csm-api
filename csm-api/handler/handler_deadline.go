package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
	"fmt"
	"net/http"
)

type HandlerDeadline struct {
	UploadService service.UploadFileService
}

// 일일마감 엑셀 자료 정보
func (h *HandlerDeadline) UploadFileList(w http.ResponseWriter, r *http.Request) {
	jno := r.URL.Query().Get("jno")
	workDate := r.URL.Query().Get("work_date")
	if jno == "" || workDate == "" {
		FailResponse(r.Context(), w, utils.CustomErrorf(fmt.Errorf("jno or work_date is empty")))
		return
	}

	file := entity.UploadFile{
		Jno:      utils.ParseNullInt(jno),
		WorkDate: utils.ParseNullDate(workDate),
	}
	list, err := h.UploadService.GetUploadFileList(r.Context(), file)
	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	SuccessValuesResponse(r.Context(), w, list)
}
