package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"fmt"
)

// func: 조직도 조회-고객사
// @param
// - JNO
func (r *Repository) GetOrganizationClientList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.OrganizationSqls, error) {
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
					S_JOB_MEMBER_LIST JM 
				INNER JOIN 
					S_JOB_INFO J 
				ON 
					J.JNO = JM.JNO
				INNER JOIN 
					SYS_CODE_SET SC 
				ON 
					JM.CHARGE = SC.MINOR_CD AND SC.MAJOR_CD = 'MEMBER_CHARGE' 
				WHERE JM.JNO = :1 AND COMP_TYPE = 'O'
				ORDER BY JM.FUNC_NAME ASC, SC.VAL5 ASC, JM.SORT_NO ASC`)

	if err := db.SelectContext(ctx, &sqlData, query, jno); err != nil {
		return nil, fmt.Errorf("GetOrganizationClientList err: %w", err)
	}

	return &sqlData, nil
}

// func: 조직도 조회-계약자
// @param
// - JNO
func (r *Repository) GetOrganizationHtencList(ctx context.Context, db Queryer, jno sql.NullInt64, funcNo sql.NullInt64) (*entity.OrganizationSqls, error) {
	sqlData := entity.OrganizationSqls{}

	query := fmt.Sprintf(`
					WITH MEMBER_LIST AS (
						SELECT * FROM S_JOB_MEMBER_LIST
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
		return nil, fmt.Errorf("GetOrganizationHtencList err: %w", err)
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
