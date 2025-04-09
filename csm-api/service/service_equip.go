package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceEquip struct {
	DB    store.Queryer
	TDB   store.Beginner
	Store store.EquipStore
}

func (s *ServiceEquip) GetEquipList(ctx context.Context) (entity.EquipTemps, error) {
	list, err := s.Store.GetEquipList(ctx, s.DB)
	if err != nil {
		return entity.EquipTemps{}, fmt.Errorf("service;GetEquipList fail: %w", err)
	}
	return list, nil
}

func (s *ServiceEquip) MergeEquipCnt(ctx context.Context, equips entity.EquipTemps) error {
	if err := s.Store.MergeEquipCnt(ctx, s.TDB, equips); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service;MergeEquipCnt failed: %v", err)
	}
	return nil
}
