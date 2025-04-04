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

// 현장 위치 주소 추가/수정
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

	query := `
			MERGE INTO IRIS_SITE_POS M
			USING (
				SELECT 
					:1 AS SNO,
					:2 AS ADDRESS_NAME_DEPTH1,
					:3 AS ADDRESS_NAME_DEPTH2,
					:4 AS ADDRESS_NAME_DEPTH3,
					:5 AS ADDRESS_NAME_DEPTH4,
					:6 AS ADDRESS_NAME_DEPTH5,
					:7 AS LATITUDE,
					:8 AS LONGITUDE,
					:9 AS ROAD_ADDRESS_NAME_DEPTH1,
					:10 AS ROAD_ADDRESS_NAME_DEPTH2,
					:11 AS ROAD_ADDRESS_NAME_DEPTH3,
					:12 AS ROAD_ADDRESS_NAME_DEPTH4,
					:13 AS ROAD_ADDRESS_NAME_DEPTH5,
					:14 AS UDF_VAL_01,
					:15 AS UDF_VAL_02,
					:16 AS UDF_VAL_03
				FROM DUAL
			) P
			ON 
				(
				M.SNO = P.SNO
			) WHEN MATCHED THEN
				UPDATE SET
					M.ADDRESS_NAME_DEPTH1 = P.ADDRESS_NAME_DEPTH1,
					M.ADDRESS_NAME_DEPTH2 = P.ADDRESS_NAME_DEPTH2,
					M.ADDRESS_NAME_DEPTH3 = P.ADDRESS_NAME_DEPTH3,
					M.ADDRESS_NAME_DEPTH4 = P.ADDRESS_NAME_DEPTH4,
					M.ADDRESS_NAME_DEPTH5 = P.ADDRESS_NAME_DEPTH5,
					M.LATITUDE = P.LATITUDE,
					M.LONGITUDE = P.LONGITUDE,
					M.ROAD_ADDRESS_NAME_DEPTH1 = P.ROAD_ADDRESS_NAME_DEPTH1,
					M.ROAD_ADDRESS_NAME_DEPTH2 = P.ROAD_ADDRESS_NAME_DEPTH2,
					M.ROAD_ADDRESS_NAME_DEPTH3 = P.ROAD_ADDRESS_NAME_DEPTH3,
					M.ROAD_ADDRESS_NAME_DEPTH4 = P.ROAD_ADDRESS_NAME_DEPTH4,
					M.ROAD_ADDRESS_NAME_DEPTH5 = P.ROAD_ADDRESS_NAME_DEPTH5,
					M.UDF_VAL_01 = P.UDF_VAL_01, -- 도로명 FULL
					M.UDF_VAL_02 = P.UDF_VAL_02, -- 우편번호
					M.UDF_VAL_03 = P.UDF_VAL_03  -- 건물 이름
			WHEN NOT MATCHED THEN
				INSERT (
				    IDX,
					SNO, 
					ADDRESS_NAME_DEPTH1, ADDRESS_NAME_DEPTH2, ADDRESS_NAME_DEPTH3, ADDRESS_NAME_DEPTH4, ADDRESS_NAME_DEPTH5,
					LATITUDE, LONGITUDE,
					ROAD_ADDRESS_NAME_DEPTH1, ROAD_ADDRESS_NAME_DEPTH2, ROAD_ADDRESS_NAME_DEPTH3, ROAD_ADDRESS_NAME_DEPTH4, ROAD_ADDRESS_NAME_DEPTH5,
					UDF_VAL_01, UDF_VAL_02,	UDF_VAL_03
				)
				VALUES (
				    SEQ_IRIS_SITE_POS.NEXTVAL,
					P.SNO, 
					P.ADDRESS_NAME_DEPTH1,
					P.ADDRESS_NAME_DEPTH2,
					P.ADDRESS_NAME_DEPTH3,
					P.ADDRESS_NAME_DEPTH4,
					P.ADDRESS_NAME_DEPTH5,
					P.LATITUDE,
					P.LONGITUDE,
					P.ROAD_ADDRESS_NAME_DEPTH1,
					P.ROAD_ADDRESS_NAME_DEPTH2,
					P.ROAD_ADDRESS_NAME_DEPTH3,
					P.ROAD_ADDRESS_NAME_DEPTH4,
					P.ROAD_ADDRESS_NAME_DEPTH5,
					P.UDF_VAL_01,
					P.UDF_VAL_02,
					P.UDF_VAL_03
				)`

	_, err = tx.ExecContext(ctx, query,
		sno,
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
		sitePosSql.BuildingName)

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
