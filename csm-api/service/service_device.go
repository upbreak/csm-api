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
	SafeQueryer  store.Queryer
	SafeBeginner store.Beginner
	Store        store.DeviceStore
}

// func: 근태인식기 조회
// @param
// - page entity.PageSql: 현재페이지 번호, 리스트 목록 개수
func (s *ServiceDevice) GetDeviceList(ctx context.Context, page entity.Page, search entity.Device, retry string) (*entity.Devices, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_device/GetDeviceList err: %w", err)
	}

	list, err := s.Store.GetDeviceList(ctx, s.SafeQueryer, pageSql, search, retry)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_device/GetDeviceList err: %w", err)
	}

	return list, nil
}

// func: 근태인식기 전체 개수 조회
// @param
// -
func (s *ServiceDevice) GetDeviceListCount(ctx context.Context, search entity.Device, retry string) (int, error) {
	count, err := s.Store.GetDeviceListCount(ctx, s.SafeQueryer, search, retry)
	if err != nil {
		//TODO: 에러 아카이브
		return 0, fmt.Errorf("service_device/GetDeviceListCount err: %w", err)
	}

	return count, nil
}

// func: 근태인식기 추가
// @param
// - device entity.Device: SNO, DEVICE_SN, DEVICE_NM, ETC, IS_USE, REG_USER
func (s *ServiceDevice) AddDevice(ctx context.Context, device entity.Device) (err error) {
	tx, err := s.SafeBeginner.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_device/AddDevice BeginTx fail err: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_device/AddDevice Rollback err: %w", rollbackErr)
			}
			err = fmt.Errorf("service_device/AddDevice err: %w", err)
		} else {
			if err = tx.Commit(); err != nil {
				err = fmt.Errorf("service_device/AddDevice Commit err: %w", err)
			}
		}
	}()

	if err = s.Store.AddDevice(ctx, tx, device); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_device/AddDevice err: %w", err)
	}
	return
}

// func: 근태인식기 수정
// @param
// - device entity.DeviceSql: DNO, SNO, DEVICE_SN, DEVICE_NM, ETC, IS_USE, MOD_USER
func (s *ServiceDevice) ModifyDevice(ctx context.Context, device entity.Device) (err error) {
	tx, err := s.SafeBeginner.BeginTx(ctx, nil)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_device/ModifyDevice BeginTx fail err: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				//TODO: 에러 아카이브
				err = fmt.Errorf("service_device/ModifyDevice Rollback err: %w", rollbackErr)
			}
			err = fmt.Errorf("service_device/ModifyDevice Rollback err: %w", err)
		} else {
			if err = tx.Commit(); err != nil {
				err = fmt.Errorf("service_device/ModifyDevice Commit err: %w", err)
			}
		}
	}()

	if err = s.Store.ModifyDevice(ctx, tx, device); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_device/UpdateDevice err: %w", err)
	}
	return
}

// func: 근태인식기 삭제
// @param
// - dno int64: 홍채인식기 고유번호
func (s *ServiceDevice) RemoveDevice(ctx context.Context, dno int64) (err error) {
	tx, err := s.SafeBeginner.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_device/RemoveDevice BeginTx fail err: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_device/RemoveDevice Rollback err: %w", rollbackErr)
			}
			err = fmt.Errorf("service_device/RemoveDevice Rollback err: %w", err)
		} else {
			if err = tx.Commit(); err != nil {
				err = fmt.Errorf("service_device/RemoveDevice Commit err: %w", err)
			}
		}
	}()

	var dnoSql sql.NullInt64
	if dno != 0 {
		dnoSql = sql.NullInt64{Valid: true, Int64: dno}
	} else {
		dnoSql = sql.NullInt64{Valid: false}
	}

	if err = s.Store.RemoveDevice(ctx, tx, dnoSql); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_device/RemoveDevice err: %w", err)
	}
	return
}
