package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"fmt"
)

// TBM 엑셀 정보 저장
func (r *Repository) AddTbmExcel(ctx context.Context, tx Execer, tbms []entity.Tbm) error {
	agent := utils.GetAgent()

	query := `
		INSERT INTO IRIS_TBM_SET(JNO, DEPARTMENT, DISC_NAME, USER_NM, TBM_DATE, REG_DATE, REG_USER, REG_UNO, REG_AGENT)
		VALUES(:1, :2, :3, :4, :5, SYSDATE, :7, :8, :9)`

	for _, tbm := range tbms {
		if _, err := tx.ExecContext(ctx, query, tbm.Jno, tbm.Department, tbm.DiscName, tbm.UserNm, tbm.TbmDate, tbm.RegUser, tbm.RegUno, agent); err != nil {
			return fmt.Errorf("AddTbmExcel: %w", err)
		}
	}
	return nil
}

// TBM 엑셀 정보 사용 안함
func (r *Repository) ModifyTbmExcel(ctx context.Context, tx Execer, tbm entity.Tbm) error {
	agent := utils.GetAgent()

	query := `
		UPDATE IRIS_TBM_SET
		SET
		    IS_USE = 'N',
			MOD_DATE = SYSDATE,
			MOD_USER = :1,
			MOD_UNO = :2,
			MOD_AGENT = :3
		WHERE JNO = :4 
		AND TRUNC(TBM_DATE) = TRUNC(:5)
		AND DEPARTMENT = :6`

	if _, err := tx.ExecContext(ctx, query, tbm.RegUser, tbm.RegUno, agent, tbm.Jno, tbm.TbmDate, tbm.Department); err != nil {
		return fmt.Errorf("ModifyTbmExcel: %w", err)
	}
	return nil
}
