package main

import (
	"context"
	"csm-api/clock"
	"csm-api/service"
	"csm-api/store"
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
	WorkerService service.WorkerService
	cron          *cron.Cron
}

func NewScheduler(safeDb *sqlx.DB) (*Scheduler, error) {
	r := store.Repository{Clocker: clock.RealClock{}}
	c := cron.New(cron.WithSeconds())

	scheduler := &Scheduler{
		WorkerService: &service.ServiceWorker{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
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

	// ... 추가 job 등록

	s.cron.Start()

	log.Println("[Scheduler] Cron started")

	// ctx.Done() 기다리다가 종료
	<-ctx.Done()
	log.Println("[Scheduler] Stopping scheduler...")

	s.cron.Stop()
	return nil
}
