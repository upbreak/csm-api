package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceUser struct {
	SafeDB store.Queryer
	Store  store.UserStore
}

// 프로젝트 pm, pe 정보 조회
//
// @param unoList: ?? 고유번호 리스트
func (u *ServiceUser) GetUserInfoPmPeList(ctx context.Context, unoList []int) (*entity.UserPmPeInfos, error) {
	userPmPeInfos, err := u.Store.GetUserInfoPmPeList(ctx, u.SafeDB, unoList)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_user/GetUserInfoPmPeList err: %w", err)
	}

	return userPmPeInfos, nil
}
