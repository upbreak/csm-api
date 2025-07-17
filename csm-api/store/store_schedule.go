package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
)

// func: 휴무일 조회
// @param
// -
func (r *Repository) GetRestScheduleList(ctx context.Context, db Queryer, jno int64, year string, month string) (entity.RestSchedules, error) {
	list := entity.RestSchedules{}

	query := `
			SELECT 
			    CNO,
				JNO,
				IS_EVERY_YEAR,
				REST_YEAR,
				REST_MONTH,
				REST_DAY,
				REASON
			FROM IRIS_SCH_REST_SET
			WHERE 
			  (
				(IS_EVERY_YEAR = 'Y' AND TO_CHAR(REST_MONTH) = :1 OR :2 IS NULL)
				OR
				(IS_EVERY_YEAR = 'N' AND TO_CHAR(REST_YEAR) = :3 AND TO_CHAR(REST_MONTH) = :4 OR :5 IS NULL)
			  )
			AND (
			  :6 = 0 OR (JNO = :7 OR JNO = 0)
			)`

	if err := db.SelectContext(ctx, &list, query, month, month, year, month, month, jno, jno); err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return list, nil
}

// func: 휴무일 추가
// @param
// -
func (r *Repository) AddRestSchedule(ctx context.Context, tx Execer, schedule entity.RestSchedules) error {
	agent := utils.GetAgent()

	query := `
			INSERT INTO IRIS_SCH_REST_SET(
				JNO, IS_EVERY_YEAR, REST_YEAR, REST_MONTH, REST_DAY, 
			    REASON, REG_DATE, REG_AGENT, REG_UNO, REG_USER
			) VALUES (
				:1, :2, :3, :4, :5, 
				:6, SYSDATE, :8, :9, :10
			)`

	for _, rest := range schedule {
		if _, err := tx.ExecContext(ctx, query,
			rest.Jno, rest.IsEveryYear, rest.RestYear, rest.RestMonth, rest.RestDay,
			rest.Reason /*SYSDATE*/, agent, rest.RegUno, rest.RegUser,
		); err != nil {
			return utils.CustomErrorf(err)
		}
	}

	return nil
}

// func: 휴무일 수정
// @param
// -
func (r *Repository) ModifyRestSchedule(ctx context.Context, tx Execer, schedule entity.RestSchedule) error {
	agent := utils.GetAgent()

	query := `
			UPDATE IRIS_SCH_REST_SET
			SET	
				JNO = :1,
				IS_EVERY_YEAR = :2,
				REST_YEAR = :3,
				REST_MONTH = :4,
				REST_DAY = :5,
				REASON = :6,
				MOD_DATE = SYSDATE,
				MOD_AGENT = :7,
				MOD_UNO = :8,
				MOD_USER = :9
			WHERE CNO = :10`

	if _, err := tx.ExecContext(ctx, query, schedule.Jno, schedule.IsEveryYear, schedule.RestYear, schedule.RestMonth, schedule.RestDay, schedule.Reason, agent, schedule.ModUno, schedule.ModUser, schedule.Cno); err != nil {
		return utils.CustomErrorf(err)
	}

	return nil
}

// func: 휴무일 삭제
// @param
// -
func (r *Repository) RemoveRestSchedule(ctx context.Context, tx Execer, cno int64) error {
	query := `DELETE FROM IRIS_SCH_REST_SET WHERE CNO = :1`

	if _, err := tx.ExecContext(ctx, query, cno); err != nil {
		return utils.CustomErrorf(err)
	}
	return nil
}
