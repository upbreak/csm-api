package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"database/sql"
	"fmt"
)

// func: 현장 프로젝트 조회
// @param
// - sno int64 현장 번호
func (r *Repository) GetProjectList(ctx context.Context, db Queryer, sno int64) (*entity.ProjectInfoSqls, error) {
	projectInfoSqls := entity.ProjectInfoSqls{}

	// sno 변환: 0이면 NULL 처리, 아니면 Valid 값으로 설정
	var snoParam sql.NullInt64
	if sno != 0 {
		snoParam = sql.NullInt64{Valid: true, Int64: sno}
	} else {
		snoParam = sql.NullInt64{Valid: false}
	}

	sql := `SELECT
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
				t2.JOB_TYPE as PROJECT_TYPE,
				t6.CD_NM AS PROJECT_TYPE_NM,
				t2.JOB_NO as PROJECT_NO,
				t2.JOB_NAME as PROJECT_NM,
				t2.JOB_YEAR as PROJECT_YEAR,
				t2.JOB_LOC as PROJECT_LOC,
				t2.JOB_CODE as PROJECT_CODE,
				t4.KIND_NAME as PROJECT_CODE_NAME,
				t3.SITE_NM,
				t2.COMP_CODE,
				t2.COMP_NICK,
				t2.COMP_NAME,
				t2.COMP_ETC,
				t2.ORDER_COMP_CODE,
				t2.ORDER_COMP_NICK,
				t2.ORDER_COMP_NAME,
				t2.ORDER_COMP_JOB_NAME,
				t2.JOB_LOC_NAME as PROJECT_LOC_NAME,
				t2.JOB_PM,
				t2.JOB_PM_NAME,
				t2.JOB_PE,
				TO_DATE(t2.JOB_SD, 'YYYY-MM-DD') as PROJECT_STDT,
				TO_DATE(t2.JOB_ED, 'YYYY-MM-DD') as PROJECT_EDDT,
				TO_DATE(t2.REG_DATE, 'YYYY-MM-DD HH24:MI:SS') as PROJECT_REG_DATE,
				TO_DATE(t2.MOD_DATE, 'YYYY-MM-DD HH24:MI:SS') as PROJECT_MOD_DATE,
				t2.JOB_STATE as PROJECT_STATE,
				t5.CD_NM as PROJECT_STATE_NM,
				t2.MOC_NO,
				t2.WO_NO
			FROM
				IRIS_SITE_JOB t1
				INNER JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
				INNER JOIN IRIS_SITE_SET t3 ON t1.SNO = t3.SNO
				INNER JOIN TIMESHEET.JOB_KIND_CODE t4 ON t2.JOB_CODE = t4.KIND_CODE
				INNER JOIN TIMESHEET.SYS_CODE_SET t5 ON t5.MINOR_CD = t2.job_state AND t5.major_cd = 'JOB_STATE'
				INNER JOIN TIMESHEET.SYS_CODE_SET t6 ON t6.MINOR_CD = t2.JOB_TYPE AND t6.major_cd = 'JOB_TYPE'
			WHERE
				t1.SNO > 100
				AND (:1 IS NULL OR t1.SNO = :2)
			ORDER BY
				t1.IS_DEFAULT DESC`
	if err := db.SelectContext(ctx, &projectInfoSqls, sql, snoParam, snoParam); err != nil {
		return &projectInfoSqls, fmt.Errorf("getProjectList fail: %v", err)
	}

	return &projectInfoSqls, nil
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

// func: 프로젝트 전체 조회
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

// func: 프로젝트 전체 조회 개수
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

// func: 프로젝트 개수 조회
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

// func: 조직도 확인
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

// func: 조직도 확인 개수
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

// func: 조직도 조회
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
					JM.DEPT_NAME, 
					JM.CELL, 
					JM.TEL, 
					JM.EMAIL, 
					JM.IS_USE, 
					JM.CO_ID, 
					SC.CD_NM, 
					JM.UNO 
				FROM 
					JOB_MEMBER_LIST JM 
				INNER JOIN 
					SYS_CODE_SET SC 
				ON 
					JM.CHARGE = SC.MINOR_CD AND SC.MAJOR_CD = 'MEMBER_CHARGE' 
				WHERE JNO = :1 AND COMP_TYPE = 'O'`)

	if err := db.SelectContext(ctx, &sqlData, query, jno); err != nil {
		return nil, fmt.Errorf("GetClientOrganization err: %w", err)
	}

	return &sqlData, nil

}

func (r *Repository) GetHitechOrganization(ctx context.Context, db Queryer, jno sql.NullInt64, funcNo sql.NullInt64) (*entity.OrganizationSqls, error) {
	sqlData := entity.OrganizationSqls{}

	query := fmt.Sprintf(`
					WITH MEMBER_LIST AS (
						SELECT * FROM JOB_MEMBER_LIST
						WHERE JNO = :1
					)
					,HITECH AS (
							SELECT 
								M.JNO, M.FUNC_CODE, M.CHARGE_DETAIL, U.USER_NAME, U.DUTY_NAME, U.DEPT_NAME, U.CELL, U.TEL, U.EMAIL, M.IS_USE, M.CO_ID, SC.CD_NM, M.UNO 
							FROM 
								MEMBER_LIST M 
							INNER JOIN 
								(SELECT 
									UNO, USER_NAME, DUTY_NAME, DEPT_NAME, CELL, TEL, EMAIL 
								FROM 
									S_SYS_USER_SET 
								ORDER BY 
									DUTY_ID) U 
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
								M.JNO, M.FUNC_CODE, M.CHARGE_DETAIL, M.MEMBER_NAME AS USER_NAME, M.GRADE_NAME AS DUTY_NAME, M.DEPT_NAME, M.CELL, M.TEL, M.EMAIL, M.IS_USE, M.CO_ID, SC.CD_NM, M.UNO 
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
						H.CELL, 
						H.TEL, 
						H.EMAIL, 
						H.IS_USE,
						H.CO_ID, 
						H.CD_NM, 
						H.UNO
					FROM 
						HITECH H
					WHERE
						H.FUNC_CODE = :2
					ORDER BY  
						CD_NM DESC
					`)

	if err := db.SelectContext(ctx, &sqlData, query, jno, funcNo); err != nil {
		return nil, fmt.Errorf("GetHitechOrganization err: %w", err)
	}
	return &sqlData, nil
}

func (r *Repository) GetFuncNameList(ctx context.Context, db Queryer) (*entity.FuncNameSqls, error) {

	sqlData := entity.FuncNameSqls{}

	query := fmt.Sprintf(`
		SELECT 
			FUNC_NO, FUNC_TITLE
		FROM
			COMMON.V_COMM_FUNC_CODE
		ORDER BY 
			CASE WHEN 
				FUNC_TITLE = 'PM' 
				THEN 0 
				ELSE 1 
				END, 
			FUNC_NO
	`)

	if err := db.SelectContext(ctx, &sqlData, query); err != nil {
		return nil, fmt.Errorf("GetFuncNameList err: %w", err)
	}

	return &sqlData, nil
}
