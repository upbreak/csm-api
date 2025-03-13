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

// 현장 날씨 테이블 조회
//
// @param sno: 현장 고유번호
func (s *ServiceSiteDate) GetSiteDateData(ctx context.Context, sno int64) (*entity.SiteDate, error) {
	siteDateSql, err := s.Store.GetSiteDateData(ctx, s.DB, sno)
	if err != nil {
		return nil, fmt.Errorf("service_site_date/GetSiteDateData err: %w", err)
	}
	siteDate := &entity.SiteDate{}
	siteDate.ToSiteDate(siteDateSql)

	return siteDate, nil
}

//func ModifySiteDate(ctx context.Context, sno int64, siteDate entity.SiteDate) error
