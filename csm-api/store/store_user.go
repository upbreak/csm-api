package store

import (
	"context"
	"csm-api/entity"
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
