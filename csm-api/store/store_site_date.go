package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) GetSiteDateData(ctx context.Context, db Queryer, sno int64) (*entity.SiteDate, error) {
	siteDate := entity.SiteDate{}

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

	if err := db.GetContext(ctx, &siteDate, query, sno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &siteDate, nil
		}
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("GetSiteDateData fail: %w", err)
	}
	return &siteDate, nil
}

// 현장 날짜 테이블 수정
//
// @param
// - sno: 현장고유번호
// - siteDate: 현장 시간 (opening_date, closing_plan_date, closing_forecast_date, closing_actual_date)
func (r *Repository) ModifySiteDate(ctx context.Context, db Beginner, sno int64, siteDateSql entity.SiteDate) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("store/site_date. Failed to begin transaction: %v", err)
	}

	query := fmt.Sprintf(`
			UPDATE IRIS_SITE_DATE 
			SET
			    OPENING_DATE = :1,
				CLOSING_PLAN_DATE = :2,
				CLOSING_FORECAST_DATE = :3,
				CLOSING_ACTUAL_DATE = :4
			WHERE
				SNO = :5 
				AND IS_USE = 'Y'
			`)

	_, err = tx.ExecContext(ctx, query, siteDateSql.OpeningDate, siteDateSql.ClosingPlanDate, siteDateSql.ClosingForecastDate, siteDateSql.ClosingActualDate, sno)
	if err != nil {
		origErr := err
		if err = tx.Rollback(); err != nil {
			return err
		}
		//TODO: 에러 아카이브
		return fmt.Errorf("store/site_date. ModifySiteDate fail: %v", origErr)
	}

	if err = tx.Commit(); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("store/site_date. failed to commit transaction: %v", err)
	}

	return nil
}
