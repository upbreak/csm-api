package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"database/sql"
	"fmt"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일: 2025-02-21
 * @modifiedBy 최종 수정자: 정지영
 * @modified description
 * - 검색 및 정렬 조건 추가
 */

// struct: 근태인식기 서비스 구조체
type ServiceDevice struct {
	DB    store.Queryer
	TDB   store.Beginner
	Store store.DeviceStore
}

// func: 근태인식기 조회
// @param
// - page entity.PageSql: 현재페이지 번호, 리스트 목록 개수
func (s *ServiceDevice) GetDeviceList(ctx context.Context, page entity.Page, search entity.Device, retry string) (*entity.Devices, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	searchSql := entity.DeviceSql{}
	searchSql = *searchSql.OfDeviceSql(search)

	if err != nil {
		return nil, fmt.Errorf("service_device/GetDeviceList err: %w", err)
	}
	dbList, err := s.Store.GetDeviceList(ctx, s.DB, pageSql, searchSql, retry)
	if err != nil {
		return nil, fmt.Errorf("service_device/GetDeviceList err: %w", err)
	}

	list := &entity.Devices{}
	list.ToDevices(dbList)

	return list, nil
}

// func: 근태인식기 전체 개수 조회
// @param
// -
func (s *ServiceDevice) GetDeviceListCount(ctx context.Context, search entity.Device, retry string) (int, error) {
	searchSql := entity.DeviceSql{}
	searchSql = *searchSql.OfDeviceSql(search)

	count, err := s.Store.GetDeviceListCount(ctx, s.DB, searchSql, retry)
	if err != nil {
		return 0, fmt.Errorf("service_device/GetDeviceListCount err: %w", err)
	}

	return count, nil
}

// func: 근태인식기 추가
// @param
// - device entity.Device: SNO, DEVICE_SN, DEVICE_NM, ETC, IS_USE, REG_USER
func (s *ServiceDevice) AddDevice(ctx context.Context, device entity.Device) error {
	deviceSql := &entity.DeviceSql{}
	deviceSql = deviceSql.OfDeviceSql(device)
	if err := s.Store.AddDevice(ctx, s.TDB, *deviceSql); err != nil {
		return fmt.Errorf("service_device/AddDevice err: %w", err)
	}
	return nil
}

// func: 근태인식기 수정
// @param
// - device entity.DeviceSql: DNO, SNO, DEVICE_SN, DEVICE_NM, ETC, IS_USE, MOD_USER
func (s *ServiceDevice) ModifyDevice(ctx context.Context, device entity.Device) error {
	deviceSql := &entity.DeviceSql{}
	deviceSql = deviceSql.OfDeviceSql(device)
	if err := s.Store.ModifyDevice(ctx, s.TDB, *deviceSql); err != nil {
		return fmt.Errorf("service_device/UpdateDevice err: %w", err)
	}
	return nil
}

// func: 근태인식기 삭제
// @param
// - dno int64: 홍채인식기 고유번호
func (s *ServiceDevice) RemoveDevice(ctx context.Context, dno int64) error {
	var dnoSql sql.NullInt64
	if dno != 0 {
		dnoSql = sql.NullInt64{Valid: true, Int64: dno}
	} else {
		dnoSql = sql.NullInt64{Valid: false}
	}
	if err := s.Store.RemoveDevice(ctx, s.TDB, dnoSql); err != nil {
		return fmt.Errorf("service_device/RemoveDevice err: %w", err)
	}
	return nil
}
