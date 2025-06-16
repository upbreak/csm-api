package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceManHour struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.ManHourStore
}

// func: 프로젝트에 설정된 공수 조회
// @param
// - jno: 프로젝트pk
func (s *ServiceManHour) GetManHourList(ctx context.Context, jno int64) (*entity.ManHours, error) {

	manhours, err := s.Store.GetManHourList(ctx, s.SafeDB, jno)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.ManHours{}, fmt.Errorf("service_manHour/GetManHourList err: %w", err)
	}

	return manhours, nil

}

// func: 공수 수정 및 추가
// @param
// - manHours: 공수 정보 배열
func (s *ServiceManHour) MergeManHour(ctx context.Context, manHour entity.ManHour) error {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_project/ModifyProjectSetting Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_project/ModifyProjectSetting Commit error: %w", commitErr)
			}
		}
	}()

	if err = s.Store.MergeManHour(ctx, tx, manHour); err != nil {
		return fmt.Errorf("service_project/CheckProjectSetting error: %w", err)
	}

	return nil
}
