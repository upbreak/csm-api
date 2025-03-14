package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
	"strings"
)

type ServiceSitePos struct {
	DB    store.Queryer
	TDB   store.Beginner
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

	if sitePos.RoadAddress == "" {
		depthArray := []string{sitePos.RoadAddressNameDepth1, sitePos.RoadAddressNameDepth2, sitePos.RoadAddressNameDepth3, sitePos.RoadAddressNameDepth4, sitePos.RoadAddressNameDepth5}
		roadAddress := ""
		for _, depth := range depthArray {
			if depth != "" {
				roadAddress = roadAddress + " " + depth
			}
		}
		sitePos.RoadAddress = strings.Trim(roadAddress, " ")
	}

	return sitePos, nil
}

// 현장 위치 주소 수정
//
// @params
//   - sno : 현장 고유번호
//   - sitePos: 현장 정보 (ADDRESS_NAME_DEPTH1, ADDRESS_NAME_DEPTH2, ADDRESS_NAME_DEPTH3, ADDRESS_NAME_DEPTH4, ADDRESS_NAME_DEPTH5,
//     ROAD_ADDRESS_NAME_DEPTH1, ROAD_ADDRESS_NAME_DEPTH2, ROAD_ADDRESS_NAME_DEPTH3, ROAD_ADDRESS_NAME_DEPTH4, ROAD_ADDRESS_NAME_DEPTH5,
//     ROAD_ADDRESS, ZONE_CODE, BUILDING_NAME)
func (s *ServiceSitePos) ModifySitePos(ctx context.Context, sno int64, sitePos entity.SitePos) error {

	sitePosSql := &entity.SitePosSql{}

	if err := entity.ConvertToSQLNulls(sitePos, sitePosSql); err != nil {
		return fmt.Errorf("service_site_pos/ConvertSliceToSQLNulls err: %w", err)
	}
	if err := s.Store.ModifySitePosData(ctx, s.TDB, sno, *sitePosSql); err != nil {
		return fmt.Errorf("service_site_pos/ModifySitePosData err: %w", err)
	}

	return nil
}
