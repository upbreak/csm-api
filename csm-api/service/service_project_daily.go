package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceProjectDaily struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.ProjectDailyStore
}

// 작업내용 조회
func (s *ServiceProjectDaily) GetDailyJobList(ctx context.Context, jno int64, targetDate string) (entity.ProjectDailys, error) {
	list, err := s.Store.GetDailyJobList(ctx, s.SafeDB, jno, targetDate)
	if err != nil {
		return nil, fmt.Errorf("service;GetDailyJobList fail: %w", err)
	}
	return list, nil
}

// 작업내용 추가
func (s *ServiceProjectDaily) AddDailyJob(ctx context.Context, project entity.ProjectDailys) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service;AddDailyJob fail: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service;AddDailyJob panic: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service;AddDailyJob rollback fail: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service;AddDailyJob commit fail: %w", commitErr)
			}
		}
	}()

	if err = s.Store.AddDailyJob(ctx, tx, project); err != nil {
		return fmt.Errorf("service;AddDailyJob fail: %w", err)
	}
	return
}

// 작업내용 수정
func (s *ServiceProjectDaily) ModifyDailyJob(ctx context.Context, project entity.ProjectDaily) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service;ModifyDailyJob fail: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service;ModifyDailyJob panic: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service;ModifyDailyJob rollback fail: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service;ModifyDailyJob commit fail: %w", commitErr)
			}
		}
	}()

	if err = s.Store.ModifyDailyJob(ctx, tx, project); err != nil {
		return fmt.Errorf("service;ModifyDailyJob fail: %w", err)
	}
	return
}

// 작업내용 삭제
func (s *ServiceProjectDaily) RemoveDailyJob(ctx context.Context, idx int64) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service;RemoveDailyJob fail: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service;RemoveDailyJob panic: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service;RemoveDailyJob rollback fail: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
			}
		}
	}()
	if err = s.Store.RemoveDailyJob(ctx, tx, idx); err != nil {
		return fmt.Errorf("service;RemoveDailyJob fail: %w", err)
	}
	return
}
