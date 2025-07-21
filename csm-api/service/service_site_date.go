package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/txutil"
	"csm-api/utils"
)

type ServiceSiteDate struct {
	DB    store.Queryer
	TDB   store.Beginner
	Store store.SiteDateStore
}

// 현장 날짜 테이블 조회
//
// @param sno: 현장 고유번호
func (s *ServiceSiteDate) GetSiteDateData(ctx context.Context, sno int64) (*entity.SiteDate, error) {
	siteDate, err := s.Store.GetSiteDateData(ctx, s.DB, sno)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return siteDate, nil
}

// 현장 날짜 테이블 수정
//
// @param
// - sno: 현장고유번호
// - siteDate: 현장 시간 (opening_date, closing_plan_date, closing_forecast_date, closing_actual_date)
func (s *ServiceSiteDate) ModifySiteDate(ctx context.Context, sno int64, siteDate entity.SiteDate) (err error) {
	tx, cleanup, err := txutil.BeginTxWithCleanMode(ctx, s.TDB, false)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer func() {
		txutil.DeferTx(tx, &err)
		cleanup()
	}()

	if err := s.Store.ModifySiteDate(ctx, tx, sno, siteDate); err != nil {
		return utils.CustomErrorf(err)
	}

	return
}
