package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceNotice struct {
	DB    store.Queryer
	TDB   store.Beginner
	Store store.NoticeStore
}

// func: 공지사항 전체 조회
// @param
// - page entity.PageSql : 현재 페이지번호, 리스트 목록 개수
func (s *ServiceNotice) GetNoticeList(ctx context.Context, page entity.Page, search entity.Notice) (*entity.Notices, error) {

	pageSql := entity.PageSql{}
	searchSql := &entity.NoticeSql{}
	pageSql, err := pageSql.OfPageSql(page)
	searchSql = searchSql.OfNoticeSql(search)

	if err != nil {
		return nil, fmt.Errorf("service_notice/GetNoticeList err : %w", err)
	}

	noticeSqls, err := s.Store.GetNoticeList(ctx, s.DB, pageSql, *searchSql)
	if err != nil {
		return &entity.Notices{}, fmt.Errorf("fail to list notice: %w", err)
	}

	notices := &entity.Notices{}
	notices.ToNotices(noticeSqls)

	return notices, nil
}

// func: 공지사항 전체 개수 조회
// @param
// -
func (s *ServiceNotice) GetNoticeListCount(ctx context.Context, search entity.Notice) (int, error) {
	searchSql := &entity.NoticeSql{}
	searchSql = searchSql.OfNoticeSql(search)

	count, err := s.Store.GetNoticeListCount(ctx, s.DB, *searchSql)
	if err != nil {
		return 0, fmt.Errorf("service_notice/GetNoticeListCount err : %w", err)
	}

	return count, nil

}

// func: 공지사항 추가
// @param
// - notice entity.Notice: SNO, TITLE, CONTENT, SHOW_YN, REG_UNO, REG_USER
func (s *ServiceNotice) AddNotice(ctx context.Context, notice entity.Notice) error {
	noticeSql := &entity.NoticeSql{}
	noticeSql = noticeSql.OfNoticeSql(notice)

	if err := s.Store.AddNotice(ctx, s.TDB, *noticeSql); err != nil {
		return fmt.Errorf("service_notice/AddNotice err : %w", err)
	}

	return nil
}

// func: 공지사항 수정
// @param
// -notice entity.Notice: IDX, SNO, TITLE, CONTENT, SHOW_YN, MOD_UNO, MOD_USER
func (s *ServiceNotice) ModifyNotice(ctx context.Context, notice entity.Notice) error {
	noticeSql := &entity.NoticeSql{}
	noticeSql = noticeSql.OfNoticeSql(notice)

	if err := s.Store.ModifyNotice(ctx, s.TDB, *noticeSql); err != nil {
		return fmt.Errorf("service_notice/ModifyNotice err: %w", err)
	}

	return nil
}

// func: 공지사항 삭제
// @param
// - IDX: 공지사항 인덱스
func (s *ServiceNotice) RemoveNotice(ctx context.Context, idx int64) error {
	var idxSql entity.NoticeID

	if idx != 0 {
		idxSql = entity.NoticeID(idx)
	} else {
		return fmt.Errorf("idx parameter is missing")
	}

	if err := s.Store.RemoveNotice(ctx, s.TDB, idxSql); err != nil {
		return fmt.Errorf("service_notice/RemomveNotice err: %w", err)
	}

	return nil
}
