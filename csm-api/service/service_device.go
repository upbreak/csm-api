package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"database/sql"
	"encoding/json"
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
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.DeviceStore
}

// func: 근태인식기 조회
// @param
// - page entity.PageSql: 현재페이지 번호, 리스트 목록 개수
func (s *ServiceDevice) GetDeviceList(ctx context.Context, page entity.Page, search entity.Device, retry string) (*entity.Devices, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	list, err := s.Store.GetDeviceList(ctx, s.SafeDB, pageSql, search, retry)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return list, nil
}

// func: 근태인식기 전체 개수 조회
// @param
// -
func (s *ServiceDevice) GetDeviceListCount(ctx context.Context, search entity.Device, retry string) (int, error) {
	count, err := s.Store.GetDeviceListCount(ctx, s.SafeDB, search, retry)
	if err != nil {
		return 0, utils.CustomErrorf(err)
	}

	return count, nil
}

// func: 근태인식기 추가
// @param
// - device entity.Device: SNO, DEVICE_SN, DEVICE_NM, ETC, IS_USE, REG_USER
func (s *ServiceDevice) AddDevice(ctx context.Context, device entity.Device) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer utils.DeferTx(tx, &err)

	if err = s.Store.AddDevice(ctx, tx, device); err != nil {
		return utils.CustomErrorf(err)
	}
	return
}

// func: 근태인식기 수정
// @param
// - device entity.DeviceSql: DNO, SNO, DEVICE_SN, DEVICE_NM, ETC, IS_USE, MOD_USER
func (s *ServiceDevice) ModifyDevice(ctx context.Context, device entity.Device) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer utils.DeferTx(tx, &err)

	if err = s.Store.ModifyDevice(ctx, tx, device); err != nil {
		return utils.CustomErrorf(err)
	}
	return
}

// func: 근태인식기 삭제
// @param
// - dno int64: 홍채인식기 고유번호
func (s *ServiceDevice) RemoveDevice(ctx context.Context, dno int64) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer utils.DeferTx(tx, &err)

	var dnoSql sql.NullInt64
	if dno != 0 {
		dnoSql = sql.NullInt64{Valid: true, Int64: dno}
	} else {
		dnoSql = sql.NullInt64{Valid: false}
	}

	if err = s.Store.RemoveDevice(ctx, tx, dnoSql); err != nil {
		return utils.CustomErrorf(err)
	}
	return
}

// func: 근태인식기 미등록장치 확인
// @param
func (s *ServiceDevice) GetCheckRegisteredDevices(ctx context.Context) ([]string, error) {

	// 당일 iris_recd_log 가져오기
	devices, err := s.Store.GetDeviceLog(ctx, s.SafeDB)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	var log entity.RecdLog
	deviceList := make(map[string]int)
	// iris_recd_log 테이블의 iris_data값인 json에 들어온 deviceName을 파싱해서 deviceList 얻기
	for _, device := range *devices {
		if err = json.Unmarshal([]byte(device.IrisData.String), &log); err != nil {
			return nil, utils.CustomErrorf(err)
		}

		// 중복 제거
		_, ok := deviceList[log.DeviceName.String]
		if !ok {
			deviceList[log.DeviceName.String] = 1
		}
	}

	// 미등록 확인.
	var respond []string
	var check int
	for deviceName, _ := range deviceList {
		check, err = s.Store.GetCheckRegistered(ctx, s.SafeDB, deviceName)
		if err != nil {
			return nil, utils.CustomErrorf(err)
		}
		if check == 0 {
			respond = append(respond, deviceName)
		}
	}

	return respond, nil
}
