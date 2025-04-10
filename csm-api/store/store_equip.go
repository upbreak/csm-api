package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

func (r *Repository) GetEquipList(ctx context.Context, db Queryer) (entity.EquipTemps, error) {
	list := entity.EquipTemps{}

	query := `
			SELECT 
			    T1.SNO, 
			    T1.JNO, 
			    NVL(T2.CNT, 0) AS CNT, 
			    T3.JOB_NAME
			FROM IRIS_SITE_JOB T1
			LEFT JOIN IRIS_EQUIP_TEMP T2 ON T1.SNO = T2.SNO AND T1.JNO = T2.JNO 
			LEFT JOIN S_JOB_INFO T3 ON T1.JNO = T3.JNO`

	if err := db.SelectContext(ctx, &list, query); err != nil {
		return list, fmt.Errorf("GetEquipList fail: %w", err)
	}
	return list, nil
}

func (r *Repository) MergeEquipCnt(ctx context.Context, tx Execer, equips entity.EquipTemps) error {
	query := `
			MERGE INTO IRIS_EQUIP_TEMP T1
			USING(
				SELECT
					:1 AS SNO,
					:2 AS JNO,
					:3 AS CNT
				FROM DUAL
			) T2 
			ON (
				T1.SNO = T2.SNO
				AND T1.JNO = T2.JNO
			)
			WHEN MATCHED THEN
				UPDATE SET
					T1.CNT = T2.CNT
				WHERE T1.SNO = T2.SNO
				AND T1.JNO = T2.JNO
			WHEN NOT MATCHED THEN
				INSERT (SNO, JNO, CNT)
				VALUES (T2.SNO, T2.JNO, T2.CNT)`

	for _, equip := range equips {
		if _, err := tx.ExecContext(ctx, query, equip.Sno, equip.Jno, equip.Cnt); err != nil {
			//TODO: 에러 아카이브
			return fmt.Errorf("MergeEquipCnt ExecContext fail: %v", err)
		}
	}

	return nil
}
