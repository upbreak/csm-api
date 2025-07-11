package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
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
		return nil, fmt.Errorf("service_user_role/GetUserRoleListByUno err: %w", err)
	}
	return list, nil
}

// 사용자 권한 조회
// param: 권한코드, 프로젝트 번호
func (s *ServiceUserRole) GetUserRoleListByCodeAndJno(ctx context.Context, code string, jno int64) ([]entity.UserRoleMap, error) {
	list, err := s.Store.GetUserRoleListByCodeAndJno(ctx, s.SafeDB, code, jno)
	if err != nil {
		return nil, fmt.Errorf("service_user_role/GetUserRoleListByCodeAndJno err: %w", err)
	}
	return list, nil
}

// 사용자 권한 추가
func (s *ServiceUserRole) AddUserRole(ctx context.Context, userRoles []entity.UserRoleMap) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("service_user_role/AddUserRole BeginTx error: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_user_role/AddUserRole panic: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_user_role/AddUserRole Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_user_role/AddUserRole Commit error: %w", commitErr)
			}
		}
	}()

	if err = s.Store.AddUserRole(ctx, tx, userRoles); err != nil {
		err = fmt.Errorf("service_user_role/AddUserRole error: %w", err)
	}
	return
}

// 사용자 권한 삭제
func (s *ServiceUserRole) RemoveUserRole(ctx context.Context, userRoles []entity.UserRoleMap) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("service_user_role/RemoveUserRole BeginTx error: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_user_role/RemoveUserRole panic error: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_user_role/RemoveUserRole Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_user_role/RemoveUserRole Commit error: %w", commitErr)
			}
		}
	}()
	if err = s.Store.RemoveUserRole(ctx, tx, userRoles); err != nil {
		err = fmt.Errorf("service_user_role/RemoveUserRole error: %w", err)
	}
	return
}
