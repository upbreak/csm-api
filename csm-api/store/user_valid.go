package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

func (r *Repository) GetUserValid(ctx context.Context, db Queryer, userId string, userPwd string) (entity.User, error) {
	user := entity.User{}

	sql := `
		SELECT
			T1.UNO,
			T1.USER_ID,
			T1.USER_NAME,
			T1.DEPT_NAME,
			T1.TEAM_NAME,
			T2.ROLE_CODE
		FROM
			COMMON.V_BIZ_USER_INFO T1
			LEFT JOIN IRIS_USER_ROLE_MAP T2 ON T1.UNO = T2.USER_UNO
		WHERE T1.IS_USE = 'Y'
		AND T1.USER_ID = :1
		AND t1.USER_PWD = :2`

	if err := db.GetContext(ctx, &user, sql, userId, userPwd); err != nil {
		return user, fmt.Errorf("user.get user fail: %w", err)
	}
	return user, nil
}
