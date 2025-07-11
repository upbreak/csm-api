package store

import (
	"context"
	"csm-api/entity"
	"fmt"
	"time"
)

// func: weather 저장
func (r *Repository) SaveWeather(ctx context.Context, tx Execer, weather entity.Weather) error {
	query := `	
		INSERT INTO IRIS_WEATHER (
		    SNO, LGT, PTY, RN1, SKY, T1H, 
		    REH, UUU, VVV, VEC, WSD, RECOG_TIME
		)
		VALUES (
		    :1, :2, :3, :4, :5, :6, 
		    :7, :8, :9, :10, :11, :12
		)
		`

	if _, err := tx.ExecContext(ctx, query, weather.Sno, weather.Lgt, weather.Pty, weather.Rn1, weather.Sky, weather.T1h, weather.Reh, weather.Uuu, weather.Vvv, weather.Vec, weather.Wsd, weather.RecogTime); err != nil {
		return fmt.Errorf("IRIS_WEATHER INSERT failed: %w", err)
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
		return nil, fmt.Errorf("IRIS_WEATHER LIST failed: %w", err)
	}

	return &weathers, nil
}
