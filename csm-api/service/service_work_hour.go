package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/txutil"
	"csm-api/utils"
)

type ServiceWorkHour struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.WorkHourStore
}

// 특정 프로젝트 및 근로자의 공수 계산: jno는 필수, ids는 없으면 jno의 모든 근로자 계산 있으면 해당 id의 근로자만 계산
func (s *ServiceWorkHour) ModifyWorkHourByJno(ctx context.Context, jno int64, user entity.Base, ids []string) (err error) {
	tx, err := txutil.BeginTxWithMode(ctx, s.SafeTDB, false)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer txutil.DeferTx(tx, &err)

	err = s.Store.ModifyWorkHourByJno(ctx, tx, jno, user, ids)
	if err != nil {
		return utils.CustomErrorf(err)
	}
	return
}

// 출퇴근이 둘다 있는 모든 근로자의 공수 계산
func (s *ServiceWorkHour) ModifyWorkHour(ctx context.Context, user entity.Base) (err error) {
	tx, err := txutil.BeginTxWithMode(ctx, s.SafeTDB, false)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer txutil.DeferTx(tx, &err)

	err = s.Store.ModifyWorkHour(ctx, tx, user)
	if err != nil {
		return utils.CustomErrorf(err)
	}
	return
}
