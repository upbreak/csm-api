package main

import (
	"context"
	"csm-api/clock"
	"csm-api/config"
	"csm-api/entity"
	"csm-api/service"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"log"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-04-28
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @description: 스캐줄러 설정
 * -
 */

type Scheduler struct {
	WorkerService         service.WorkerService
	WorkHourService       service.WorkHourService
	ProjectService        service.ProjectService
	ProjectSettingService service.ProjectSettingService
	WeatherService        service.WeatherApiService
	cron                  *cron.Cron
}

func NewScheduler(safeDb *sqlx.DB, apiCfg *config.ApiConfig) (*Scheduler, error) {
	r := store.Repository{Clocker: clock.RealClock{}}
	c := cron.New(cron.WithSeconds())

	scheduler := &Scheduler{
		WorkerService: &service.ServiceWorker{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
		WorkHourService: &service.ServiceWorkHour{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
		ProjectService: &service.ServiceProject{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
		ProjectSettingService: &service.ServiceProjectSetting{
			SafeDB:        safeDb,
			SafeTDB:       safeDb,
			Store:         &r,
			WorkHourStore: &r,
		},
		WeatherService: &service.ServiceWeather{
			ApiKey:       apiCfg,
			SafeDB:       safeDb,
			SafeTDB:      safeDb,
			Store:        &r,
			SitePosStore: &r,
		},

		cron: c,
	}

	return scheduler, nil
}

func (s *Scheduler) Run(ctx context.Context) error {
	// 0시 0분 0초에 실행
	// 근로자 마감 처리 (퇴근한 근로자만 처리)
	_, err := s.cron.AddFunc("0 0 0 * * *", func() {
		log.Println("[Scheduler] Running ModifyWorkerDeadlineSchedule")

		if err := s.WorkerService.ModifyWorkerDeadlineInit(ctx); err != nil {
			log.Printf("[Scheduler] ModifyWorkerDeadlineSchedule fail: %+v", err)
		} else {
			log.Println("[Scheduler] ModifyWorkerDeadlineSchedule completed")
		}
	})
	if err != nil {
		return fmt.Errorf("[Scheduler] failed to add cron job: %w", err)
	}

	// 1분마다 실행
	// 철야 확인 작업
	_, err = s.cron.AddFunc("0 0/1 * * * *", func() {
		var count int
		if count, err = s.WorkerService.ModifyWorkerOverTime(ctx); err != nil {
			log.Printf("[Scheduler] ModifyWorkerOverTime fail: %+v", err)
		} else if count != 0 {
			log.Println("[Scheduler] ModifyWorkerOverTime completed")
		}
	})
	if err != nil {
		// TODO: 에러아카이브
		return fmt.Errorf("[Scheduler] failed to add cron job: %w", err)
	}

	// 5분 마다 실행
	// 프로젝트 정보 업데이트(초기 세팅)
	_, err = s.cron.AddFunc("0 0/5 * * * *", func() {
		var count int
		if count, err = s.ProjectSettingService.CheckProjectSetting(ctx); err != nil {
			log.Printf("[Scheduler] CheckProjectSettings fail: %+v", err)
		} else if count != 0 {
			log.Println("[Scheduler] CheckProjectSettings completed")
		}
	})
	if err != nil {
		// TODO: 에러아카이브
		return fmt.Errorf("[Scheduler] failed to add cron job: %w", err)
	}

	// 0시 1분 0초에 실행
	// 근로자 공수 계산 (마감 처리가 안되고 출퇴근이 모두 있는 근로자)
	_, err = s.cron.AddFunc("0 1 0 * * *", func() {
		log.Println("[Scheduler] Running ModifyWorkHour")
		user := entity.Base{
			ModUser: utils.ParseNullString("SYSTEM_BATCH"),
		}
		if err = s.WorkHourService.ModifyWorkHour(ctx, user); err != nil {
			log.Printf("[Scheduler] ModifyWorkHour fail: %+v", err)
		} else {
			log.Println("[Scheduler] ModifyWorkHour completed")
		}
	})
	if err != nil {
		return fmt.Errorf("[Scheduler] failed to add cron job: %w", err)
	}

	_, err = s.cron.AddFunc("0 0 8,10,13,15,16,17 * * *", func() {
		log.Println("[Scheduler] 날씨 저장")

		err = s.WeatherService.SaveWeather(ctx)
		if err != nil {
			log.Printf("[Scheduler] GetWeatherSrtNcst fail: %w", err)
		}

	})
	if err != nil {
		return fmt.Errorf("[Scheduler] failed to add cron job: %w", err)
	}

	// ... 추가 job 등록
	s.cron.Start()

	log.Println("[Scheduler] Cron started")

	// ctx.Done() 기다리다가 종료
	<-ctx.Done()
	log.Println("[Scheduler] Stopping scheduler...")

	s.cron.Stop()
	return nil
}
