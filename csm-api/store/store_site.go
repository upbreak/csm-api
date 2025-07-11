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
func (r *Repository) GetSiteList(ctx context.Context, db Queryer, targetDate time.Time, role int, uno int64) (*entity.Sites, error) {
	sites := entity.Sites{}

	sql := `
			WITH USER_IN_JNO AS (
					SELECT JNO 
					FROM S_JOB_MEMBER_LIST
					WHERE 1 = :1 OR UNO = :2
				UNION
					SELECT JNO
					FROM JOB_SUBCON_INFO 
					WHERE ID = :3
			)
			SELECT
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
						WHEN EXISTS (
							SELECT 1
							FROM IRIS_SCH_REST_SET r
							WHERE r.JNO = t2.JNO
							  AND TO_DATE(r.REST_YEAR || LPAD(r.REST_MONTH, 2, '0') || LPAD(r.REST_DAY, 2, '0'), 'YYYYMMDD') = TRUNC(:4)
						) THEN 'H'
						WHEN (
							SELECT COUNT(*)
							FROM IRIS_WORKER_DAILY_SET d
							WHERE d.SNO = t1.SNO
							  AND TRUNC(d.RECORD_DATE) = TRUNC(:5)
							  AND d.WORK_STATE = '01'
						) >= 5 THEN 'Y'
						ELSE 'C'
					END AS CURRENT_SITE_STATS
				FROM IRIS_SITE_SET t1
				INNER JOIN IRIS_SITE_JOB t2 ON t1.SNO = t2.SNO AND t2.IS_DEFAULT = 'Y'
				INNER JOIN S_JOB_INFO t3 ON t2.JNO = t3.JNO AND t3.JNO IN (SELECT * FROM USER_IN_JNO)
				INNER JOIN (SELECT * FROM IRIS_SITE_DATE WHERE (:6 BETWEEN OPENING_DATE AND CLOSING_ACTUAL_DATE) OR (:7 >= OPENING_DATE AND CLOSING_ACTUAL_DATE IS NULL) OR (:8 <= CLOSING_ACTUAL_DATE AND OPENING_DATE IS NULL)) t4 ON t1.SNO = t4.SNO
				WHERE t1.SNO > -1
				AND t1.IS_USE = 'Y'
				ORDER BY t1.SNO DESC`


	if err := db.SelectContext(ctx, &sites, sql, role, uno, uno, targetDate, targetDate, targetDate, targetDate, targetDate); err != nil {
		//TODO: 에러 아카이브

		return &sites, fmt.Errorf("getSiteList fail: %w", err)
	}

	return &sites, nil
}

// func: 현장 데이터 리스트
// @param
// -
func (r *Repository) GetSiteNmList(ctx context.Context, db Queryer, page entity.PageSql, search entity.Site, nonSite int) (*entity.Sites, error) {

	condition := ""

	condition = utils.Int64WhereConvert(condition, search.Sno.NullInt64, "t1.SNO")
	condition = utils.StringWhereConvert(condition, search.SiteNm.NullString, "t1.SITE_NM")
	condition = utils.StringWhereConvert(condition, search.Etc.NullString, "t1.ETC")
	condition = utils.StringWhereConvert(condition, search.LocName.NullString, "t1.LOC_NAME")

	sites := entity.Sites{}

	order := ""
	if page.Order.Valid {
		order = page.Order.String
	} else {
		order = "''"
	}
	query := fmt.Sprintf(`
				SELECT * FROM (
				    SELECT ROWNUM AS RNUM, sorted_data.*
				    FROM (				        
						SELECT 
							t1.SNO,
							t1.SITE_NM,
							t1.LOC_CODE,
							t1.LOC_NAME,
							t1.ETC,
							t1.REG_DATE,
							t1.MOD_DATE
						FROM IRIS_SITE_SET t1
						WHERE (sno > 100
							OR ( 1=:1 AND sno = 100))
							%s
				    ) sorted_data
					WHERE ROWNUM <= :2
					ORDER BY 
					    CASE WHEN 
							SNO = 100 
							THEN 0 
							ELSE 1 
						END,
					    %s,
						SNO DESC
				) WHERE RNUM > :3`, condition, order)

	if err := db.SelectContext(ctx, &sites, query, nonSite, page.EndNum, page.StartNum); err != nil {
		return &sites, fmt.Errorf("getSiteNmList fail: %w", err)
	}
	return &sites, nil
}

