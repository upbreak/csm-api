package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceCode struct {
	DB    store.Queryer
	Store store.CodeStore
}

func (s *ServiceCode) GetCodeList(ctx context.Context, pCode string) (*entity.Codes, error) {
	list, err := s.Store.GetCodeList(ctx, s.DB, pCode)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_code/GetCodeList err: %w", err)
	}

	return list, nil
}

// 코드트리 조회
func (s *ServiceCode) GetCodeTree(ctx context.Context) (*entity.CodeTrees, error) {

	// 코드리스트 조회
	codes, err := s.Store.GetCodeTree(ctx, s.DB)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_code/GetCodeSetList err: %w", err)
	}

	// 트리구조로 반환
	trees := convertCodesToCodeTree(*codes, "")

	return &trees, nil

}

// 코드를 코드트리 구조로 변환
func convertCodesToCodeTree(codes entity.Codes, pCode string) entity.CodeTrees {
	codeTrees := entity.CodeTrees{}

	for i := 0; i < len(codes); i++ {

		// 부모가 다른 경우 넘기기
		if codes[i].PCode.Valid && codes[i].PCode.String != pCode {
			continue
		}

		// code의 값 tree 형태로 변환
		codeTree := &entity.CodeTree{}

		codeTree.IDX = codes[i].IDX
		codeTree.Level = codes[i].Level
		codeTree.Code = codes[i].Code

		// root는 root값 입력하기
		if codes[i].PCode.String == "" {
			codeTree.PCode.String = "root"
			codeTree.PCode.Valid = true
		} else {
			codeTree.PCode = codes[i].PCode
		}

		// 하위 code 넣기
		codeTree.Children = &entity.CodeTrees{}
		child := convertCodesToCodeTree(codes[i+1:], codes[i].Code.String) // 현재 코드가 다음 레벨의 부모코드
		codeTree.Children = &child
		if len(child) > 0 {
			codeTree.Expand.Bool = true
			codeTree.Expand.Valid = true
		} else {
			codeTree.Expand.Bool = false
			codeTree.Expand.Valid = true
		}
		codeTree.CodeSet = codes[i]
		codeTrees = append(codeTrees, codeTree)
	}

	return codeTrees
}
