package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

func (r *Repository) GetUserValid(ctx context.Context, db Queryer, userId string, userPwd string) (entity.User, error) {
	user := entity.User{}

	sql := `SELECT
			t1.UNO,
		    t1.USER_ID,
			t1.USER_NAME
		FROM
			COMMON.V_BIZ_USER_INFO t1
		WHERE
		    t1.USER_ID = :1
			AND t1.USER_PWD = :2`

	if err := db.GetContext(ctx, &user, sql, userId, userPwd); err != nil {
		return user, fmt.Errorf("user.get user fail: %w", err)
	}
	return user, nil
}
