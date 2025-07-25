package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"database/sql"
	"fmt"
	"time"
)

// 현장관리 당일 작업 내용 조회
func (r *Repository) GetProjectDailyContentList(ctx context.Context, db Queryer, jno int64, targetDate time.Time) (*entity.ProjectDailys, error) {
	projectDailys := entity.ProjectDailys{}

	// jno 변환: 0이면 NULL 처리, 아니면 Valid 값으로 설정
	var jnoParam sql.NullInt64
	if jno != 0 {
		jnoParam = sql.NullInt64{Valid: true, Int64: jno}
	} else {
		jnoParam = sql.NullInt64{Valid: false}
	}

	// targetDate 변환: zero 값이면 NULL 처리, 아니면 Valid 값으로 설정
	var targetDateParam sql.NullTime
	if !targetDate.IsZero() {
		targetDateParam = sql.NullTime{Valid: true, Time: targetDate}
	} else {
		targetDateParam = sql.NullTime{Valid: false}
	}

	sql := `SELECT 
				t1.JNO,
				t1.IDX,
				t1.CONTENT,
				t1.IS_USE,
				t1.TARGET_DATE,
				t1.REG_DATE,
				t1.MOD_DATE,
				t1.REG_UNO,
				t1.REG_USER,
				t1.MOD_UNO,
				t1.MOD_USER
			FROM
				IRIS_DAILY_JOB t1
			WHERE
-- 				t1.IS_USE = 'Y'
-- 				AND 
			    t1.JNO = :2
				AND TO_CHAR(t1.TARGET_DATE, 'YYYY-MM-DD') = TO_CHAR(:2 , 'YYYY-MM-DD')
			ORDER BY
				NVL(t1.REG_DATE, t1.MOD_DATE) DESC`

	if err := db.SelectContext(ctx, &projectDailys, sql, jnoParam, targetDateParam); err != nil {
		return nil, utils.CustomErrorf(err)
	}
	return &projectDailys, nil
}

// 작업내용 조회
func (r *Repository) GetDailyJobList(ctx context.Context, db Queryer, jno int64, targetDate string) (entity.ProjectDailys, error) {
	projectDailys := entity.ProjectDailys{}

	query := `
			SELECT 
				IDX,
				JNO,
				CONTENT,
				TARGET_DATE
			FROM IRIS_DAILY_JOB
			WHERE TO_CHAR(TARGET_DATE, 'YYYY-MM') = :1
			AND (:2 = 0 OR (JNO = :3 OR JNO = 0))`

	if err := db.SelectContext(ctx, &projectDailys, query, targetDate, jno, jno); err != nil {
		return entity.ProjectDailys{}, utils.CustomErrorf(err)
	}
	return projectDailys, nil
}

// 작업내용 추가
func (r *Repository) AddDailyJob(ctx context.Context, tx Execer, project entity.ProjectDailys) error {
	query := `
		INSERT INTO IRIS_DAILY_JOB(JNO, CONTENT, TARGET_DATE, REG_DATE, REG_UNO, REG_USER)
			SELECT
				:1, :2, :3, SYSDATE, :4, :5
			FROM dual
			WHERE NOT EXISTS (
				SELECT 1
				FROM IRIS_DAILY_JOB
				WHERE 
					JNO = :6
					AND TRUNC(TARGET_DATE) = TRUNC(:7)
			)
		`

	for _, job := range project {
		if result, err := tx.ExecContext(ctx, query, job.Jno, job.Content, job.TargetDate, job.RegUno, job.RegUser, job.Jno, job.TargetDate); err != nil {
			return utils.CustomErrorf(err)
		} else {
			count, _ := result.RowsAffected()
			if count == 0 {
				return utils.CustomErrorf(fmt.Errorf("중복데이터 존재"))
			}
		}
	}

	return nil
}

// 작업내용 수정
func (r *Repository) ModifyDailyJob(ctx context.Context, tx Execer, project entity.ProjectDaily) error {
	query := `
			UPDATE IRIS_DAILY_JOB 
			SET 
				JNO = :1,
				CONTENT = :2,
				TARGET_DATE = :3,
				MOD_DATE = SYSDATE,
				MOD_UNO = :4,
				MOD_USER = :5
			WHERE 
			    IDX = :6
				AND NOT EXISTS (
					SELECT	1
					FROM IRIS_DAILY_JOB
					WHERE
						JNO = :7
						AND TRUNC(TARGET_DATE) = TRUNC(:8)
						AND IDX != :9			
					)
			`

	if result, err := tx.ExecContext(ctx, query, project.Jno, project.Content, project.TargetDate, project.RegUno, project.RegUser, project.Idx, project.Jno, project.TargetDate, project.Idx); err != nil {
		return utils.CustomErrorf(err)
	} else {
		count, _ := result.RowsAffected()
		if count == 0 {
			return utils.CustomErrorf(fmt.Errorf("중복데이터 존재"))
		}
	}

	return nil
}

// 작업내용 삭제
func (r *Repository) RemoveDailyJob(ctx context.Context, tx Execer, idx int64) error {
	query := `DELETE FROM IRIS_DAILY_JOB WHERE IDX = :1`

	if _, err := tx.ExecContext(ctx, query, idx); err != nil {
		return utils.CustomErrorf(err)
	}
	return nil
}
