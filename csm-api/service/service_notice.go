package service

import (
	"context"
	"csm-api/auth"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"github.com/guregu/null"
)

type ServiceNotice struct {
	SafeDB    store.Queryer
	SafeTDB   store.Beginner
	Store     store.NoticeStore
	UserStore store.UserStore
}

// func: 공지사항 전체 조회
// @param
// - page entity.PageSql : 현재 페이지번호, 리스트 목록 개수
func (s *ServiceNotice) GetNoticeList(ctx context.Context, page entity.Page, search entity.Notice) (*entity.Notices, error) {

	// 사용자 정보 가져오기
	role, _ := auth.GetContext(ctx, auth.Role{})
	unoString, _ := auth.GetContext(ctx, auth.Uno{})
	uno := utils.ParseNullInt(unoString)

	// 권한 조회
	list, err := s.UserStore.GetAuthorizationList(ctx, s.SafeDB, "/notice")
	if err != nil {
		// TODO: 에러 아카이브 처리
		return nil, fmt.Errorf("service_notice/GetNoticeList GetAuthorizationList err: %w", err)
	}

	var roleInt int
	if entity.AuthorizationCheck(*list, role) {
		roleInt = 1
	} else {
		roleInt = 0
	}

	// 페이지 변환
	pageSql := entity.PageSql{}
	pageSql, err = pageSql.OfPageSql(page)

	if err != nil {
		//TODO: 에러 아카이브 처리
		return nil, fmt.Errorf("service_notice/GetNoticeList err : %w", err)
	}

	notices, err := s.Store.GetNoticeList(ctx, s.SafeDB, uno, roleInt, pageSql, search)
	if err != nil {
		//TODO: 에러 아카이브 처리
		return &entity.Notices{}, fmt.Errorf("fail to list notice: %w", err)
	}

	return notices, nil
}

// func: 공지사항 전체 개수 조회
// @param
// -
func (s *ServiceNotice) GetNoticeListCount(ctx context.Context, search entity.Notice) (int, error) {

	role, _ := auth.GetContext(ctx, auth.Role{})
	unoString, _ := auth.GetContext(ctx, auth.Uno{})
	uno := utils.ParseNullInt(unoString)

	// 권한 조회
	list, err := s.UserStore.GetAuthorizationList(ctx, s.SafeDB, "/notice")
	if err != nil {
		// TODO: 에러 아카이브 처리
		return 0, fmt.Errorf("service_notice/GetNoticeList GetAuthorizationList err: %w", err)
	}

	var roleInt int
	if entity.AuthorizationCheck(*list, role) {
		roleInt = 1
	} else {
		roleInt = 0
	}

	count, err := s.Store.GetNoticeListCount(ctx, s.SafeDB, uno, roleInt, search)
	if err != nil {
		//TODO: 에러 아카이브 처리
		return 0, fmt.Errorf("service_notice/GetNoticeListCount err : %w", err)
	}

	return count, nil

}

// func: 공지사항 추가
// @param
// - notice entity.Notice: JNO, TITLE, CONTENT, SHOW_YN, PERIOD_CODE, REG_UNO, REG_USER
func (s *ServiceNotice) AddNotice(ctx context.Context, notice entity.Notice) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_notice/AddNotice err : %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_notice/AddNotice rollback err : %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_notice/AddNotice commit err : %w", commitErr)
			}
		}
	}()
	if err = s.Store.AddNotice(ctx, tx, notice); err != nil {
		//TODO: 에러 아카이브 처리
		return fmt.Errorf("service_notice/AddNotice err : %w", err)
	}

	return
}

// func: 공지사항 수정
// @param
// -notice entity.Notice: IDX, JNO, TITLE, CONTENT, SHOW_YN, MOD_UNO, MOD_USER
func (s *ServiceNotice) ModifyNotice(ctx context.Context, notice entity.Notice) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_notice/ModifyNotice err : %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_notice/ModifyNotice rollback err : %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_notice/ModifyNotice commit err : %w", commitErr)
			}
		}
	}()

	if err = s.Store.ModifyNotice(ctx, tx, notice); err != nil {
		//TODO: 에러 아카이브 처리
		return fmt.Errorf("service_notice/ModifyNotice err: %w", err)
	}

	return
}

// func: 공지사항 삭제
// @param
// - IDX: 공지사항 인덱스
func (s *ServiceNotice) RemoveNotice(ctx context.Context, idx null.Int) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_notice/RemoveNotice err : %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_notice/RemoveNotice rollback err : %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_notice/RemoveNotice commit err : %w", commitErr)
			}
		}
	}()

	if err = s.Store.RemoveNotice(ctx, tx, idx); err != nil {
		//TODO: 에러 아카이브 처리
		return fmt.Errorf("service_notice/RemomveNotice err: %w", err)
	}

	return
}
