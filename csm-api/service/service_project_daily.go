package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
	"time"
)

type ServiceProjectDaily struct {
	DB    store.Queryer
	Store store.ProjectDailyStore
}

// 현장 고유번호로 현장에 해당하는 프로젝트 리스트 조회 비즈니스
//
// @param jno: 프로젝트 관리번호
// @param targetDate: 현재 시간
func (s *ServiceProjectDaily) GetProjectDailyContentList(ctx context.Context, jno int64, targetDate time.Time) (*entity.ProjectDailys, error) {
	projectDailys, err := s.Store.GetProjectDailyContentList(ctx, s.DB, jno, targetDate)
	if err != nil {
		return nil, fmt.Errorf("service_project_daily/GetProjectDailyContentList err: %w", err)
	}
	return projectDailys, nil
}
