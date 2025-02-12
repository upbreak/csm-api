package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ServiceDevice struct {
	DB    store.Queryer
	Store store.DeviceStore
}

func (s *ServiceDevice) GetDeviceList(ctx context.Context, page entity.Page) (*entity.Devices, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		return nil, fmt.Errorf("service_device/GetDeviceList err: %w", err)
	}
	dbList, err := s.Store.GetDeviceList(ctx, s.DB, pageSql)
	if err != nil {
		return nil, fmt.Errorf("service_device/GetDeviceList err: %w", err)
	}

	list := &entity.Devices{}
	list.ToDevices(dbList)

	return list, nil
}

func (s *ServiceDevice) GetDeviceListCount(ctx context.Context) (int, error) {
	count, err := s.Store.GetDeviceListCount(ctx, s.DB)
	if err != nil {
		return 0, fmt.Errorf("service_device/GetDeviceListCount err: %w", err)
	}

	return count, nil
}
