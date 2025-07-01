package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceMenu struct {
	SafeDB store.Queryer
	Store  store.MenuStore
}

// 권한별 메뉴 리스트
func (s *ServiceMenu) GetMenu(ctx context.Context, roles []string) (entity.MenuRes, error) {
	parentList, err := s.Store.GetParentMenu(ctx, s.SafeDB, roles)
	if err != nil {
		return entity.MenuRes{}, fmt.Errorf("service_menu/GetParentMenu err: %w", err)
	}

	childList, err := s.Store.GetChildMenu(ctx, s.SafeDB, roles)
	if err != nil {
		return entity.MenuRes{}, fmt.Errorf("service_menu/GetChildMenu err: %w", err)
	}

	list := entity.MenuRes{
		Parent: parentList,
		Child:  childList,
	}

	return list, nil
}
