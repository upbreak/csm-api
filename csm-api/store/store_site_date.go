package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) GetSiteDateData(ctx context.Context, db Queryer, sno int64) (*entity.SiteDateSql, error) {
	siteDateSql := entity.SiteDateSql{}

	query := `SELECT
				t1.OPENING_DATE,
				t1.CLOSING_PLAN_DATE,
				t1.CLOSING_FORECAST_DATE,
				t1.CLOSING_ACTUAL_DATE,
				t1.REG_UNO,
				t1.REG_USER,
				t1.REG_DATE
			FROM
				IRIS_SITE_DATE t1
			WHERE
				t1.SNO = :1
				AND t1.IS_USE = 'Y'`

	if err := db.GetContext(ctx, &siteDateSql, query, sno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &siteDateSql, nil
		}
		return nil, fmt.Errorf("GetSiteDateData fail: %w", err)
	}
	return &siteDateSql, nil
}
