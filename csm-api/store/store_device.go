package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"database/sql"
	"errors"
	"fmt"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// func: 근태인식기 전체 조회
// @param
// - page entity.PageSql: 현재페이지 번호, 리스트 목록 개수
func (r *Repository) GetDeviceList(ctx context.Context, db Queryer, page entity.PageSql) (*entity.DeviceSqls, error) {
	sqls := entity.DeviceSqls{}

	query := `SELECT *
				FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						SELECT 
							t1.DNO,
							t1.SNO,
							t2.SITE_NM,
							t1.DEVICE_SN,
							t1.DEVICE_NM,
							t1.IS_USE,
							t1.REG_DATE AS REG_DATE,
							t1.MOD_DATE AS MOD_DATE
						FROM IRIS_DEVICE_SET t1
						LEFT OUTER JOIN IRIS_SITE_SET t2 ON t1.SNO = t2.SNO
						WHERE t1.IS_USE = 'Y'
						ORDER BY t1.REG_DATE DESC
					) sorted_data
					WHERE ROWNUM <= :1
				)
				WHERE RNUM > :2`

	if err := db.SelectContext(ctx, &sqls, query, page.EndNum, page.StartNum); err != nil {
		return nil, fmt.Errorf("GetDeviceList err: %v", err)
	}

	return &sqls, nil
}

// func: 근태인식기 전체 개수 조회
// @param
// -
func (r *Repository) GetDeviceListCount(ctx context.Context, db Queryer) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM IRIS_DEVICE_SET WHERE IS_USE = 'Y'`

	if err := db.GetContext(ctx, &count, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("GetDeviceListCount fail: %w", err)
	}
	return count, nil
}

// func: 근태인식기 추가
// @param
// - device entity.DeviceSql: SNO, DEVICE_SN, DEVICE_NM, ETC, IS_USE, REG_USER
func (r *Repository) AddDevice(ctx context.Context, db Beginner, device entity.DeviceSql) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Failed to begin transaction: %v", err)
	}

	agent := utils.GetAgent()

	query := `
				INSERT INTO IRIS_DEVICE_SET(
					DNO, 
					SNO, 
					DEVICE_SN, 
					DEVICE_NM, 
					ETC, 
					IS_USE, 
					REG_DATE, 
					REG_AGENT,
					REG_USER
				) VALUES (
					SEQ_IRIS_DEVICE_SET.NEXTVAL,
					:1,
				    :2,
				    :3,
				    :4,
				    :5,
				    SYSDATE,
				    :6,
				    :7    
				)`
	_, err = tx.ExecContext(ctx, query, device.Sno, device.DeviceSn, device.DeviceNm, device.Etc, device.IsUse, agent, device.RegUser)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("AddDevice fail: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

// func: 근태인식기 수정
// @param
// - device entity.DeviceSql: DNO, SNO, DEVICE_SN, DEVICE_NM, ETC, IS_USE, MOD_USER
func (r *Repository) ModifyDevice(ctx context.Context, db Beginner, device entity.DeviceSql) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Failed to begin transaction: %v", err)
	}

	agent := utils.GetAgent()

	query := `
				UPDATE IRIS_DEVICE_SET 
				SET 
					SNO = :1, 
					DEVICE_SN = :2, 
					DEVICE_NM = :3, 
					ETC = :4, 
					IS_USE = :5, 
					MOD_DATE = SYSDATE,
					MOD_USER = :6,
					MOD_AGENT = :7 
				WHERE DNO = :8`

	_, err = tx.ExecContext(ctx, query, device.Sno, device.DeviceSn, device.DeviceNm, device.Etc, device.IsUse, device.ModUser, agent, device.Dno)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("ModifyDevice fail: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

// func: 근태인식기 삭제
// @param
// - dno sql.NullInt64: 홍채인식기 고유번호
func (r *Repository) RemoveDevice(ctx context.Context, db Beginner, dno sql.NullInt64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Failed to begin transaction: %v", err)
	}

	query := `DELETE FROM IRIS_DEVICE_SET WHERE DNO = :1`

	_, err = tx.ExecContext(ctx, query, dno)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("RemoveDevice fail: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}
