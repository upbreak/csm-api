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
