package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"database/sql"
	"fmt"
	"time"
)

// func: 현장 프로젝트 조회
// @param
// - sno int64 현장 번호, targetDate time.Time: 현재시간
func (r *Repository) GetProjectList(ctx context.Context, db Queryer, sno int64, targetDate time.Time) (*entity.ProjectInfoSqls, error) {
	projectInfoSqls := entity.ProjectInfoSqls{}

	// sno 변환: 0이면 NULL 처리, 아니면 Valid 값으로 설정
	var snoParam sql.NullInt64
	if sno != 0 {
		snoParam = sql.NullInt64{Valid: true, Int64: sno}
	} else {
		snoParam = sql.NullInt64{Valid: false}
	}

	sql := `
			WITH base AS (
				SELECT
					t1.SNO,
					t1.JNO,
					t1.USER_ID,
					NVL(t2.USER_NM, ' ') AS USER_NM,
					TO_CHAR(t1.RECORD_DATE, 'YYYY-MM-DD') AS RECORD_DATE,
					NVL(t2.DEPARTMENT, ' ') AS DEPARTMENT
				FROM IRIS_WORKER_DAILY_SET t1
				LEFT JOIN IRIS_WORKER_SET t2 ON t1.SNO = t2.SNO AND t1.JNO = t2.JNO AND t1.USER_ID = t2.USER_ID
				WHERE t1.SNO > 100
				AND (:1 IS NULL OR t1.SNO = :2)
			),
			worker_counts AS (
				SELECT
					SNO,
					JNO,
					COUNT(DISTINCT USER_ID || '-' || RECORD_DATE) AS WORKER_COUNT_ALL,
					COUNT(DISTINCT CASE 
									 WHEN RECORD_DATE = TO_CHAR(:3, 'YYYY-MM-DD')
									 THEN USER_ID || '-' || RECORD_DATE
									 END) AS WORKER_COUNT_DATE,
					COUNT(DISTINCT CASE 
									 WHEN RECORD_DATE = TO_CHAR(:4, 'YYYY-MM-DD')
									  AND (INSTR(DEPARTMENT, '하이테크') > 0 
									   OR INSTR(DEPARTMENT, 'HTENC') > 0 
									   OR INSTR(DEPARTMENT, 'HTE') > 0)
									 THEN USER_ID || '-' || RECORD_DATE
									 END) AS WORKER_COUNT_HTENC,
					COUNT(DISTINCT CASE 
									 WHEN RECORD_DATE = TO_CHAR(:5, 'YYYY-MM-DD')
									  AND INSTR(DEPARTMENT, '하이테크') = 0
									  AND INSTR(DEPARTMENT, 'HTENC') = 0
									  AND INSTR(DEPARTMENT, 'HTE') = 0
									  AND (INSTR(DEPARTMENT, '관리') > 0 
										   OR INSTR(USER_NM, '관리') > 0)
									 THEN USER_ID || '-' || RECORD_DATE
									 END) AS WORKER_COUNT_MANAGER,
					COUNT(DISTINCT CASE 
									 WHEN RECORD_DATE = TO_CHAR(:6, 'YYYY-MM-DD')
									  AND INSTR(DEPARTMENT, '하이테크') = 0
									  AND INSTR(DEPARTMENT, 'HTENC') = 0
									  AND INSTR(DEPARTMENT, 'HTE') = 0
									  AND (INSTR(DEPARTMENT, '관리') = 0 
										   AND INSTR(USER_NM, '관리') = 0)
									 THEN USER_ID || '-' || RECORD_DATE
									 END) AS WORKER_COUNT_NOT_MANAGER
				FROM base
				GROUP BY SNO, JNO
			)
			SELECT
				t1.SNO,
				t1.JNO,
				t1.IS_USE,
				t1.IS_DEFAULT,
				t1.REG_DATE,
				t1.REG_USER,
				t1.REG_UNO,
				t1.MOD_DATE,
				t1.MOD_USER,
				t1.MOD_UNO,
				t2.JOB_TYPE AS PROJECT_TYPE,
				t6.CD_NM AS PROJECT_TYPE_NM,
				t2.JOB_NO AS PROJECT_NO,
				t2.JOB_NAME AS PROJECT_NM,
				t2.JOB_YEAR AS PROJECT_YEAR,
				t2.JOB_LOC AS PROJECT_LOC,
				t2.JOB_CODE AS PROJECT_CODE,
				t4.KIND_NAME AS PROJECT_CODE_NAME,
				t3.SITE_NM,
				t2.COMP_CODE,
				t2.COMP_NICK,
				t2.COMP_NAME,
				t2.COMP_ETC,
				t2.ORDER_COMP_CODE,
				t2.ORDER_COMP_NICK,
				t2.ORDER_COMP_NAME,
				t2.ORDER_COMP_JOB_NAME,
				t2.JOB_LOC_NAME AS PROJECT_LOC_NAME,
				t2.JOB_PM,
				t2.JOB_PM_NAME,
				t2.JOB_PE,
				TO_DATE(t2.JOB_SD, 'YYYY-MM-DD') AS PROJECT_STDT,
				TO_DATE(t2.JOB_ED, 'YYYY-MM-DD') AS PROJECT_EDDT,
				TO_DATE(t2.REG_DATE, 'YYYY-MM-DD HH24:MI:SS') AS PROJECT_REG_DATE,
				TO_DATE(t2.MOD_DATE, 'YYYY-MM-DD HH24:MI:SS') AS PROJECT_MOD_DATE,
				t2.JOB_STATE AS PROJECT_STATE,
				t5.CD_NM AS PROJECT_STATE_NM,
				t2.MOC_NO,
				t2.WO_NO,
				NVL(wc.WORKER_COUNT_ALL, 0) AS WORKER_COUNT_ALL,
				NVL(wc.WORKER_COUNT_DATE, 0) AS WORKER_COUNT_DATE,
				NVL(wc.WORKER_COUNT_HTENC, 0) AS WORKER_COUNT_HTENC,
				NVL(wc.WORKER_COUNT_MANAGER, 0) AS WORKER_COUNT_MANAGER,
				NVL(wc.WORKER_COUNT_NOT_MANAGER, 0) AS WORKER_COUNT_NOT_MANAGER
			FROM IRIS_SITE_JOB t1
			INNER JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
			INNER JOIN IRIS_SITE_SET t3 ON t1.SNO = t3.SNO
			INNER JOIN TIMESHEET.JOB_KIND_CODE t4 ON t2.JOB_CODE = t4.KIND_CODE
			INNER JOIN TIMESHEET.SYS_CODE_SET t5 ON t5.MINOR_CD = t2.job_state 
				 AND t5.major_cd = 'JOB_STATE'
			INNER JOIN TIMESHEET.SYS_CODE_SET t6 ON t6.MINOR_CD = t2.JOB_TYPE 
				 AND t6.major_cd = 'JOB_TYPE' 
			LEFT JOIN worker_counts wc ON t1.SNO = wc.SNO 
				 AND t1.JNO = wc.JNO
			WHERE t1.SNO > 100
			AND (:7 IS NULL OR t1.SNO = :8)
			ORDER BY IS_DEFAULT DESC`
	if err := db.SelectContext(ctx, &projectInfoSqls, sql, snoParam, snoParam, targetDate, targetDate, targetDate, targetDate, snoParam, snoParam); err != nil {
		return &projectInfoSqls, fmt.Errorf("getProjectList fail: %v", err)
	}

	return &projectInfoSqls, nil
}

