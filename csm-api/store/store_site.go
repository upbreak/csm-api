package store

import (
	"context"
	"csm-api/entity"
	"fmt"
	"time"
)

func (r *Repository) GetSiteList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.SiteSqls, error) {
	siteSqls := entity.SiteSqls{}

	sql := `SELECT
				*
			FROM (
				WITH sour AS (
					SELECT
						t1.SNO
						,t1.SITE_NM
						,t1.ETC
						,t1.LOC_CODE
						,t1.LOC_NAME
						,t1.IS_USE
						,t1.REG_DATE
						,t1.REG_USER
						,t1.REG_UNO
						,t1.MOD_DATE
						,t1.MOD_USER
						,t1.MOD_UNO
						,t2.JNO AS DEFAULT_JNO
						,t3.JOB_NAME AS DEFAULT_PROJECT_NAME
						,t3.JOB_NO AS DEFAULT_PROJECT_NO
					FROM
						IRIS_SITE_SET t1
						INNER JOIN IRIS_SITE_JOB t2 ON t1.SNO = t2.SNO AND t2.IS_DEFAULT = 'Y'
						INNER JOIN S_JOB_INFO t3 ON t2.JNO = t3.JNO
					WHERE
						t1.sno > 100
					ORDER BY
						t1.SNO DESC
				)
				SELECT
					sour.*,
					NVL(iss.STATS, 'Y') AS CURRENT_SITE_STATS
				FROM
					sour
				LEFT JOIN IRIS_SITE_STATS iss
					ON sour.SNO = iss.SNO
					AND iss.START_DATE <= :1
					AND iss.END_DATE >= :2
					AND iss.IS_USE = 'Y'
			) A
			WHERE
				1 = 1`

	if err := db.SelectContext(ctx, &siteSqls, sql, targetDate, targetDate); err != nil {
		return &siteSqls, fmt.Errorf("getSiteList fail: %w", err)
	}

	return &siteSqls, nil
}
