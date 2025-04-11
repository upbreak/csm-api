package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
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
	sitePos, err := s.Store.GetSitePosData(ctx, s.DB, sno)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_site_pos/GetSitePosData err: %w", err)
	}

	if sitePos.RoadAddress.String == "" {
		depthArray := []string{sitePos.RoadAddressNameDepth1.String, sitePos.RoadAddressNameDepth2.String, sitePos.RoadAddressNameDepth3.String, sitePos.RoadAddressNameDepth4.String, sitePos.RoadAddressNameDepth5.String}
		roadAddress := ""
		for _, depth := range depthArray {
			if depth != "" {
				roadAddress = roadAddress + " " + depth
			}
		}
		sitePos.RoadAddress = utils.ParseNullString(strings.Trim(roadAddress, " "))
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
func (s *ServiceSitePos) ModifySitePos(ctx context.Context, sno int64, sitePos entity.SitePos) (err error) {
	tx, err := s.TDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_site_pos/TBeginTx err: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_site_pos/TRollback err: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_site_pos/TXCommit err: %w", commitErr)
			}
		}
	}()

	if err := s.Store.ModifySitePosData(ctx, tx, sno, sitePos); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_site_pos/ModifySitePosData err: %w", err)
	}

	return nil
}
