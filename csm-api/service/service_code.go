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
	sqls, err := s.Store.GetCodeList(ctx, s.DB, pCode)
	if err != nil {
		return nil, fmt.Errorf("service_code/GetCodeList err: %w", err)
	}
	codes := &entity.Codes{}
	codes.ToCodes(sqls)

	return codes, nil
}
