package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

func (r *Repository) MergeEquipCnt(ctx context.Context, db Beginner, equips entity.EquipTemps) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("MergeEquipCnt BeginTx fail: %v", err)
	}

	query := `
			MERGE INTO IRIS_EQUIP_TEMP T1
			USING(
				SELECT
					:1 AS SNO,
					:2 AS JNO,
					:3 AS CNT,
					TRUNC(SYSDATE) AS RECORD_DATE,
				    :4 AS REG_USER
				FROM DUAL
			) T2 
			ON (
				T1.SNO = T2.SNO
				AND T1.JNO = T2.JNO
				AND T1.RECORD_DATE = T2.RECORD_DATE
			)
			WHEN MATCHED THEN
				UPDATE SET
					T1.CNT = T2.CNT,
			    	T1.REG_USER = T2.REG_USER
				WHERE T1.SNO = T2.SNO
				AND T1.JNO = T2.JNO
				AND T1.RECORD_DATE = T2.RECORD_DATE
			WHEN NOT MATCHED THEN
				INSERT (SNO, JNO, CNT, RECORD_DATE, REG_USER)
				VALUES (T2.SNO, T2.JNO, T2.CNT, T2.RECORD_DATE, T2.REG_USER)`

	for _, equip := range equips {
		if _, err = tx.QueryContext(ctx, query, equip.Sno, equip.Jno, equip.Cnt, equip.RegUser); err != nil {
			origErr := err
			err = tx.Rollback()
			if err != nil {
				//TODO: 에러 아카이브 처리
				return fmt.Errorf("MergeEquipCnt Rollback fail: %v\n", err)
			}
			//TODO: 에러 아카이브
			return fmt.Errorf("MergeEquipCnt QueryContext fail: %v", origErr)
		}
	}

	if err = tx.Commit(); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("MergeEquipCnt Commit fail: %v\n", err)
	}
	return nil
}
