package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
)

type ServiceMenu struct {
	SafeDB store.Queryer
	Store  store.MenuStore
}

// 권한별 메뉴 리스트
func (s *ServiceMenu) GetMenu(ctx context.Context, roles []string) (entity.MenuRes, error) {
	parentList, err := s.Store.GetParentMenu(ctx, s.SafeDB, roles)
	if err != nil {
		return entity.MenuRes{}, utils.CustomErrorf(err)
	}

	childList, err := s.Store.GetChildMenu(ctx, s.SafeDB, roles)
	if err != nil {
		return entity.MenuRes{}, utils.CustomErrorf(err)
	}

	list := entity.MenuRes{
		Parent: parentList,
		Child:  childList,
	}

	return list, nil
}
