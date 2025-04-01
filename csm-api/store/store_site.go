package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
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
				t1.SNO,
				t1.SITE_NM,
				t1.ETC,
				t1.LOC_CODE,
				t1.LOC_NAME,
				t1.IS_USE,
				t1.REG_DATE,
				t1.REG_USER,
				t1.REG_UNO,
				t1.MOD_DATE,
				t1.MOD_USER,
				t1.MOD_UNO,
				t2.JNO AS DEFAULT_JNO,
				t3.JOB_NAME AS DEFAULT_PROJECT_NAME,
				t3.JOB_NO AS DEFAULT_PROJECT_NO,
				CASE 
					WHEN COUNT(CASE WHEN t4.TRANS_TYPE = 'Clock in' THEN 1 END) >= 5 THEN 'Y'
					ELSE 'C'
				END AS CURRENT_SITE_STATS
			FROM
				IRIS_SITE_SET t1
				INNER JOIN IRIS_SITE_JOB t2 
					ON t1.SNO = t2.SNO 
					AND t2.IS_DEFAULT = 'Y'
				INNER JOIN S_JOB_INFO t3 
					ON t2.JNO = t3.JNO
				LEFT JOIN IRIS_RECD_SET t4
					ON t1.SNO = t4.SNO
					AND TO_CHAR(t4.RECOG_TIME, 'YYYY-MM-DD') = TO_CHAR(:1, 'YYYY-MM-DD')
			WHERE
				t1.SNO > 100
			GROUP BY
				t1.SNO,
				t2.JNO,
				t1.SITE_NM,
				t1.ETC,
				t1.LOC_CODE,
				t1.LOC_NAME,
				t1.IS_USE,
				t1.REG_DATE,
				t1.REG_USER,
				t1.REG_UNO,
				t1.MOD_DATE,
				t1.MOD_USER,
				t1.MOD_UNO,
				t3.JOB_NAME,
				t3.JOB_NO
			ORDER BY
				t1.SNO DESC`

	if err := db.SelectContext(ctx, &siteSqls, sql, targetDate); err != nil {
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

// func: 현장 상태 조회
// @param
// -
func (r *Repository) GetSiteStatsList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.SiteSqls, error) {
	siteSqls := entity.SiteSqls{}

	query := `
				SELECT 
					t1.SNO,
					NVL(t2.CURRENT_SITE_STATS, 'C') CURRENT_SITE_STATS
				FROM IRIS_RECD_SET t1
				LEFT JOIN (
					SELECT SNO, 
						CASE 
							WHEN COUNT(CASE WHEN TRANS_TYPE = 'Clock in' THEN 1 END) >= 5 THEN 'Y'
							ELSE 'C'
						END AS CURRENT_SITE_STATS
					FROM IRIS_RECD_SET 
					WHERE SNO > 100 
					AND TO_CHAR(RECOG_TIME, 'YYYY-MM-DD') = TO_CHAR(:1, 'YYYY-MM-DD')
					GROUP by SNO
				) t2 ON t1.SNO = t2.SNO
				WHERE t1.SNO > 100
				GROUP by t1.SNO, t2.CURRENT_SITE_STATS`
	if err := db.SelectContext(ctx, &siteSqls, query, targetDate); err != nil {
		return &siteSqls, fmt.Errorf("getSiteStatsList fail: %w", err)
	}
	return &siteSqls, nil
}

// func: 현장 수정
// @param
// -
func (r *Repository) ModifySite(ctx context.Context, db Beginner, site entity.Site) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("store/site. Failed to begin transaction: %w", err)
	}

	agent := utils.GetAgent()

	query := `
			UPDATE IRIS_SITE_SET 
			SET
			    SITE_NM = :1,
			    ETC = :2,
				MOD_UNO = :3,
				MOD_USER = :4,
				MOD_AGENT = :5,
				MOD_DATE = SYSDATE
			WHERE
			    SNO = :6
			`
	if _, err := tx.ExecContext(ctx, query, site.SiteNm, site.Etc, site.ModUno, site.ModUser, agent, site.Sno); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("store/site. ModifySite fail: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store/site. Failed to commit transaction: %w", err)
	}

	return nil
}

// func: 현장 생성
// @param
// -
func (r *Repository) AddSite(ctx context.Context, db Queryer, tdb Beginner, jno int64, user entity.User) error {
	var generatedSNO int64

	// sno 생성
	query := `SELECT SEQ_IRIS_SITE_SET.NEXTVAL FROM DUAL`
	if err := db.GetContext(ctx, &generatedSNO, query); err != nil {
		return fmt.Errorf("store/site. Failed to get generated SITE_SET_SEQ.NEXTVAL: %w", err)
	}

	tx, err := tdb.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddSite. Failed to begin transaction: %w", err)
	}

	// IRIS_SITE_SET 생성
	query = `
			INSERT INTO IRIS_SITE_SET(
				SNO, SITE_NM, LOC_CODE, LOC_NAME, IS_USE, 
			    REG_DATE, REG_AGENT, REG_USER, REG_UNO
			) 
			SELECT 
				:1, JOB_NAME, JOB_LOC, JOB_LOC_NAME, 'Y', 
				SYSDATE, :2, :3, :4
			FROM s_job_info 
			WHERE JNO = :5`
	if _, err = tx.ExecContext(ctx, query, generatedSNO, user.Agent, user.UserName, user.Uno, jno); err != nil {
		origErr := err
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("IRIS_SITE_SET. Failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("IRIS_SITE_SET INSERT failed: %w", origErr)
	}

	// IRIS_SITE_JOB 생성
	query = `
			INSERT INTO IRIS_SITE_JOB(
				SNO, JNO, IS_USE, IS_DEFAULT, REG_DATE,
				REG_AGENT, REG_USER, REG_UNO
			) VALUES (
				:1, :2, 'Y', 'Y', SYSDATE,
				:3, :4, :5
			)`
	if _, err = tx.ExecContext(ctx, query, generatedSNO, jno, user.Agent, user.UserName, user.Uno); err != nil {
		origErr := err
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("IRIS_SITE_JOB. Failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("IRIS_SITE_JOB INSERT failed: %w", origErr)
	}

	// IRIS_SITE_DATE 생성
	query = `
			INSERT INTO IRIS_SITE_DATE(
				SNO, OPENING_DATE, CLOSING_PLAN_DATE, IS_USE, REG_DATE,
				REG_AGENT, REG_USER, REG_UNO
			)
			SELECT
				:1,	TO_DATE(JOB_SD, 'YYYY-MM-DD'), TO_DATE(JOB_ED, 'YYYY-MM-DD'), 'Y', SYSDATE,
				:2, :3, :4
			FROM s_job_info
			WHERE JNO = :5`
	if _, err = tx.ExecContext(ctx, query, generatedSNO, user.Agent, user.UserName, user.Uno, jno); err != nil {
		origErr := err
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("IRIS_SITE_DATE. Failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("IRIS_SITE_DATE INSERT failed: %w", origErr)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("AddSite. Failed to commit transaction: %w", err)
	}

	return nil
}
