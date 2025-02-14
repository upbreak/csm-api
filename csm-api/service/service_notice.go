package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceNotice struct {
	DB    store.Queryer
	Store store.NoticeStore
}

func (s *ServiceNotice) GetNoticeList(ctx context.Context, page entity.Page) (*entity.Notices, error) {

	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		return nil, fmt.Errorf("service_notice/GetNoticeList err : %w", err)
	}

	noticeSqls, err := s.Store.GetNoticeList(ctx, s.DB, pageSql)
	if err != nil {
		return &entity.Notices{}, fmt.Errorf("fail to list notice: %w", err)
	}

	notices := &entity.Notices{}
	notices.ToNotices(noticeSqls)

	return notices, nil
}

func (s *ServiceNotice) GetNoticeListCount(ctx context.Context) (int, error) {
	count, err := s.Store.GetNoticeListCount(ctx, s.DB)
	if err != nil {
		return 0, fmt.Errorf("service_notice/GetNoticeListCount err : %w", err)
	}

	return count, nil

}
