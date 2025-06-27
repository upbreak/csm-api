package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceWorkHour struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.WorkHourStore
}

// 특정 프로젝트 및 근로자의 공수 계산: jno는 필수, ids는 없으면 jno의 모든 근로자 계산 있으면 해당 id의 근로자만 계산
func (s *ServiceWorkHour) ModifyWorkHourByJno(ctx context.Context, jno int64, user entity.Base, ids []string) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_work_hour.ModifyWorkHourByJno BeginTx err: %v", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_work_hour.ModifyWorkHourByJno err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_work_hour.ModifyWorkHourByJno err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	err = s.Store.ModifyWorkHourByJno(ctx, tx, jno, user, ids)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_work_hour.ModifyWorkHourByJno err: %v", err)
	}
	return
}

// 출퇴근이 둘다 있는 모든 근로자의 공수 계산
func (s *ServiceWorkHour) ModifyWorkHour(ctx context.Context, user entity.Base) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_work_hour.ModifyWorkHour BeginTx err: %v", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_work_hour.ModifyWorkHour err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_work_hour.ModifyWorkHour err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	err = s.Store.ModifyWorkHour(ctx, tx, user)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_work_hour.ModifyWorkHour err: %v", err)
	}
	return
}