// func: 프로젝트 근로자 수 조회
// @param
// - sno int64 현장 번호, targetDate time.Time: 현재시간
func (r *Repository) GetProjectWorkerCountList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.ProjectInfoSqls, error) {
	sqlList := entity.ProjectInfoSqls{}

	query := `
				WITH base AS (
					SELECT
						t1.SNO,
						t1.JNO,
						t1.USER_ID,
						NVL(t2.USER_NM, ' ') AS USER_NM,
						TO_CHAR(t1.RECORD_DATE, 'YYYY-MM-DD') AS RECORD_DATE,
						NVL(t2.DEPARTMENT, ' ') AS DEPARTMENT
					FROM IRIS_WORKER_DAILY_SET t1
					LEFT JOIN IRIS_WORKER_SET t2 ON t1.SNO = t2.SNO AND t1.JNO = t2.JNO AND t1.USER_ID = t2.USER_ID
					WHERE t1.SNO > 100
				),
				worker_counts AS (
					SELECT
						SNO,
						JNO,
						COUNT(DISTINCT USER_ID || '-' || RECORD_DATE) AS WORKER_COUNT_ALL,
						COUNT(DISTINCT CASE 
										 WHEN RECORD_DATE = TO_CHAR(:1, 'YYYY-MM-DD')
										 THEN USER_ID || '-' || RECORD_DATE
										 END) AS WORKER_COUNT_DATE,
						COUNT(DISTINCT CASE 
										 WHEN RECORD_DATE = TO_CHAR(:2, 'YYYY-MM-DD')
										  AND (INSTR(DEPARTMENT, '하이테크') > 0 
										   OR INSTR(DEPARTMENT, 'HTENC') > 0 
										   OR INSTR(DEPARTMENT, 'HTE') > 0)
										 THEN USER_ID || '-' || RECORD_DATE
										 END) AS WORKER_COUNT_HTENC,
						COUNT(DISTINCT CASE 
										 WHEN RECORD_DATE = TO_CHAR(:3, 'YYYY-MM-DD')
										  AND INSTR(DEPARTMENT, '하이테크') = 0
										  AND INSTR(DEPARTMENT, 'HTENC') = 0
										  AND INSTR(DEPARTMENT, 'HTE') = 0
										  AND (INSTR(DEPARTMENT, '관리') > 0 
											   OR INSTR(USER_NM, '관리') > 0)
										 THEN USER_ID || '-' || RECORD_DATE
										 END) AS WORKER_COUNT_MANAGER,
						COUNT(DISTINCT CASE 
										 WHEN RECORD_DATE = TO_CHAR(:4, 'YYYY-MM-DD')
										  AND INSTR(DEPARTMENT, '하이테크') = 0
										  AND INSTR(DEPARTMENT, 'HTENC') = 0
										  AND INSTR(DEPARTMENT, 'HTE') = 0
										  AND (INSTR(DEPARTMENT, '관리') = 0 
											   AND INSTR(USER_NM, '관리') = 0)
										 THEN USER_ID || '-' || RECORD_DATE
										 END) AS WORKER_COUNT_NOT_MANAGER
					FROM base
					GROUP BY SNO, JNO
				)
				SELECT 
					t1.SNO,
					t1.JNO,
					NVL(wc.WORKER_COUNT_ALL, 0) AS WORKER_COUNT_ALL,
					NVL(wc.WORKER_COUNT_DATE, 0) AS WORKER_COUNT_DATE,
					NVL(wc.WORKER_COUNT_HTENC, 0) AS WORKER_COUNT_HTENC,
					NVL(wc.WORKER_COUNT_MANAGER, 0) AS WORKER_COUNT_MANAGER,
					NVL(wc.WORKER_COUNT_NOT_MANAGER, 0) AS WORKER_COUNT_NOT_MANAGER
				FROM IRIS_SITE_JOB t1
				LEFT JOIN worker_counts wc ON t1.SNO = wc.SNO 
					 AND t1.JNO = wc.JNO
				WHERE t1.SNO > 100`

	if err := db.SelectContext(ctx, &sqlList, query, targetDate, targetDate, targetDate, targetDate); err != nil {
		return nil, fmt.Errorf("GetProjectWorkerCountList err: %v", err)
	}

	return &sqlList, nil
}

