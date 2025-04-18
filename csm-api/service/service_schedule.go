package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceSchedule struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.ScheduleStore
}

// func: 휴무일 조회
// @param
// -
func (s *ServiceSchedule) GetRestScheduleList(ctx context.Context, jno int64, year string, month string) (entity.RestSchedules, error) {
	list, err := s.Store.GetRestScheduleList(ctx, s.SafeDB, jno, year, month)
	if err != nil {
		return entity.RestSchedules{}, fmt.Errorf("service;GetRestScheduleList: %w", err)
	}

	return list, nil
}

// func: 휴무일 추가
// @param
// -
func (s *ServiceSchedule) AddRestSchedule(ctx context.Context, schedule entity.RestSchedules) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service;add;BeginTx fail: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service;add;Rollback fail: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service;add;Commit fail: %w", commitErr)
			}
		}
	}()

	err = s.Store.AddRestSchedule(ctx, tx, schedule)
	if err != nil {
		return fmt.Errorf("service;Add fail: %w", err)
	}

	return
}

// func: 휴무일 수정
// @param
// -
func (s *ServiceSchedule) ModifyRestSchedule(ctx context.Context, schedule entity.RestSchedule) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service;modify;BeginTx fail: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service;modify;Rollback fail: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service;modify;Commit fail: %w", commitErr)
			}
		}
	}()

	err = s.Store.ModifyRestSchedule(ctx, tx, schedule)
	if err != nil {
		return fmt.Errorf("service;Modify fail: %w", err)
	}

	return
}

// func: 휴무일 삭제
// @param
// -
func (s *ServiceSchedule) RemoveRestSchedule(ctx context.Context, cno int64) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service;remove;BeginTx fail: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service;remove;Rollback fail: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service;remove;Commit fail: %w", commitErr)
			}
		}
	}()

	err = s.Store.RemoveRestSchedule(ctx, tx, cno)
	if err != nil {
		return fmt.Errorf("service;Remove fail: %w", err)
	}
	return
}
