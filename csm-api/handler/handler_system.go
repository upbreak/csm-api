package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type SystemHandler struct {
	WorkerService         service.WorkerService
	WorkHourService       service.WorkHourService
	ProjectService        service.ProjectService
	ProjectSettingService service.ProjectSettingService
	WeatherService        service.WeatherApiService
	SiteService           service.SiteService
}

// 근로자 마감 처리
func (h *SystemHandler) WorkerDeadline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.WorkerService.ModifyWorkerDeadlineInit(ctx); err != nil {
		_ = entity.WriteErrorLog(ctx, fmt.Errorf("ModifyWorkerDeadlineInit fail: %+v", err))
		FailResponse(ctx, w, err)
		return

	} else {
		log.Println("ModifyWorkerDeadlineInit completed")
	}

	SuccessResponse(ctx, w)
}

// 철야 확인 작업
func (h *SystemHandler) WorkerOverTime(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if count, err := h.WorkerService.ModifyWorkerOverTime(ctx); err != nil {
		_ = entity.WriteErrorLog(ctx, fmt.Errorf("ModifyWorkerOverTime fail: %+v", err))
		FailResponse(ctx, w, err)
		return

	} else if count != 0 {
		log.Println("ModifyWorkerOverTime completed")
	}

	SuccessResponse(ctx, w)
}

// 프로젝트 정보 업데이트(초기 세팅)
func (h *SystemHandler) ProjectInitSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if count, err := h.ProjectSettingService.CheckProjectSetting(ctx); err != nil {
		_ = entity.WriteErrorLog(ctx, fmt.Errorf("CheckProjectSettings fail: %+v", err))
		FailResponse(ctx, w, err)
		return

	} else if count != 0 {
		log.Printf("CheckProjectSettings %d completed \n", count)
	}

	SuccessResponse(ctx, w)
}

// 근로자 공수 계산(마감처리 안되고, 출퇴근이 모두 있는 근로자)
func (h *SystemHandler) UpdateWorkHour(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := entity.Base{
		ModUser: utils.ParseNullString("SYSTEM_BATCH"),
	}

	if err := h.WorkHourService.ModifyWorkHour(ctx, user); err != nil {
		_ = entity.WriteErrorLog(ctx, fmt.Errorf("ModifyWorkHour fail: %+v", err))
		FailResponse(ctx, w, err)
		return

	} else {
		log.Println("ModifyWorkHour completed")
	}

	SuccessResponse(ctx, w)
}

// 당일 공정률 기록
func (h *SystemHandler) SettingWorkRate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	targetDateString := r.URL.Query().Get("targetDate")
	if targetDateString == "" {
		BadRequestResponse(ctx, w)
		return
	}
	targetDate, err := time.Parse("2006-01-02", targetDateString)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	var count int64 = 0
	if count, err = h.SiteService.SettingWorkRate(ctx, targetDate); err != nil {
		log.Printf("SettingWorkRate fail: %w", err)
		FailResponse(ctx, w, err)
		return

	} else if count > 0 {
		log.Printf("SettingWorkRate success: %d", count)
	}

	SuccessResponse(ctx, w)
}

// 공수 추가
func (h *SystemHandler) AddManHour(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	manhour := entity.ManHour{}
	if err := json.NewDecoder(r.Body).Decode(&manhour); err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	if err := h.ProjectSettingService.AddManHour(ctx, manhour); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)

}