// func: 프로젝트별 출근 안전관리자 수
// @param
// - targetDate: 현재시간
func (r *Repository) GetProjectSafeWorkerCountList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.ProjectSafeCountSqls, error) {
	sqlList := entity.ProjectSafeCountSqls{}

	query := `
				WITH htenc_cnt AS (
				  SELECT DISTINCT 
						 SNO,
						 JNO,
						 USER_NM
				  FROM IRIS_RECD_SET
				  WHERE SNO > 100
					AND TO_CHAR(RECOG_TIME, 'YYYY-MM-DD') = TO_CHAR(:1, 'YYYY-MM-DD')
					AND INSTR(NVL(DEPARTMENT, ' '), '하이테크') > 0
				),
				safe_cnt AS (
				  SELECT t1.JNO, USER_NAME
				  FROM JOB_MANAGER t1
				  JOIN S_SYS_USER_SET t2 ON t2.UNO = t1.UNO
				  WHERE t1.AUTH = 'SAFETY_MANAGER'
				  UNION
				  SELECT t1.JNO, USER_NAME
				  FROM TIMESHEET.JOB_MEMBER_LIST t1
				  JOIN S_SYS_USER_SET t2 ON t2.UNO = t1.UNO
				  WHERE t1.COMP_TYPE = 'H'
					AND t1.FUNC_CODE = 510
					AND t1.CHARGE = '21'
					AND t1.IS_USE = 'Y'
				)
				SELECT 
				  ht.SNO,
				  ht.JNO,
				  COUNT(*) AS SAFE_COUNT
				FROM htenc_cnt ht
				WHERE EXISTS (
				  SELECT 1
				  FROM safe_cnt s
				  WHERE s.JNO = ht.JNO
					AND INSTR(ht.USER_NM, s.USER_NAME) > 0
				)
				GROUP BY ht.SNO, ht.JNO`

	if err := db.SelectContext(ctx, &sqlList, query, targetDate); err != nil {
		return nil, fmt.Errorf("GetProjectSafeWorkerCountList err: %v", err)
	}

	return &sqlList, nil
}

// func: 프로젝트 조회(이름)
// @param
// -
func (r *Repository) GetProjectNmList(ctx context.Context, db Queryer) (*entity.ProjectInfoSqls, error) {
	projectInfoSqls := entity.ProjectInfoSqls{}

	sql := `SELECT
    			t1.SNO,
				t1.JNO,
				t2.JOB_NAME as PROJECT_NM
			FROM
				IRIS_SITE_JOB t1
				INNER JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
				INNER JOIN IRIS_SITE_SET t3 ON t1.SNO = t3.SNO
				INNER JOIN TIMESHEET.JOB_KIND_CODE t4 ON t2.JOB_CODE = t4.KIND_CODE
			WHERE t1.sno > 100
			ORDER BY
				t1.IS_DEFAULT DESC`
	if err := db.SelectContext(ctx, &projectInfoSqls, sql); err != nil {
		return &projectInfoSqls, fmt.Errorf("GetProjectNmList fail: %v", err)
	}

	return &projectInfoSqls, nil
}

