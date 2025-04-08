package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"errors"
	"fmt"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-18
 * @modified 최종 수정일: 2025-02-26
 * @modifiedBy 최종 수정자: 정지영
 * @modified description
 * - 현장소장 및 안전관리자 UserId, UserInfo 추가
 */

// func: job 정보 조회
// @param
// - jno sql.NullInt64: 프로젝트 고유번호
func (r *Repository) GetJobInfo(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.JobInfo, error) {
	data := entity.JobInfo{}

	query := `
				SELECT 
					t1.JNO,
					t2.job_name,
					t2.job_no,
					t2.job_sd,
					t2.job_ed,
					t2.comp_name,
					t2.order_comp_name,
					t2.job_pm_name,
					t6.duty_name as JOB_PM_DUTY_NAME,
					t5.cd_nm
				FROM
					IRIS_SITE_JOB t1
					INNER JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
					INNER JOIN IRIS_SITE_SET t3 ON t1.SNO = t3.SNO
					INNER JOIN TIMESHEET.JOB_KIND_CODE t4 ON t2.JOB_CODE = t4.KIND_CODE
					INNER JOIN TIMESHEET.SYS_CODE_SET t5 ON t5.MINOR_CD = t2.job_state AND t5.major_cd = 'JOB_STATE'
					INNER JOIN S_SYS_USER_SET t6 ON t2.JOB_PM = t6.UNO
				WHERE 
					t1.JNO = :1`
	if err := db.GetContext(ctx, &data, query, jno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &data, nil
		}
		return nil, fmt.Errorf("GetJobInfo fail: %w", err)
	}

	return &data, nil
}

// func: 현장소장 조회
// @param
// - jno sql.NullInt64: 프로젝트 고유번호
func (r *Repository) GetSiteManagerList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.Managers, error) {
	list := entity.Managers{}

	query := `
				SELECT 
					 m.JNO,
					 U.UNO, 
					 U.USER_NAME, 
					 U.DUTY_NAME,
					 U.USER_ID,
					 U.USER_NAME || ' ' || U.DUTY_NAME || ' (' || U.USER_ID || ')' AS USER_INFO
				FROM JOB_MEMBER_LIST M, 
					 S_SYS_USER_SET U
				WHERE M.COMP_TYPE = 'H'
				AND U.UNO = M.UNO
				AND M.FUNC_CODE = 510
				AND M.CHARGE = '21'
				AND M.JNO = :1
				AND M.IS_USE = 'Y'`
	if err := db.SelectContext(ctx, &list, query, jno); err != nil {
		return nil, fmt.Errorf("GetSiteManagerList err: %v", err)
	}
	return &list, nil
}

// func: 안전관리자 조회
// @param
// - jno sql.NullInt64: 프로젝트 고유번호
func (r *Repository) GetSafeManagerList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.Managers, error) {
	list := entity.Managers{}

	query := `
				SELECT 
						 U.UNO, 
						 U.USER_NAME, 
						 U.DUTY_NAME, 
 						 U.USER_ID,
						 J.TEAM_LEADER
					FROM JOB_MANAGER J,
						 S_SYS_USER_SET U
				   WHERE U.UNO = J.UNO
					 AND J.AUTH = 'SAFETY_MANAGER'
					 AND J.JNO = :1
				ORDER BY TEAM_LEADER DESC, 
						 U.DUTY_CD, 
						 U.JOBDUTY_ID, 
						 U.JOIN_DATE, 
						 U.USER_NAME`
	if err := db.SelectContext(ctx, &list, query, jno); err != nil {
		return nil, fmt.Errorf("GetSafeManagerList err: %v", err)
	}
	return &list, nil
}

// func: 관리감독자 조회
// @param
// - jno sql.NullInt64: 프로젝트 고유번호
func (r *Repository) GetSupervisorList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.Supervisors, error) {
	list := entity.Supervisors{}

	query := `
				SELECT 
					 M.UNO, 
					 M.JNO, 
					 U.USER_NAME, 
					 U.USER_ID, 
					 U.DUTY_NAME, 
					 U.DUTY_CD, 
					 U.JOBDUTY_ID, 
					 U.JOIN_DATE,
					 LISTAGG(F.FUNC_NO, '|') WITHIN GROUP(ORDER BY M.UNO) AS FUNC_NO
				FROM JOB_MANAGER M
				RIGHT OUTER JOIN JOB_MANAGER_FUNC F 
				  ON F.JNO = M.JNO AND F.UNO = M.UNO
				JOIN S_SYS_USER_SET U ON M.UNO = U.UNO
				WHERE M.JNO = :1
				 AND M.AUTH = 'SUPERVISOR' 
				GROUP BY M.UNO, 
					 M.JNO, 
					 U.USER_NAME, 
					 U.USER_ID, 
					 U.DUTY_NAME, 
					 U.DUTY_CD, 
					 U.JOBDUTY_ID, 
					 U.JOIN_DATE, 
					 U.USER_NAME
				ORDER BY U.DUTY_CD, 
					 U.JOBDUTY_ID, 
					 U.JOIN_DATE, 
					 U.USER_NAME`
	if err := db.SelectContext(ctx, &list, query, jno); err != nil {
		return nil, fmt.Errorf("GetSupervisorList err: %v", err)
	}
	return &list, nil
}

// func: 공종 정보 조회
// @param
func (r *Repository) GetWorkInfoList(ctx context.Context, db Queryer) (*entity.WorkInfos, error) {
	list := entity.WorkInfos{}

	query := `
				SELECT 
						 FUNC_NO, 
						 FUNC_NAME
					FROM COMMON.COMM_FUNC_QHSE
				   WHERE IS_USE = 'Y'
				ORDER BY SORT_NO`
	if err := db.SelectContext(ctx, &list, query); err != nil {
		return nil, fmt.Errorf("GetWorkInfoList err: %v", err)
	}
	return &list, nil
}

// func: 협력업체 정보 조회
// @param
// - jno sql.NullInt64: 프로젝트 고유번호
func (r *Repository) GetCompanyInfoList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.CompanyInfos, error) {
	list := entity.CompanyInfos{}

	query := `
				SELECT S.JNO, 
					 S.CNO, 
					 S.ID, 
					 S.PW, 
					 S.CELLPHONE, 
					 S.EMAIL, 
					 U.USER_NAME, 
					 U.DUTY_NAME
				FROM JOB_SUBCON_INFO S, 
					 S_SYS_USER_SET U
				WHERE S.UNO = U.UNO(+)
				 AND S.IS_USE = 'Y'
				 AND S.JNO = :1
				ORDER BY S.CNO`
	if err := db.SelectContext(ctx, &list, query, jno); err != nil {
		return nil, fmt.Errorf("GetCompanyInfoList err: %v", err)
	}
	return &list, nil
}

// func: 협력업체별 공종 조회
// @param
// - jno sql.NullInt64: 프로젝트 고유번호
func (r *Repository) GetCompanyWorkInfoList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.WorkInfos, error) {
	list := entity.WorkInfos{}

	query := `
				SELECT F.JNO, 
					 F.CNO, 
					 C.FUNC_NO, 
					 C.FUNC_NAME
					 
				FROM JOB_SUBCON_FUNC F, 
					 COMMON.COMM_FUNC_QHSE C
				WHERE F.FUNC_NO = C.FUNC_NO
				 AND C.IS_USE = 'Y'
				 AND F.JNO = :1`
	if err := db.SelectContext(ctx, &list, query, jno); err != nil {
		return nil, fmt.Errorf("GetCompanyWorkInfoList err: %v", err)
	}
	return &list, nil
}
