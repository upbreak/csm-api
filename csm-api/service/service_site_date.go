package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceSiteDate struct {
	DB    store.Queryer
	TDB   store.Beginner
	Store store.SiteDateStore
}

// 현장 날짜 테이블 조회
//
// @param sno: 현장 고유번호
func (s *ServiceSiteDate) GetSiteDateData(ctx context.Context, sno int64) (*entity.SiteDate, error) {
	siteDate, err := s.Store.GetSiteDateData(ctx, s.DB, sno)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_site_date/GetSiteDateData err: %w", err)
	}

	return siteDate, nil
}

// 현장 날짜 테이블 수정
//
// @param
// - sno: 현장고유번호
// - siteDate: 현장 시간 (opening_date, closing_plan_date, closing_forecast_date, closing_actual_date)
func (s *ServiceSiteDate) ModifySiteDate(ctx context.Context, sno int64, siteDate entity.SiteDate) error {

	if err := s.Store.ModifySiteDate(ctx, s.TDB, sno, siteDate); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_site_date/ModifySiteDate err: %w", err)
	}

	return nil
}
