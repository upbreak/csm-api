package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
)

type ServiceCode struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.CodeStore
}

func (s *ServiceCode) GetCodeList(ctx context.Context, pCode string) (*entity.Codes, error) {
	list, err := s.Store.GetCodeList(ctx, s.SafeDB, pCode)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return list, nil
}

// func: 코드트리 조회
// @param
// -
func (s *ServiceCode) GetCodeTree(ctx context.Context, pCode string) (*entity.CodeTrees, error) {

	// 코드리스트 조회
	codes, err := s.Store.GetCodeTree(ctx, s.SafeDB, pCode)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	// 트리구조로 반환
	trees, err := entity.ConvertCodesToCodeTree(*codes, pCode)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return &trees, nil

}

// func: 코드트리 수정 및 저장
// @param
// -
func (s *ServiceCode) MergeCode(ctx context.Context, code entity.Code) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer utils.DeferTx(tx, &err)

	if err = s.Store.MergeCode(ctx, tx, code); err != nil {
		return utils.CustomErrorf(err)
	}

	return
}

// func: 코드 삭제
// @param
// -
func (s *ServiceCode) RemoveCode(ctx context.Context, idx int64) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer utils.DeferTx(tx, &err)

	if err = s.Store.RemoveCode(ctx, tx, idx); err != nil {
		return utils.CustomErrorf(err)
	}

	return
}

// func: 코드순서 변경
// @param
// -
func (s *ServiceCode) ModifySortNo(ctx context.Context, codeSorts entity.CodeSorts) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer utils.DeferTx(tx, &err)

	for _, codeSort := range codeSorts {
		if err = s.Store.ModifySortNo(ctx, tx, *codeSort); err != nil {
			return utils.CustomErrorf(err)
		}
	}

	return
}

// func: 코드 중복 검사
// @param
// - code
func (s *ServiceCode) DuplicateCheckCode(ctx context.Context, code string) (bool, error) {

	count, err := s.Store.DuplicateCheckCode(ctx, s.SafeDB, code)
	if err != nil {
		return false, utils.CustomErrorf(err)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil

}
