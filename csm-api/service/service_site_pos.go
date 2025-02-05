package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceSitePos struct {
	DB    store.Queryer
	Store store.SitePosStore
}

// 현장 위치 테이블 조회
//
// @param sno: 현장 고유번호
func (s *ServiceSitePos) GetSitePosData(ctx context.Context, sno int64) (*entity.SitePos, error) {
	sitePosSql, err := s.Store.GetSitePosData(ctx, s.DB, sno)
	if err != nil {
		return nil, fmt.Errorf("service_site_pos/GetSitePosData err: %w", err)
	}
	sitePos := &entity.SitePos{}
	sitePos.ToSitePos(sitePosSql)

	return sitePos, nil
}
