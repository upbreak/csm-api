package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
	"fmt"
	"log"
	"net/http"
	"time"
)

type InitApiHandler struct {
	WorkerService         service.WorkerService
	WorkHourService       service.WorkHourService
	ProjectService        service.ProjectService
	ProjectSettingService service.ProjectSettingService
	WeatherService        service.WeatherApiService
	SiteService           service.SiteService
}

func (h *InitApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("[InitApi] Running InitApi")

	// 근로자 마감 처리
	if err := h.WorkerService.ModifyWorkerDeadlineInit(ctx); err != nil {
		_ = entity.WriteErrorLog(ctx, utils.CustomErrorf(fmt.Errorf("[InitApi] ModifyWorkerDeadlineInit fail: %+v", err)))
	} else {
		log.Println("[InitApi] ModifyWorkerDeadlineInit completed")
	}

	// 철야 확인 작업
	if count, err := h.WorkerService.ModifyWorkerOverTime(ctx); err != nil {
		_ = entity.WriteErrorLog(ctx, utils.CustomErrorf(fmt.Errorf("[InitApi] ModifyWorkerOverTime fail: %+v", err)))
	} else if count != 0 {
		log.Println("[InitApi] ModifyWorkerOverTime completed")
	}

	// 프로젝트 정보 업데이트(초기 세팅)
	if count, err := h.ProjectSettingService.CheckProjectSetting(ctx); err != nil {
		_ = entity.WriteErrorLog(ctx, utils.CustomErrorf(fmt.Errorf("[InitApi] CheckProjectSettings fail: %+v", err)))
	} else if count != 0 {
		log.Printf("[InitApi] CheckProjectSettings %d completed \n", count)
	}

	// 근로자 공수 계산(마감처리 안되고, 출퇴근이 모두 있는 근로자)
	user := entity.Base{
		ModUser: utils.ParseNullString("SYSTEM_BATCH"),
	}
	if err := h.WorkHourService.ModifyWorkHour(ctx, user); err != nil {
		_ = entity.WriteErrorLog(ctx, utils.CustomErrorf(fmt.Errorf("[InitApi] ModifyWorkHour fail: %+v", err)))
	} else {
		log.Println("[InitApi] ModifyWorkHour completed")
	}

	// 날씨 저장
	if err := h.WeatherService.SaveWeather(ctx); err != nil {
		_ = entity.WriteErrorLog(ctx, fmt.Errorf("[InitApi] SaveWeather fail: %w", err))
	} else {
		log.Printf("[InitApi] SaveWeather completed")
	}

	// 당일 공정률 기록
	now := time.Now()
	if count, err := h.SiteService.SettingWorkRate(ctx, now); err != nil {
		log.Printf("[InitApi] SettingWorkRate fail: %w", err)
	} else if count > 0 {
		log.Printf("[InitApi] SettingWorkRate success: %d", count)
	}

	SuccessResponse(ctx, w)
}
