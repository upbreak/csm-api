package main

import (
	"context"
	"csm-api/clock"
	"csm-api/service"
	"csm-api/store"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
	"log"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-04-28
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @description: 서버 실행시 초기화를 위한 설정
 * - 현장 근로자 마감처리 (당일 이전 날짜 중에서 퇴근을 한 근로자들만 마감처리)
 */
type Init struct {
	WorkerService service.WorkerService
}

func NewInit(safeDb *sqlx.DB) (*Init, error) {
	r := store.Repository{Clocker: clock.RealClock{}}

	init := &Init{
		WorkerService: &service.ServiceWorker{
			SafeDB:  safeDb,
			SafeTDB: safeDb,
			Store:   &r,
		},
	}

	return init, nil
}

func (i *Init) RunInitializations(ctx context.Context) (err error) {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err = i.WorkerService.ModifyWorkerDeadlineInit(ctx); err != nil {
			return fmt.Errorf("[init] RunInitializations fail: %w", err)
		}
		log.Println("[init] ModifyWorkerDeadlineInit completed")
		return nil
	})

	if err = eg.Wait(); err != nil {
		return err
	}
	return
}
