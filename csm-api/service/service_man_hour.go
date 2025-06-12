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

	// jno 기본 공수 넣기
	if err = s.Store.MergeManHour(ctx, tx, manHour); err != nil {
		return fmt.Errorf("service_project/CheckProjectSetting error: %w", err)
	}

	return nil
}
