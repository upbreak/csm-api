package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/txutil"
	"csm-api/utils"
)

type ServiceUserRole struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.UserRoleStore
}

// 사용자 권한 조회
// param: 사용자번호
func (s *ServiceUserRole) GetUserRoleListByUno(ctx context.Context, uno int64) ([]entity.UserRoleMap, error) {
	list, err := s.Store.GetUserRoleListByUno(ctx, s.SafeDB, uno)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}
	return list, nil
}

// 사용자 권한 조회
// param: 권한코드, 프로젝트 번호
func (s *ServiceUserRole) GetUserRoleListByCodeAndJno(ctx context.Context, code string, jno int64) ([]entity.UserRoleMap, error) {
	list, err := s.Store.GetUserRoleListByCodeAndJno(ctx, s.SafeDB, code, jno)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}
	return list, nil
}

// 사용자 권한 추가
func (s *ServiceUserRole) AddUserRole(ctx context.Context, userRoles []entity.UserRoleMap) (err error) {
	tx, err := txutil.BeginTxWithMode(ctx, s.SafeTDB, false)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer txutil.DeferTx(tx, &err)

	if err = s.Store.AddUserRole(ctx, tx, userRoles); err != nil {
		err = utils.CustomErrorf(err)
	}
	return
}

// 사용자 권한 삭제
func (s *ServiceUserRole) RemoveUserRole(ctx context.Context, userRoles []entity.UserRoleMap) (err error) {
	tx, err := txutil.BeginTxWithMode(ctx, s.SafeTDB, false)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer txutil.DeferTx(tx, &err)

	if err = s.Store.RemoveUserRole(ctx, tx, userRoles); err != nil {
		err = utils.CustomErrorf(err)
	}
	return
}
