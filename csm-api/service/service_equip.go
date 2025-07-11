package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceEquip struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.EquipStore
}

func (s *ServiceEquip) GetEquipList(ctx context.Context) (entity.EquipTemps, error) {
	list, err := s.Store.GetEquipList(ctx, s.SafeDB)
	if err != nil {
		return entity.EquipTemps{}, fmt.Errorf("service;GetEquipList fail: %w", err)
	}
	return list, nil
}

func (s *ServiceEquip) MergeEquipCnt(ctx context.Context, equips entity.EquipTemps) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service;MergeEquipCnt fail: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service;MergeEquipCnt panic: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service;MergeEquipCnt rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service;MergeEquipCnt commit error: %w", commitErr)
			}
		}
	}()
	if err = s.Store.MergeEquipCnt(ctx, tx, equips); err != nil {
		return fmt.Errorf("service;MergeEquipCnt failed: %v", err)
	}
	return
}
