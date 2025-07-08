package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

func (r *Repository) SaveWeather(ctx context.Context, tx Execer, weather entity.Weather) error {
	query :=
		`	
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
		// TODO: 에러 아카이브
		return fmt.Errorf("IRIS_WEATHER INSERT failed: %w", err)
	}

	return nil
}
