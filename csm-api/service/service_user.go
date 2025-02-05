package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceUser struct {
	DB    store.Queryer
	Store store.UserStore
}

// 프로젝트 pm, pe 정보 조회
//
// @param unoList: ?? 고유번호 리스트
func (u *ServiceUser) GetUserInfoPmPeList(ctx context.Context, unoList []int) (*entity.UserPmPeInfos, error) {
	userPmPeInfoSqls, err := u.Store.GetUserInfoPmPeList(ctx, u.DB, unoList)
	if err != nil {
		return nil, fmt.Errorf("service_user/GetUserInfoPmPeList err: %w", err)
	}
	userPmPeInfos := &entity.UserPmPeInfos{}
	userPmPeInfos.ToUserPmPeInfos(userPmPeInfoSqls)

	return userPmPeInfos, nil
}
