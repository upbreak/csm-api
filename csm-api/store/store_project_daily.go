package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"fmt"
	"time"
)

func (r *Repository) GetProjectDailyContentList(ctx context.Context, db Queryer, jno int64, targetDate time.Time) (*entity.ProjectDailys, error) {
	projectDailys := entity.ProjectDailys{}

	// jno 변환: 0이면 NULL 처리, 아니면 Valid 값으로 설정
	var jnoParam sql.NullInt64
	if jno != 0 {
		jnoParam = sql.NullInt64{Valid: true, Int64: jno}
	} else {
		jnoParam = sql.NullInt64{Valid: false}
	}

	// targetDate 변환: zero 값이면 NULL 처리, 아니면 Valid 값으로 설정
	var targetDateParam sql.NullTime
	if !targetDate.IsZero() {
		targetDateParam = sql.NullTime{Valid: true, Time: targetDate}
	} else {
		targetDateParam = sql.NullTime{Valid: false}
	}

	sql := `SELECT 
				t1.JNO,
				t1.CONTENT,
				t1.IS_USE,
				t1.REG_DATE,
				t1.MOD_DATE,
				t1.REG_UNO,
				t1.REG_USER,
				t1.MOD_UNO,
				t1.MOD_USER
			FROM
				IRIS_DAILY_JOB t1
			WHERE
				t1.IS_USE = 'Y'
				AND t1.JNO = :2
				AND TO_CHAR(t1.TARGET_DATE, 'YYYY-MM-DD') = TO_CHAR(:2 , 'YYYY-MM-DD')
			ORDER BY
				NVL(t1.REG_DATE, t1.MOD_DATE) DESC`

	if err := db.SelectContext(ctx, &projectDailys, sql, jnoParam, targetDateParam); err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("GetProjectDailyContentList fail: %w", err)
	}
	return &projectDailys, nil
}
