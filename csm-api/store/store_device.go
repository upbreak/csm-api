package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) GetDeviceList(ctx context.Context, db Queryer, page entity.PageSql) (*entity.DeviceSqls, error) {
	sqls := entity.DeviceSqls{}

	query := `SELECT *
				FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						SELECT 
							t1.DNO,
							t1.SNO,
							t2.SITE_NM,
							t1.DEVICE_SN,
							t1.DEVICE_NM,
							t1.IS_USE,
							t1.REG_DATE AS REG_DATE,
							t1.MOD_DATE AS MOD_DATE
						FROM IRIS_DEVICE_SET t1
						LEFT OUTER JOIN IRIS_SITE_SET t2 ON t1.SNO = t2.SNO
						WHERE t1.IS_USE = 'Y'
						ORDER BY t1.REG_DATE DESC
					) sorted_data
					WHERE ROWNUM <= :1
				)
				WHERE RNUM > :2`

	if err := db.SelectContext(ctx, &sqls, query, page.EndNum, page.StartNum); err != nil {
		return nil, fmt.Errorf("GetDeviceList err: %v", err)
	}

	return &sqls, nil
}

func (r *Repository) GetDeviceListCount(ctx context.Context, db Queryer) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM IRIS_DEVICE_SET WHERE IS_USE = 'Y'`

	if err := db.GetContext(ctx, &count, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("GetDeviceListCount fail: %w", err)
	}
	return count, nil
}
