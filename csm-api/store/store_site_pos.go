package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) GetSitePosData(ctx context.Context, db Queryer, sno int64) (*entity.SitePosSql, error) {
	sitePosSql := entity.SitePosSql{}

	query := `SELECT
				t1.ADDRESS_NAME_DEPTH1,
				t1.ADDRESS_NAME_DEPTH2,
				t1.ADDRESS_NAME_DEPTH3,
				t1.ADDRESS_NAME_DEPTH4,
				t1.ADDRESS_NAME_DEPTH5,
				t1.ROAD_ADDRESS_NAME_DEPTH1,
				t1.ROAD_ADDRESS_NAME_DEPTH2,
				t1.ROAD_ADDRESS_NAME_DEPTH3,
				t1.ROAD_ADDRESS_NAME_DEPTH4,
				t1.ROAD_ADDRESS_NAME_DEPTH5,
				t1.LATITUDE,
				t1.LONGITUDE,
				t1.REG_DATE
			FROM
				IRIS_SITE_POS t1
			WHERE
				t1.SNO = :1
				AND t1.IS_USE = 'Y'`

	if err := db.GetContext(ctx, &sitePosSql, query, sno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &sitePosSql, nil
		}
		return nil, fmt.Errorf("GetSitePosData fail: %v", err)
	}

	return &sitePosSql, nil
}

// 현장 위치 주소 수정
//
// @params
//   - sno : 현장 고유번호
//   - sitePos: 현장 정보 (ADDRESS_NAME_DEPTH1, ADDRESS_NAME_DEPTH2, ADDRESS_NAME_DEPTH3, ADDRESS_NAME_DEPTH4, ADDRESS_NAME_DEPTH5,
//     LATITUDE, LONGTITUDE,
//     ROAD_ADDRESS_NAME_DEPTH1, ROAD_ADDRESS_NAME_DEPTH2, ROAD_ADDRESS_NAME_DEPTH3, ROAD_ADDRESS_NAME_DEPTH4, ROAD_ADDRESS_NAME_DEPTH5,
//     ROAD_ADDRESS, ZONE_CODE, BUILDING_NAME)
func (r *Repository) ModifySitePosData(ctx context.Context, db Beginner, sno int64, sitePosSql entity.SitePosSql) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store/site_pos. Failed to begin transaction: %v", err)
	}

	query := fmt.Sprintf(`
			UPDATE IRIS_SITE_POS 
			SET
				ADDRESS_NAME_DEPTH1 = :1,
				ADDRESS_NAME_DEPTH2 = :2,
				ADDRESS_NAME_DEPTH3 = :3,
				ADDRESS_NAME_DEPTH4 = :4,
				ADDRESS_NAME_DEPTH5 = :5,
				LATITUDE = :6,
				LONGITUDE = :7,
				ROAD_ADDRESS_NAME_DEPTH1 = :8,
				ROAD_ADDRESS_NAME_DEPTH2 = :9,
				ROAD_ADDRESS_NAME_DEPTH3 = :10,
				ROAD_ADDRESS_NAME_DEPTH4 = :11,
				ROAD_ADDRESS_NAME_DEPTH5 = :12,
				UDF_VAL_01 = :13, -- 도로명 FULL
				UDF_VAL_02 = :14, -- 우편번호
				UDF_VAL_03 = :15  -- 건물 이름
			WHERE
			    SNO = :16
			`)

	_, err = tx.ExecContext(ctx, query,
		sitePosSql.AddressNameDepth1,
		sitePosSql.AddressNameDepth2,
		sitePosSql.AddressNameDepth3,
		sitePosSql.AddressNameDepth4,
		sitePosSql.AddressNameDepth5,
		sitePosSql.Latitude,
		sitePosSql.Longitude,
		sitePosSql.RoadAddressNameDepth1,
		sitePosSql.RoadAddressNameDepth2,
		sitePosSql.RoadAddressNameDepth3,
		sitePosSql.RoadAddressNameDepth4,
		sitePosSql.RoadAddressNameDepth5,
		sitePosSql.RoadAddress,
		sitePosSql.ZoneCode,
		sitePosSql.BuildingName, sno)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("store/site_pos. ModifySitePosData fail: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store/site_pos. Failed to commit transaction: %v", err)
	}

	return nil

}
