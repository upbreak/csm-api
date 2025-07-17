package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/txutil"
	"csm-api/utils"
)

type ServiceEquip struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.EquipStore
}

func (s *ServiceEquip) GetEquipList(ctx context.Context) (entity.EquipTemps, error) {
	list, err := s.Store.GetEquipList(ctx, s.SafeDB)
	if err != nil {
		return entity.EquipTemps{}, utils.CustomErrorf(err)
	}
	return list, nil
}

func (s *ServiceEquip) MergeEquipCnt(ctx context.Context, equips entity.EquipTemps) (err error) {
	tx, err := txutil.BeginTxWithMode(ctx, s.SafeTDB, false)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer txutil.DeferTx(tx, &err)

	if err = s.Store.MergeEquipCnt(ctx, tx, equips); err != nil {
		return utils.CustomErrorf(err)
	}
	return
}
