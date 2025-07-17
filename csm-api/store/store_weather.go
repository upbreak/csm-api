package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"time"
)

// func: weather 저장
func (r *Repository) SaveWeather(ctx context.Context, tx Execer, weather entity.Weather) error {
	if weather.RecogTime.Valid {
		weather.RecogTime.Time = weather.RecogTime.Time.Truncate(time.Hour)
	}

	query := `
		MERGE INTO IRIS_WEATHER target
		USING (
			SELECT :1 AS SNO, :2 AS RECOG_TIME FROM DUAL
		) source
		ON (
			target.SNO = source.SNO AND
			TO_CHAR(target.RECOG_TIME, 'YYYY-MM-DD HH24') = TO_CHAR(source.RECOG_TIME, 'YYYY-MM-DD HH24')
		)
		WHEN NOT MATCHED THEN
		INSERT (
			SNO, LGT, PTY, RN1, SKY, T1H, 
			REH, UUU, VVV, VEC, WSD, RECOG_TIME
		) VALUES (
			:3, :4, :5, :6, :7, :8, 
			:9, :10, :11, :12, :13, :14
		)
	`

	// 3. 실행
	if _, err := tx.ExecContext(ctx, query,
		weather.Sno,       // :1
		weather.RecogTime, // :2
		weather.Sno,       // :3
		weather.Lgt,       // :4
		weather.Pty,       // :5
		weather.Rn1,       // :6
		weather.Sky,       // :7
		weather.T1h,       // :8
		weather.Reh,       // :9
		weather.Uuu,       // :10
		weather.Vvv,       // :11
		weather.Vec,       // :12
		weather.Wsd,       // :13
		weather.RecogTime, // :14
	); err != nil {
		return utils.CustomErrorf(err)
	}

	return nil
}

// func: 저장된 날씨 리스트 조회
// params
// - sno: 현장 PK
// - targetDate: 조회할 날짜
func (r *Repository) GetWeatherList(ctx context.Context, db Queryer, sno int64, targetDate time.Time) (*entity.Weathers, error) {
	weathers := entity.Weathers{}
	query := `
			SELECT *
			FROM IRIS_WEATHER
			WHERE SNO = :1 
			AND TRUNC(RECOG_TIME) = TO_DATE(TO_CHAR(:2, 'YYYY-MM-DD'), 'YYYY-MM-DD')
			ORDER BY RECOG_TIME ASC
		`

	if err := db.SelectContext(ctx, &weathers, query, sno, targetDate); err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return &weathers, nil
}
