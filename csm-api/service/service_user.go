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
func (u *ServiceUser) GetUserInfoPeList(ctx context.Context, unoList []int) (*entity.UserPeInfos, error) {
	userPeInfos, err := u.Store.GetUserInfoPeList(ctx, u.SafeDB, unoList)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_user/GetUserInfoPeList err: %w", err)
	}

	return userPeInfos, nil
}
