package store

import (
	"context"
	"csm-api/entity"
	"fmt"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// func: 현장 관리 조회
// @param
// - targetDate: 현재시간
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

// func: 현장 데이터 리스트
// @param
// -
func (r *Repository) GetSiteNmList(ctx context.Context, db Queryer) (*entity.SiteSqls, error) {
	siteSqls := entity.SiteSqls{}

	query := `
				SELECT 
					t1.SNO,
					t1.SITE_NM,
					t1.LOC_CODE,
					t1.LOC_NAME,
					t1.ETC,
					t1.REG_DATE,
					t1.MOD_DATE
				FROM IRIS_SITE_SET t1
				WHERE sno > 100`
	//WHERE t1.IS_USE ='Y'`

	if err := db.SelectContext(ctx, &siteSqls, query); err != nil {
		return &siteSqls, fmt.Errorf("getSiteNmList fail: %w", err)
	}
	return &siteSqls, nil
}
