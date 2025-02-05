package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

func (r *Repository) GetUserInfoPmPeList(ctx context.Context, db Queryer, unoList []int) (*entity.UserPmPeInfoSqls, error) {
	userPmPeInfoSqls := entity.UserPmPeInfoSqls{}

	sql := `SELECT
    			t1.UNO,
    			t1.USER_ID,
    			t1.USER_NAME
			FROM COMMON.V_BIZ_USER_INFO t1
			WHERE t1.UNO IN (:1)`

	if err := db.SelectContext(ctx, &userPmPeInfoSqls, sql, unoList); err != nil {
		return nil, fmt.Errorf("GetUserInfoPmPeList fail: %w", err)
	}

	return &userPmPeInfoSqls, nil
}
