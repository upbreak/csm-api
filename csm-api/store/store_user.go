package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (r *Repository) GetUserInfoPeList(ctx context.Context, db Queryer, unoList []int) (*entity.UserPeInfos, error) {
	userPeInfos := entity.UserPeInfos{}

	if len(unoList) == 0 {
		return &entity.UserPeInfos{}, nil
	}

	placeholders := make([]string, len(unoList))
	args := make([]interface{}, len(unoList))

	for i, uno := range unoList {
		placeholder := fmt.Sprintf(":p%d", i+1)
		placeholders[i] = placeholder
		args[i] = uno
	}

	sql := fmt.Sprintf(`SELECT
    			t1.UNO,
    			t1.USER_ID,
    			t1.USER_NAME
			FROM COMMON.V_BIZ_USER_INFO t1
			WHERE t1.UNO IN (%s)`, strings.Join(placeholders, ","))

	if err := db.SelectContext(ctx, &userPeInfos, sql, args...); err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return &userPeInfos, nil
}

// 현장소장, 현장 관리자 권한 조회
// @param
// - uno: 유저PK
// - jno: 프로젝트PK (최초 로그인 시 jno를 0으로 부여하여 전체 프로젝트에서 권한 확인)
func (r *Repository) GetSiteRole(ctx context.Context, db Queryer, jno int64, uno int64) (string, error) {
	var role string

	query := `
		WITH MEMBER_LIST AS (
			SELECT * FROM TIMESHEET.JOB_MEMBER_LIST
			WHERE UNO = :1
			AND (0 = :2 OR JNO = :3)
		),
		TRIMMED_CODE AS (
			SELECT 
				M.*,
				REPLACE(TRIM(SC.CD_NM), ' ', '') AS CLEAN_CD_NM
			FROM MEMBER_LIST M
			INNER JOIN TIMESHEET.SYS_CODE_SET SC ON M.CHARGE = SC.MINOR_CD 
			WHERE SC.MAJOR_CD = 'MEMBER_CHARGE'
			AND M.COMP_TYPE = 'H'
		)
		SELECT *
		FROM (
		  SELECT 
		    CASE
		      WHEN CLEAN_CD_NM LIKE '%현장소장%' THEN 'SITE_DIRECTOR'
		      WHEN CLEAN_CD_NM LIKE '%공무%' OR CLEAN_CD_NM LIKE '%사무보조%' THEN 'SITE_MANAGER'
		      ELSE NULL
		    END AS ROLE
		  FROM TRIMMED_CODE
		  ORDER BY ROLE ASC
		)
		WHERE ROWNUM = 1`

	if err := db.GetContext(ctx, &role, query, uno, jno, jno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", utils.CustomErrorf(err)
	}
	return role, nil
}

// 안전관리자, 관리감독자 조회
// @param
// - uno: 유저PK
// - jno: 프로젝트PK (최초 로그인 시 jno를 0으로 부여하여 전체 프로젝트에서 권한 확인)
func (r *Repository) GetOperationalRole(ctx context.Context, db Queryer, jno int64, uno int64) (string, error) {
	var role string

	query := `
		SELECT * FROM (
			SELECT ROLE FROM (
					SELECT 'SAFETY_MANAGER' AS ROLE, COUNT(*) AS CNT
					FROM JOB_MANAGER J
					JOIN S_SYS_USER_SET U ON U.UNO = J.UNO
					WHERE J.AUTH = 'SAFETY_MANAGER'
					 AND J.UNO = :1
					 AND (0 = :2 OR J.JNO = :3)
				UNION ALL
					SELECT 'SUPERVISOR' AS ROLE, COUNT(*) AS CNT
					FROM JOB_MANAGER M
					JOIN S_SYS_USER_SET U ON M.UNO = U.UNO
					JOIN JOB_MANAGER_FUNC F ON F.JNO = M.JNO AND F.UNO = M.UNO
					WHERE M.AUTH = 'SUPERVISOR'
					 AND M.UNO = :4
					 AND (0 = :5 OR M.JNO = :6)
				)
			WHERE CNT > 0
		) WHERE ROWNUM = 1
				`

	if err := db.GetContext(ctx, &role, query, uno, jno, jno, uno, jno, jno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", utils.CustomErrorf(err)
	}
	return role, nil
}

// func: 기능 별로 권한 조회
// @parms
// - api : 조회할 기능 문자열
func (r *Repository) GetAuthorizationList(ctx context.Context, db Queryer, api string) (*entity.RoleList, error) {
	list := entity.RoleList{}

	query := `
		SELECT * 
		FROM 
		    IRIS_LIST_PERMIT_ROLE
		WHERE
		    API = :1 
		`

	if err := db.SelectContext(ctx, &list, query, api); err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return &list, nil
}

// func: 안전보건 시스템에 등록되지 않은 관리감독자 조회
// @params
// - uno: 유저PK
func (r *Repository) GetSupervisorRole(ctx context.Context, db Queryer, uno int64) (string, error) {
	var role string

	query := `
		WITH MEMBER_LIST AS (
			SELECT * FROM TIMESHEET.JOB_MEMBER_LIST
			WHERE UNO = :1
		),
		TRIMMED_CODE AS (
			SELECT 
				COUNT(*) AS CNT
			FROM MEMBER_LIST M
			INNER JOIN TIMESHEET.SYS_CODE_SET SC ON M.CHARGE = SC.MINOR_CD 
			WHERE SC.MAJOR_CD = 'MEMBER_CHARGE'
			AND M.COMP_TYPE = 'H'
			AND M.FUNC_CODE = 510
		)
		SELECT 
			'SUPERVISOR' AS ROLE 
		FROM 
			TRIMMED_CODE 
		WHERE 
			CNT > 0
	`

	if err := db.GetContext(ctx, &role, query, uno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", utils.CustomErrorf(err)
	}
	return role, nil
}