// func: 공사관리시스템 등록 프로젝트 전체 조회
// @param
// -
func (r *Repository) GetUsedProjectList(ctx context.Context, db Queryer, pageSql entity.PageSql, search entity.JobInfoSql) (*entity.JobInfoSqls, error) {
	sqlData := entity.JobInfoSqls{}

	condition := ""
	condition = utils.StringWhereConvert(condition, search.JobNo, "t2.JOB_NO")
	condition = utils.StringWhereConvert(condition, search.CompName, "t2.COMP_NAME")
	condition = utils.StringWhereConvert(condition, search.OrderCompName, "t2.ORDER_COMP_NAME")
	condition = utils.StringWhereConvert(condition, search.JobName, "t2.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.JobPmName, "t2.JOB_PM_NAME")
	condition = utils.StringWhereConvert(condition, search.JobSd, "t2.JOB_SD")
	condition = utils.StringWhereConvert(condition, search.JobEd, "t2.JOB_ED")
	condition = utils.StringWhereConvert(condition, search.CdNm, "t5.CD_NM")

	var order string
	if pageSql.Order.Valid {
		order = pageSql.Order.String
	} else {
		order = "JOB_NO ASC"
	}

	query := fmt.Sprintf(`
				SELECT *
				FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						SELECT 
							t1.JNO,
							t1.SNO,
							t2.JOB_NAME,
							t2.JOB_NO,
							t2.JOB_SD,
							t2.JOB_ED,
							t2.COMP_NAME,
							t2.ORDER_COMP_NAME,
							t2.JOB_PM_NAME,
							t5.CD_NM
						FROM
							IRIS_SITE_JOB t1
							INNER JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
							INNER JOIN IRIS_SITE_SET t3 ON t1.SNO = t3.SNO
							INNER JOIN TIMESHEET.JOB_KIND_CODE t4 ON t2.JOB_CODE = t4.KIND_CODE
							INNER JOIN TIMESHEET.SYS_CODE_SET t5 ON t5.MINOR_CD = t2.job_state AND t5.major_cd = 'JOB_STATE'
						WHERE t1.SNO > 100
						%s
						ORDER BY %s
					) sorted_data
					WHERE ROWNUM <= :1
				)
				WHERE RNUM > :2`, condition, order)

	if err := db.SelectContext(ctx, &sqlData, query, pageSql.EndNum, pageSql.StartNum); err != nil {
		return nil, fmt.Errorf("GetUsedProjectList err: %w", err)
	}

	return &sqlData, nil
}

// func: 공사관리시스템 등록 프로젝트 전체 조회 개수
// @param
// -
func (r *Repository) GetUsedProjectCount(ctx context.Context, db Queryer, search entity.JobInfoSql) (int, error) {
	var count int

	condition := ""
	condition = utils.StringWhereConvert(condition, search.JobNo, "t2.JOB_NO")
	condition = utils.StringWhereConvert(condition, search.CompName, "t2.COMP_NAME")
	condition = utils.StringWhereConvert(condition, search.OrderCompName, "t2.ORDER_COMP_NAME")
	condition = utils.StringWhereConvert(condition, search.JobName, "t2.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.JobPmName, "t2.JOB_PM_NAME")
	condition = utils.StringWhereConvert(condition, search.JobSd, "t2.JOB_SD")
	condition = utils.StringWhereConvert(condition, search.JobEd, "t2.JOB_ED")
	condition = utils.StringWhereConvert(condition, search.CdNm, "t5.CD_NM")

	query := fmt.Sprintf(`
				SELECT 
					COUNT(*)
				FROM
					IRIS_SITE_JOB t1
					INNER JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
					INNER JOIN IRIS_SITE_SET t3 ON t1.SNO = t3.SNO
					INNER JOIN TIMESHEET.JOB_KIND_CODE t4 ON t2.JOB_CODE = t4.KIND_CODE
					INNER JOIN TIMESHEET.SYS_CODE_SET t5 ON t5.MINOR_CD = t2.job_state AND t5.major_cd = 'JOB_STATE'
				WHERE t1.SNO > 100
				%s`, condition)

	if err := db.GetContext(ctx, &count, query); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("GetUsedProjectCount err: %v", err)
	}

	return count, nil
}

// func: 프로젝트 전체 조회
// @param
// -
func (r *Repository) GetAllProjectList(ctx context.Context, db Queryer, pageSql entity.PageSql, search entity.JobInfoSql) (*entity.JobInfoSqls, error) {
	sqlData := entity.JobInfoSqls{}

	condition := "1 = 1"
	condition = utils.StringWhereConvert(condition, search.JobNo, "J.JOB_NO")
	condition = utils.StringWhereConvert(condition, search.CompName, "J.COMP_NAME")
	condition = utils.StringWhereConvert(condition, search.OrderCompName, "J.ORDER_COMP_NAME")
	condition = utils.StringWhereConvert(condition, search.JobName, "J.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.JobPmName, "J.JOB_PM_NAME")
	condition = utils.StringWhereConvert(condition, search.JobSd, "J.JOB_SD")
	condition = utils.StringWhereConvert(condition, search.JobEd, "J.JOB_ED")
	condition = utils.StringWhereConvert(condition, search.CdNm, "SC.CD_NM")

	var order string
	if pageSql.Order.Valid {
		order = pageSql.Order.String
	} else {
		order = "JNO DESC"
	}

	query := fmt.Sprintf(`
				SELECT *
				FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						SELECT 
							J.JNO,
							S.SNO, 
							J.JOB_NAME, 
							J.JOB_NO, 
							J.JOB_SD, 
							J.JOB_ED, 
							J.COMP_NAME, 
							J.ORDER_COMP_NAME, 
							J.JOB_PM_NAME, 
							SC.CD_NM 
						FROM 
							S_JOB_INFO J 
						INNER JOIN 
							TIMESHEET.job_kind_code JC 
						ON 
							J.job_code = JC.kind_code 
						LEFT JOIN 
							IRIS_SITE_JOB S
						ON
							S.JNO = J.JNO
						INNER JOIN 
							TIMESHEET.SYS_CODE_SET SC 
						ON 
							J.job_state = SC.minor_cd 
							AND SC.MAJOR_CD = 'JOB_STATE' 
						WHERE %s
						ORDER BY %s
					) sorted_data
					WHERE ROWNUM <= :1
				)
				WHERE RNUM > :2`, condition, order)

	if err := db.SelectContext(ctx, &sqlData, query, pageSql.EndNum, pageSql.StartNum); err != nil {
		return nil, fmt.Errorf("GetAllProjectList err: %w", err)
	}

	return &sqlData, nil
}

// func: 프로젝트 개수 조회 개수
// @param
// -
func (r *Repository) GetAllProjectCount(ctx context.Context, db Queryer, search entity.JobInfoSql) (int, error) {
	var count int

	condition := "1 = 1"
	condition = utils.StringWhereConvert(condition, search.JobNo, "J.JOB_NO")
	condition = utils.StringWhereConvert(condition, search.CompName, "J.COMP_NAME")
	condition = utils.StringWhereConvert(condition, search.OrderCompName, "J.ORDER_COMP_NAME")
	condition = utils.StringWhereConvert(condition, search.JobName, "J.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.JobPmName, "J.JOB_PM_NAME")
	condition = utils.StringWhereConvert(condition, search.JobSd, "J.JOB_SD")
	condition = utils.StringWhereConvert(condition, search.JobEd, "J.JOB_ED")
	condition = utils.StringWhereConvert(condition, search.CdNm, "SC.CD_NM")
	query := fmt.Sprintf(`
				SELECT 
					count(*)
				FROM 
					S_JOB_INFO J 
				INNER JOIN 
					TIMESHEET.job_kind_code JC 
				ON 
					J.job_code = JC.kind_code 
				INNER JOIN 
					TIMESHEET.SYS_CODE_SET SC 
				ON 
					J.job_state = SC.minor_cd 
					AND SC.MAJOR_CD = 'JOB_STATE' 
				WHERE %s`, condition)

	if err := db.GetContext(ctx, &count, query); err != nil {
		return 0, fmt.Errorf("GetAllProjectCount err: %w", err)
	}

	return count, nil
}

// func: 본인이 속한 프로젝트 조회
// @param
// - UNO
func (r *Repository) GetStaffProjectList(ctx context.Context, db Queryer, pageSql entity.PageSql, searchSql entity.JobInfoSql, uno sql.NullInt64) (*entity.JobInfoSqls, error) {

	sqlData := entity.JobInfoSqls{}

	condition := "1=1"
	condition = utils.StringWhereConvert(condition, searchSql.JobNo, "J.JOB_NO")
	condition = utils.StringWhereConvert(condition, searchSql.CompName, "J.COMP_NAME")
	condition = utils.StringWhereConvert(condition, searchSql.OrderCompName, "J.ORDER_COMP_NAME")
	condition = utils.StringWhereConvert(condition, searchSql.JobName, "J.JOB_NAME")
	condition = utils.StringWhereConvert(condition, searchSql.JobPmName, "J.JOB_PM_NAME")
	condition = utils.StringWhereConvert(condition, searchSql.JobSd, "J.JOB_SD")
	condition = utils.StringWhereConvert(condition, searchSql.JobEd, "J.JOB_ED")
	condition = utils.StringWhereConvert(condition, searchSql.CdNm, "SC.CD_NM")

	var order string
	if pageSql.Order.Valid {
		order = pageSql.Order.String
	} else {
		order = "JNO DESC"
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM (
			SELECT ROWNUM AS RNUM, sorted_data.*
				FROM (
					SELECT 
						J.JNO,
						S.SNO,
						J.JOB_NAME, 
						J.JOB_NO, 
						J.JOB_SD, 
						J.JOB_ED, 
						J.COMP_NAME, 
						J.ORDER_COMP_NAME, 
						J.JOB_PM_NAME, 
						SC.CD_NM 
					FROM 
						S_JOB_INFO J 
					INNER JOIN 
						(SELECT * FROM TIMESHEET.JOB_MEMBER_LIST WHERE UNO = :1) JM 
					ON 
						J.JNO = JM.JNO
					LEFT JOIN
						IRIS_SITE_JOB S
					ON
						S.JNO = J.JNO
					INNER JOIN 
						TIMESHEET.SYS_CODE_SET SC 
					ON 
						J.job_state = SC.minor_cd 
						AND SC.MAJOR_CD = 'JOB_STATE'
					WHERE %s
					ORDER BY %s
					) sorted_data
				WHERE ROWNUM <= :2
			) 
			WHERE RNUM > :3`, condition, order)

	if err := db.SelectContext(ctx, &sqlData, query, uno, pageSql.EndNum, pageSql.StartNum); err != nil {
		return nil, fmt.Errorf("GetStaffProjectList err: %w", err)
	}

	return &sqlData, nil
}

// func: 본인이 속한 프로젝트 개수
// @param
// - UNO
func (r *Repository) GetStaffProjectCount(ctx context.Context, db Queryer, searchSql entity.JobInfoSql, uno sql.NullInt64) (int, error) {
	var count int

	condition := "1=1"
	condition = utils.StringWhereConvert(condition, searchSql.JobNo, "J.JOB_NO")
	condition = utils.StringWhereConvert(condition, searchSql.CompName, "J.COMP_NAME")
	condition = utils.StringWhereConvert(condition, searchSql.OrderCompName, "J.ORDER_COMP_NAME")
	condition = utils.StringWhereConvert(condition, searchSql.JobName, "J.JOB_NAME")
	condition = utils.StringWhereConvert(condition, searchSql.JobPmName, "J.JOB_PM_NAME")
	condition = utils.StringWhereConvert(condition, searchSql.JobSd, "J.JOB_SD")
	condition = utils.StringWhereConvert(condition, searchSql.JobEd, "J.JOB_ED")
	condition = utils.StringWhereConvert(condition, searchSql.CdNm, "SC.CD_NM")

	query := fmt.Sprintf(`
				SELECT 
					COUNT(*)
				FROM S_JOB_INFO J 
				INNER JOIN 
					(SELECT * FROM TIMESHEET.JOB_MEMBER_LIST WHERE UNO = :1) JM 
				ON 
					J.JNO = JM.JNO 
				INNER JOIN 
					TIMESHEET.SYS_CODE_SET SC 
				ON 
					J.job_state = SC.minor_cd 
					AND SC.MAJOR_CD = 'JOB_STATE'
				WHERE %s`, condition)

	if err := db.GetContext(ctx, &count, query, uno); err != nil {
		return 0, fmt.Errorf("GetStaffProjectCount err: %w", err)
	}

	return count, nil
}

// func: 조직도 조회-고객사
// @param
// - JNO
func (r *Repository) GetClientOrganization(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.OrganizationSqls, error) {
	sqlData := entity.OrganizationSqls{}
	query := fmt.Sprintf(`
				SELECT 
					JM.JNO, 
					JM.FUNC_NAME, 
					JM.CHARGE_DETAIL, 
					JM.MEMBER_NAME AS USER_NAME, 
					JM.GRADE_NAME AS DUTY_NAME, 
					J.ORDER_COMP_NAME AS DEPT_NAME,
					JM.EMAIL, 
					JM.IS_USE, 
					JM.CO_ID, 
					SC.CD_NM, 
					JM.UNO,
					CASE WHEN LENGTH(JM.CELL) > 6 THEN  JM.CELL ELSE '' END CELL, 
					CASE WHEN LENGTH(JM.TEL) > 6 THEN  JM.TEL ELSE '' END TEL	
				FROM 
					JOB_MEMBER_LIST JM 
				INNER JOIN 
					JOB_INFO J 
				ON 
					J.JNO = JM.JNO
				INNER JOIN 
					SYS_CODE_SET SC 
				ON 
					JM.CHARGE = SC.MINOR_CD AND SC.MAJOR_CD = 'MEMBER_CHARGE' 
				WHERE JM.JNO = :1 AND COMP_TYPE = 'O'
				ORDER BY JM.FUNC_NAME ASC, SC.VAL5 ASC, JM.SORT_NO ASC`)

	if err := db.SelectContext(ctx, &sqlData, query, jno); err != nil {
		return nil, fmt.Errorf("GetClientOrganization err: %w", err)
	}

	return &sqlData, nil

}

// func: 조직도 조회-계약자
// @param
// - JNO
func (r *Repository) GetHitechOrganization(ctx context.Context, db Queryer, jno sql.NullInt64, funcNo sql.NullInt64) (*entity.OrganizationSqls, error) {
	sqlData := entity.OrganizationSqls{}

	query := fmt.Sprintf(`
					WITH MEMBER_LIST AS (
						SELECT * FROM JOB_MEMBER_LIST
						WHERE JNO = :1
					)
					,HITECH AS (
							SELECT 
								M.JNO, M.FUNC_CODE, M.CHARGE_DETAIL, U.USER_NAME, U.DUTY_NAME, U.DEPT_NAME, U.CELL, U.TEL, U.EMAIL, U.IS_USE, M.CO_ID, SC.CD_NM, M.UNO, SC.VAL5 AS CHARGE_SORT, U.DUTY_ID
							FROM 
								MEMBER_LIST M 
							INNER JOIN 
								(SELECT 
									UNO, USER_NAME, DUTY_NAME, TEAM_NAME AS DEPT_NAME, CELL, TEL, EMAIL, IS_USE, DUTY_ID 
								FROM 
									S_SYS_USER_SET) U 
							ON 
								M.UNO = U.UNO 
							INNER JOIN 
								SYS_CODE_SET SC 
							ON 
								M.CHARGE = SC.MINOR_CD 
								AND SC.MAJOR_CD = 'MEMBER_CHARGE' 
							WHERE 
								COMP_TYPE = 'H'
						UNION
							SELECT 
								M.JNO, M.FUNC_CODE, M.CHARGE_DETAIL, M.MEMBER_NAME AS USER_NAME, M.GRADE_NAME AS DUTY_NAME, M.DEPT_NAME, M.CELL, M.TEL, M.EMAIL, M.IS_USE, M.CO_ID, SC.CD_NM, M.UNO, SC.VAL5 AS CHARGE_SORT, 99999 AS DUTY_ID 
							FROM 
								MEMBER_LIST M 
							INNER JOIN 
								SYS_CODE_SET SC 
							ON 
								M.CHARGE = SC.MINOR_CD 
								AND SC.MAJOR_CD = 'MEMBER_CHARGE' 
							WHERE 
								COMP_TYPE = 'H' 
								AND UNO IS NULL
					)
					SELECT 
						H.JNO,
						H.CHARGE_DETAIL, 
						H.USER_NAME, 
						H.DUTY_NAME,
						H.DEPT_NAME, 
						H.EMAIL, 
						H.IS_USE,
						H.CO_ID, 
						H.CD_NM, 
						H.UNO,
						CASE WHEN LENGTH(H.CELL) > 6 THEN  H.CELL ELSE '' END CELL, 
						CASE WHEN LENGTH(H.TEL) > 6 THEN  H.TEL ELSE '' END TEL
					FROM 
						HITECH H
					WHERE
						H.FUNC_CODE = :2
					ORDER BY  
						H.CHARGE_SORT ASC,
						H.DUTY_ID ASC,
						H.USER_NAME ASC
					`)

	if err := db.SelectContext(ctx, &sqlData, query, jno, funcNo); err != nil {
		return nil, fmt.Errorf("GetHitechOrganization err: %w", err)
	}
	return &sqlData, nil
}

func (r *Repository) GetFuncNameList(ctx context.Context, db Queryer) (*entity.FuncNameSqls, error) {

	sqlData := entity.FuncNameSqls{}

	query := fmt.Sprintf(`
			SELECT FUNC_NO, FUNC_TITLE
			FROM
				COMMON.V_COMM_FUNC_CODE
			WHERE FUNC_TITLE = 'PM'
		UNION ALL
			SELECT 
				FUNC_NO, FUNC_TITLE
			FROM
				(SELECT * FROM COMMON.V_COMM_FUNC_CODE ORDER BY SORT_NO_PATH)
			WHERE IS_ORG = 'Y'
	`)

	if err := db.SelectContext(ctx, &sqlData, query); err != nil {
		return nil, fmt.Errorf("GetFuncNameList err: %w", err)
	}

	return &sqlData, nil
}

// func: 본인이 속한 프로젝트 이름 조회
// @param
// - uno
func (r *Repository) GetProjectNmUnoList(ctx context.Context, db Queryer, uno sql.NullInt64, role int) (*entity.ProjectInfoSqls, error) {
	projectInfoSqls := entity.ProjectInfoSqls{}

	query := `SELECT 
    			JNO, 
    			JOB_NAME as PROJECT_NM 
			  FROM 
			    S_JOB_INFO 
			  WHERE
			      JOB_STATE = 'Y' AND
			      1=:1 OR 
			      JNO IN (SELECT DISTINCT(JNO) 
						  FROM TIMESHEET.JOB_MEMBER_LIST 
						  WHERE UNO = :2)
			  ORDER BY 
			      JNO DESC`

	if err := db.SelectContext(ctx, &projectInfoSqls, query, role, uno); err != nil {
		return &projectInfoSqls, fmt.Errorf("GetProjectNmList fail: %v", err)
	}

	return &projectInfoSqls, nil
}

// func: 현장근태에 등록되지 않은 프로젝트
// @param
// -
func (r *Repository) GetNonUsedProjectList(ctx context.Context, db Queryer, page entity.PageSql, search entity.NonUsedProjectSql, retry string) (*entity.NonUsedProjectSqls, error) {
	nonProjects := entity.NonUsedProjectSqls{}

	condition := ""

	condition = utils.Int64WhereConvert(condition, search.Jno, "t1.JNO")
	condition = utils.StringWhereConvert(condition, search.JobNo, "t1.JOB_NO")
	condition = utils.StringWhereConvert(condition, search.JobName, "t1.JOB_NAME")
	condition = utils.Int64WhereConvert(condition, search.JobYear, "t1.JOB_YEAR")
	condition = utils.StringWhereConvert(condition, search.JobSd, "t1.JOB_SD")
	condition = utils.StringWhereConvert(condition, search.JobEd, "t1.JOB_ED")
	condition = utils.StringWhereConvert(condition, search.JobPmNm, "t2.USER_NAME")

	var columns []string
	columns = append(columns, "t1.JNO")
	columns = append(columns, "t1.JOB_NO")
	columns = append(columns, "t1.JOB_NAME")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	var order string
	if page.Order.Valid {
		order = page.Order.String
	} else {
		order = `
				CASE 
					WHEN t1.REG_DATE IS NULL THEN t1.MOD_DATE 
					WHEN t1.MOD_DATE IS NULL THEN t1.REG_DATE 
					ELSE GREATEST(t1.REG_DATE, t1.MOD_DATE) 
				END DESC NULLS LAST`
	}

	query := fmt.Sprintf(`
								SELECT *
								FROM (
									SELECT ROWNUM AS RNUM, sorted_data.*
									FROM (
									    SELECT 
											t1.JNO,
											t1.JOB_NO,
											t1.JOB_NAME,
											t1.JOB_YEAR,
											t1.JOB_SD,
											t1.JOB_ED,
											t2.USER_NAME AS JOB_PM_NM,
											t2.DUTY_NAME
										FROM s_job_info t1
										LEFT JOIN COMMON.V_BIZ_USER_INFO t2 ON t1.JOB_PM = t2.UNO
										WHERE t1.JOB_STATE = 'Y'
										AND t1.JNO NOT IN(
											SELECT JNO
											FROM IRIS_SITE_JOB
										)
										%s %s
										ORDER BY %s 
									) sorted_data
									WHERE ROWNUM <= :1
								)
								WHERE RNUM > :2`, condition, retryCondition, order)

	if err := db.SelectContext(ctx, &nonProjects, query, page.EndNum, page.StartNum); err != nil {
		return &nonProjects, fmt.Errorf("GetNonUsedProjectList fail: %v", err)
	}

	return &nonProjects, nil
}

// func: 현장근태에 등록되지 않은 프로젝트 수
// @param
// -
func (r *Repository) GetNonUsedProjectCount(ctx context.Context, db Queryer, search entity.NonUsedProjectSql, retry string) (int, error) {
	var count int

	condition := ""

	condition = utils.Int64WhereConvert(condition, search.Jno, "t1.JNO")
	condition = utils.StringWhereConvert(condition, search.JobNo, "t1.JOB_NO")
	condition = utils.StringWhereConvert(condition, search.JobName, "t1.JOB_NAME")
	condition = utils.Int64WhereConvert(condition, search.JobYear, "t1.JOB_YEAR")
	condition = utils.StringWhereConvert(condition, search.JobSd, "t1.JOB_SD")
	condition = utils.StringWhereConvert(condition, search.JobEd, "t1.JOB_ED")
	condition = utils.StringWhereConvert(condition, search.JobPmNm, "t2.USER_NAME")

	var columns []string
	columns = append(columns, "t1.JNO")
	columns = append(columns, "t1.JOB_NO")
	columns = append(columns, "t1.JOB_NAME")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	query := fmt.Sprintf(`
								SELECT 
									COUNT(*)
								FROM s_job_info t1
								LEFT JOIN COMMON.V_BIZ_USER_INFO t2 ON t1.JOB_PM = t2.UNO
								WHERE t1.JOB_STATE = 'Y'
								AND t1.JNO NOT IN(
									SELECT JNO
									FROM IRIS_SITE_JOB
								)
								%s %s`, condition, retryCondition)

	if err := db.GetContext(ctx, &count, query); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("GetNonUsedProjectCount fail: %v", err)
	}

	return count, nil
}

// func: 현장 프로젝트 추가
// @param
// -
func (r *Repository) AddProject(ctx context.Context, db Beginner, project entity.ReqProject) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddProject. Failed to begin transaction: %w", err)
	}

	agent := utils.GetAgent()

	query := `
			INSERT INTO IRIS_SITE_JOB(
				SNO, JNO, IS_USE, IS_DEFAULT, REG_DATE,
				REG_AGENT, REG_USER, REG_UNO
			) VALUES (
				:1, :2, 'Y', 'N', SYSDATE,
				:3, :4, :5
			)`

	if _, err = tx.ExecContext(ctx, query, project.Sno, project.Jno, agent, project.RegUser, project.RegUno); err != nil {
		origErr := err
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("AddProject. Failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("AddProject. Failed to add project: %w", origErr)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("AddProject. Failed to commit transaction: %w", err)
	}

	return nil
}

// func: 현장 기본 프로젝트 변경
// @param
// -
func (r *Repository) ModifyDefaultProject(ctx context.Context, db Beginner, project entity.ReqProject) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ModifyDefaultProject. Failed to begin transaction: %w", err)
	}
	agent := utils.GetAgent()

	query := `
		UPDATE IRIS_SITE_JOB
		SET IS_DEFAULT = 'N',
		    MOD_AGENT = :1,
		    MOD_USER = :2,
		    MOD_UNO = :3,
		    MOD_DATE = SYSDATE
		WHERE SNO = :4`
	if _, err = tx.ExecContext(ctx, query, agent, project.ModUser, project.ModUno, project.Sno); err != nil {
		origErr := err
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("ModifyDefaultProject; IS_DEFAULT 'N'. Failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("ModifyDefaultProject; IS_DEFAULT 'N'. Failed to modify default project: %w", origErr)
	}

	query = `
		UPDATE IRIS_SITE_JOB
		SET IS_DEFAULT = 'Y',
		    MOD_AGENT = :1,
		    MOD_USER = :2,
		    MOD_UNO = :3,
		    MOD_DATE = SYSDATE
		WHERE SNO = :4
		AND JNO = :5`
	if _, err = tx.ExecContext(ctx, query, agent, project.ModUser, project.ModUno, project.Sno, project.Jno); err != nil {
		origErr := err
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("ModifyDefaultProject; IS_DEFAULT 'Y'. Failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("ModifyDefaultProject IS_DEFAULT 'Y'. Failed to modify default project: %w", origErr)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ModifyDefaultProject: Failed to commit transaction: %w", err)
	}

	return nil
}

// func: 현장 프로젝트 사용여부 변경
// @param
// -
func (r *Repository) ModifyUseProject(ctx context.Context, db Beginner, project entity.ReqProject) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ModifyUseProject. Failed to begin transaction: %w", err)
	}
	agent := utils.GetAgent()

	query := `
		UPDATE IRIS_SITE_JOB
		SET IS_USE = :1,
		    MOD_AGENT = :2,
		    MOD_USER = :3,
		    MOD_UNO = :4,
		    MOD_DATE = SYSDATE
		WHERE SNO = :5
		AND JNO = :6`
	if _, err = tx.ExecContext(ctx, query, project.IsUsed, agent, project.ModUser, project.ModUno, project.Sno, project.Jno); err != nil {
		origErr := err
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("ModifyUseProject. Failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("ModifyUseProject. Failed to modify default project: %w", origErr)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ModifyUseProject: Failed to commit transaction: %w", err)
	}

	return nil
}

// func: 현장 프로젝트 삭제
// @param
// -
func (r *Repository) RemoveProject(ctx context.Context, db Beginner, sno int64, jno int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("RemoveProject. Failed to begin transaction: %w", err)
	}

	query := `
		DELETE FROM IRIS_SITE_JOB		
		WHERE SNO = :1
		AND JNO = :2`
	if _, err = tx.ExecContext(ctx, query, sno, jno); err != nil {
		origErr := err
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("RemoveProject. Failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("RemoveProject. Failed to modify default project: %w", origErr)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("RemoveProject: Failed to commit transaction: %w", err)
	}

	return nil
}
