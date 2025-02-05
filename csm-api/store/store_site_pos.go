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