// func: 현장 데이터 개수
// @param
// -
func (r *Repository) GetSiteNmCount(ctx context.Context, db Queryer, search entity.Site, nonSite int) (int, error) {
	var count int

	condition := ""

	condition = utils.Int64WhereConvert(condition, search.Sno.NullInt64, "t1.SNO")
	condition = utils.StringWhereConvert(condition, search.SiteNm.NullString, "t1.SITE_NM")
	condition = utils.StringWhereConvert(condition, search.Etc.NullString, "t1.ETC")
	condition = utils.StringWhereConvert(condition, search.LocName.NullString, "t1.LOC_NAME")

	query := fmt.Sprintf(`			        
						SELECT 
							count(*)
						FROM IRIS_SITE_SET t1
						WHERE (sno > 100
							OR ( 1= :1 AND sno = 100))
							%s
				    `, condition)

	if err := db.GetContext(ctx, &count, query, nonSite); err != nil {
		return 0, fmt.Errorf("getSiteNmList fail: %w", err)
	}
	return count, nil
}

// func: 현장 상태 조회
// @param
// -
func (r *Repository) GetSiteStatsList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.Sites, error) {
	sites := entity.Sites{}

	query := `
				SELECT DISTINCT 
					T1.SNO,
					CASE 
						WHEN T2.JNO IS NOT NULL THEN 'H'
						WHEN NVL(T3.WORKER_COUNT, 0) >= 5 THEN 'Y'
						ELSE 'C'
					END AS CURRENT_SITE_STATS
				FROM IRIS_SITE_JOB T1
				LEFT JOIN (
					SELECT DISTINCT JNO
					FROM IRIS_SCH_REST_SET
					WHERE TO_DATE(REST_YEAR || LPAD(REST_MONTH, 2, '0') || LPAD(REST_DAY, 2, '0'), 'YYYYMMDD') = TRUNC(:1)
				) T2 ON T1.JNO = T2.JNO
				LEFT JOIN (
					SELECT SNO, COUNT(*) AS WORKER_COUNT
					FROM IRIS_WORKER_DAILY_SET
					WHERE TRUNC(RECORD_DATE) = TRUNC(:2)
					AND WORK_STATE = '01'
					GROUP BY SNO
				) T3 ON T1.SNO = T3.SNO`
	if err := db.SelectContext(ctx, &sites, query, targetDate, targetDate); err != nil {
		return &sites, fmt.Errorf("getSiteStatsList fail: %w", err)
	}
	return &sites, nil
}

// func: 현장 수정
// @param
// -
func (r *Repository) ModifySite(ctx context.Context, tx Execer, site entity.Site) error {
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
		return fmt.Errorf("store/site. ModifySite fail: %v", err)
	}

	return nil
}

// func: 현장 생성
// @param
// -
func (r *Repository) AddSite(ctx context.Context, db Queryer, tx Execer, jno int64, user entity.User) error {
	var generatedSNO int64

	// sno 생성
	query := `SELECT SEQ_IRIS_SITE_SET.NEXTVAL FROM DUAL`
	if err := db.GetContext(ctx, &generatedSNO, query); err != nil {
		return fmt.Errorf("store/site. Failed to get generated SITE_SET_SEQ.NEXTVAL: %w", err)
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
	if _, err := tx.ExecContext(ctx, query, generatedSNO, user.Agent, user.UserName, user.Uno, jno); err != nil {
		return fmt.Errorf("IRIS_SITE_SET INSERT failed: %w", err)
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
	if _, err := tx.ExecContext(ctx, query, generatedSNO, jno, user.Agent, user.UserName, user.Uno); err != nil {
		return fmt.Errorf("IRIS_SITE_JOB INSERT failed: %w", err)
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
	if _, err := tx.ExecContext(ctx, query, generatedSNO, user.Agent, user.UserName, user.Uno, jno); err != nil {
		return fmt.Errorf("IRIS_SITE_DATE INSERT failed: %w", err)
	}

	return nil
}

// func: 현장 사용안함 변경
// @param
// -
func (r *Repository) ModifySiteIsNonUse(ctx context.Context, tx Execer, site entity.ReqSite) error {
	agent := utils.GetAgent()

	query := `
			UPDATE IRIS_SITE_SET
			SET 
			    IS_USE = 'N',
				MOD_AGENT = :1,
				MOD_USER = :2,
				MOD_UNO = :3,
				MOD_DATE = SYSDATE
			WHERE SNO = :4`
	if _, err := tx.ExecContext(ctx, query, agent, site.ModUser, site.ModUno, site.Sno); err != nil {
		return fmt.Errorf("store/site. ModifySiteIsNonUse fail: %w", err)
	}

	return nil
}
