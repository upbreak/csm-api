package store

import (
	"context"
	"csm-api/entity"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// 권한별 메뉴 리스트
func (r *Repository) GetParentMenu(ctx context.Context, db Queryer, roles []string) ([]entity.Menu, error) {
	var list []entity.Menu

	query := `
        SELECT 
            T1.MENU_ID,
			MAX(T1.MENU_NM) AS MENU_NM,
			MAX(T1.HAS_CHILD) AS HAS_CHILD,
			MAX(T1.SVG_NAME) AS SVG_NAME,
			MAX(T1.IS_TEMP) AS IS_TEMP
        FROM IRIS_MENU_SET T1
        LEFT JOIN IRIS_USER_MENU T2 ON T1.MENU_ID = T2.MENU_ID 
        WHERE T1.IS_USE = 'Y'
        AND T2.IS_USE = 'Y'
        AND T1.PARENT_ID IS NULL
        AND T2.ROLE_CODE IN (?)
        GROUP BY T1.MENU_ID
		ORDER BY MAX(MENU_ORDER)
    `

	query, args, err := sqlx.In(query, roles)
	if err != nil {
		return nil, fmt.Errorf("GetParentMenu sqlx.In error: %v", err)
	}
	query = db.Rebind(query)

	if err = db.SelectContext(ctx, &list, query, args...); err != nil {
		return nil, fmt.Errorf("GetParentMenu err: %v", err)
	}

	return list, nil
}

// 권한별 서브 메뉴 리스트
func (r *Repository) GetChildMenu(ctx context.Context, db Queryer, roles []string) ([]entity.Menu, error) {
	var list []entity.Menu

	query := `
			SELECT
				T1.MENU_ID,
				MAX(T1.PARENT_ID) AS PARENT_ID,
				MAX(T1.MENU_NM) AS MENU_NM,
				MAX(T1.SVG_NAME) AS SVG_NAME,
				MAX(T1.IS_TEMP) AS IS_TEMP
			FROM IRIS_MENU_SET T1
			LEFT JOIN IRIS_USER_MENU T2 ON T1.MENU_ID = T2.MENU_ID 
			WHERE T1.IS_USE = 'Y'
			AND T2.IS_USE = 'Y'
			AND T1.PARENT_ID IS NOT NULL
			AND T2.ROLE_CODE IN (?)
			GROUP BY T1.MENU_ID
			ORDER BY MAX(T1.SUB_ORDER)`

	query, args, err := sqlx.In(query, roles)
	if err != nil {
		return nil, fmt.Errorf("GetChildMenu sqlx.In error: %v", err)
	}
	query = db.Rebind(query)

	if err = db.SelectContext(ctx, &list, query, args...); err != nil {
		return nil, fmt.Errorf("GetChildMenu err: %v", err)
	}
	return list, nil
}
