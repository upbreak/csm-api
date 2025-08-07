package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
)

type ServiceUser struct {
	SafeDB      store.Queryer
	TimeSheetDB store.Queryer
	Store       store.UserStore
}

// 프로젝트 pm, pe 정보 조회
//
// @param unoList: ?? 고유번호 리스트
func (u *ServiceUser) GetUserInfoPeList(ctx context.Context, unoList []int) (*entity.UserPeInfos, error) {
	userPeInfos, err := u.Store.GetUserInfoPeList(ctx, u.SafeDB, unoList)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return userPeInfos, nil
}

// 사용자 권한 조회 {현장소장 | 현장관리자 | 안전관리자 | 관리감독자 | 협력업체관리자}
func (u *ServiceUser) GetUserRole(ctx context.Context, jno int64, uno int64) (string, error) {
	role1, err := u.Store.GetSiteRole(ctx, u.SafeDB, jno, uno)
	if err != nil {
		return "", utils.CustomErrorf(err)
	}
	if role1 != "" {
		return role1, nil
	}

	role2, err := u.Store.GetOperationalRole(ctx, u.SafeDB, jno, uno)
	if err != nil {
		return "", utils.CustomErrorf(err)
	}

	if role2 != "" {
		return role2, nil
	}

	return "", nil
}

// 각 기능별 권한 리스트 불러오기
func (u *ServiceUser) GetAuthorizationList(ctx context.Context, api string) (*entity.RoleList, error) {

	list, err := u.Store.GetAuthorizationList(ctx, u.SafeDB, api)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return list, nil

}
