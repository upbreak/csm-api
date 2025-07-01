package store

import (
	"context"
	"csm-api/entity"
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
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("GetUserInfoPeList fail: %w", err)
	}

	return &userPeInfos, nil
}

// 현장소장, 현장 관리자 권한 조회
func (r *Repository) GetSiteRole(ctx context.Context, db Queryer, jno int64, uno int64) (string, error) {
	var role string

	query := `
		WITH MEMBER_LIST AS (
			SELECT * FROM JOB_MEMBER_LIST
			WHERE JNO = :1
			AND UNO = :2
		),
		TRIMMED_CODE AS (
			SELECT 
				M.*,
				REPLACE(TRIM(SC.CD_NM), ' ', '') AS CLEAN_CD_NM
			FROM MEMBER_LIST M
			INNER JOIN SYS_CODE_SET SC ON M.CHARGE = SC.MINOR_CD 
			WHERE SC.MAJOR_CD = 'MEMBER_CHARGE'
			AND M.COMP_TYPE = 'H'
		)
		SELECT 
			CASE
				WHEN CLEAN_CD_NM LIKE '%현장소장%' THEN 'SITE_DIRECTOR'
				WHEN CLEAN_CD_NM LIKE '%공무%' OR CLEAN_CD_NM LIKE '%사무보조%' THEN 'SITE_MANAGER'
				ELSE NULL
			END AS ROLE
		FROM TRIMMED_CODE`

	if err := db.GetContext(ctx, &role, query, jno, uno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("GetSiteRole fail: %w", err)
	}
	return role, nil
}

// 안전관리자, 관리감독자, 협렵업체관리자 권한 조회
func (r *Repository) GetOperationalRole(ctx context.Context, db Queryer, jno int64, uno int64) (string, error) {
	var role string

	query := `
		SELECT ROLE FROM (
			SELECT 'SAFETY_MANAGER' AS ROLE, COUNT(*) AS CNT
			FROM JOB_MANAGER J
			JOIN S_SYS_USER_SET U ON U.UNO = J.UNO
			WHERE J.AUTH = 'SAFETY_MANAGER'
			 AND J.JNO = :1
			 AND J.UNO = :2
			UNION ALL
			SELECT 'SUPERVISOR' AS ROLE, COUNT(*) AS CNT
			FROM JOB_MANAGER M
			JOIN S_SYS_USER_SET U ON M.UNO = U.UNO
			JOIN JOB_MANAGER_FUNC F ON F.JNO = M.JNO AND F.UNO = M.UNO
			WHERE M.AUTH = 'SUPERVISOR'
			 AND M.JNO = :3
			 AND M.UNO = :4
			UNION ALL
			SELECT 'CO_MANAGER' AS ROLE, COUNT(*) AS CNT
			FROM JOB_SUBCON_INFO S
			LEFT JOIN S_SYS_USER_SET U ON S.UNO = U.UNO
			WHERE S.IS_USE = 'Y'
			 AND S.JNO = :5
			 AND S.ID = :6
		)
		WHERE CNT > 0`

	if err := db.GetContext(ctx, &role, query, jno, uno, jno, uno, jno, uno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("GetOperationalRole fail: %w", err)
	}
	return role, nil
}
