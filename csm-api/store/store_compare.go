package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"fmt"
	"github.com/guregu/null"
)

// 일일 근로자 비교 - 근로자 리스트
func (r *Repository) GetDailyWorkerList(ctx context.Context, db Queryer, jno int64, startDate null.Time, retry string, order string) (entity.WorkerDailys, error) {
	var list entity.WorkerDailys

	var columns []string
	columns = append(columns, "T2.USER_NM")
	columns = append(columns, "T2.DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	var orderBy string
	if order == "" {
		orderBy = `
			ORDER BY
				CASE COMPARE_STATE
					WHEN 'S' THEN 1
					WHEN 'W' THEN 2
					WHEN 'C' THEN 3
					ELSE 4
				END,
				CASE 
					WHEN IN_RECOG_TIME IS NULL THEN OUT_RECOG_TIME 
					WHEN OUT_RECOG_TIME IS NULL THEN IN_RECOG_TIME 
					ELSE GREATEST(IN_RECOG_TIME, OUT_RECOG_TIME) 
				END
				DESC NULLS LAST`
	} else {
		orderBy = fmt.Sprintf(`ORDER BY %s`, order)
	}

	query := fmt.Sprintf(`
		SELECT
			T1.SNO,
			T1.JNO,
			T1.USER_ID,
			T2.USER_NM,
			T2.REG_NO,
			CASE
				WHEN INSTR(T2.DEPARTMENT, ' ', -1) > 0 THEN SUBSTR(T2.DEPARTMENT, 1, INSTR(T2.DEPARTMENT, ' ', -1) - 1)
				ELSE T2.DEPARTMENT
			END AS DEPARTMENT,
		    T2.DISC_NAME,
			T1.IN_RECOG_TIME,
			T1.OUT_RECOG_TIME,
			TRUNC(T1.RECORD_DATE) AS RECORD_DATE,
			T1.COMPARE_STATE,
			T1.IS_DEADLINE
		FROM IRIS_WORKER_DAILY_SET T1
		LEFT JOIN IRIS_WORKER_SET T2 ON T1.SNO = T2.SNO AND T1.JNO = T2.JNO AND T1.USER_ID = T2.USER_ID
		WHERE T1.JNO = :1
		AND TRUNC(T1.RECORD_DATE) = TRUNC(:2)
		%s
		%s`, retryCondition, orderBy)

	if err := db.SelectContext(ctx, &list, query, jno, startDate); err != nil {
		return list, fmt.Errorf("GetDailyWorkerList: %w", err)
	}
	return list, nil
}

// 일일 근로자 비교 - TBM 리스트
func (r *Repository) GetTbmList(ctx context.Context, db Queryer, jno int64, startDate null.Time, retry string, order string) ([]entity.Tbm, error) {
	var list []entity.Tbm

	var columns []string
	columns = append(columns, "USER_NM")
	columns = append(columns, "DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	var orderBy string
	if order == "" {
		orderBy = `ORDER BY TBM_DATE DESC NULLS LAST`
	} else {
		orderBy = fmt.Sprintf(`ORDER BY %s`, order)
	}

	query := fmt.Sprintf(`
		SELECT 
			A.SNO,
			A.JNO,
			A.DEPARTMENT,
			A.DISC_NAME,
			A.USER_NM,
			TRUNC(A.TBM_DATE) AS TBM_DATE,
			A.TBM_ORDER
		FROM IRIS_TBM_SET A
		JOIN (
			SELECT 
				SNO, JNO, DEPARTMENT, USER_NM, MAX(TBM_ORDER) AS MAX_ORDER
			FROM IRIS_TBM_SET
			WHERE JNO = :1
			  AND TRUNC(TBM_DATE) = TRUNC(:2)
			GROUP BY SNO, JNO, DEPARTMENT, USER_NM
		) B
		  ON A.SNO = B.SNO
		 AND A.JNO = B.JNO
		 AND A.DEPARTMENT = B.DEPARTMENT
		 AND A.USER_NM = B.USER_NM
		 AND A.TBM_ORDER = B.MAX_ORDER
		WHERE A.JNO = :3
		 AND TRUNC(A.TBM_DATE) = TRUNC(:4)
		%s
		%s`, retryCondition, orderBy)

	if err := db.SelectContext(ctx, &list, query, jno, startDate, jno, startDate); err != nil {
		return list, fmt.Errorf("GetTbmList: %w", err)
	}
	return list, nil
}

// 일일 근로자 비교 - 퇴직공제 리스트
func (r *Repository) GetDeductionList(ctx context.Context, db Queryer, jno int64, startDate null.Time, retry string, order string) ([]entity.Deduction, error) {
	var list []entity.Deduction

	var columns []string
	columns = append(columns, "USER_NM")
	columns = append(columns, "DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	var orderBy string
	if order == "" {
		orderBy = `
			ORDER BY (
				CASE 
					WHEN IN_RECOG_TIME IS NULL THEN OUT_RECOG_TIME 
					WHEN OUT_RECOG_TIME IS NULL THEN IN_RECOG_TIME 
					ELSE GREATEST(IN_RECOG_TIME, OUT_RECOG_TIME) 
				END
			) DESC NULLS LAST`
	} else {
		orderBy = fmt.Sprintf(`ORDER BY %s`, order)
	}

	query := fmt.Sprintf(`
		SELECT 
			A.SNO,
			A.JNO,
			A.USER_NM,
			A.GENDER,
			A.REG_NO,
			A.DEPARTMENT,
			A.IN_RECOG_TIME,
			A.OUT_RECOG_TIME,
			TRUNC(A.RECORD_DATE) AS RECORD_DATE,
			A.DEDUCT_ORDER
		FROM IRIS_DEDUCTION_SET A
		JOIN (
			SELECT 
				SNO, JNO, USER_NM, REG_NO, DEPARTMENT, GENDER, MAX(DEDUCT_ORDER) AS MAX_ORDER
			FROM IRIS_DEDUCTION_SET
			WHERE JNO = :1
			  AND TRUNC(RECORD_DATE) = TRUNC(:2)
			GROUP BY SNO, JNO, USER_NM, REG_NO, DEPARTMENT, GENDER
		) B
		  ON A.SNO = B.SNO
		 AND A.JNO = B.JNO
		 AND A.USER_NM = B.USER_NM
		 AND A.REG_NO = B.REG_NO
		 AND A.DEPARTMENT = B.DEPARTMENT
		 AND A.GENDER = B.GENDER
		 AND A.DEDUCT_ORDER = B.MAX_ORDER
		WHERE A.JNO = :3
		  AND TRUNC(A.RECORD_DATE) = TRUNC(:4)
		%s
		%s`, retryCondition, orderBy)

	if err := db.SelectContext(ctx, &list, query, jno, startDate, jno, startDate); err != nil {
		return list, fmt.Errorf("GetDeductionList: %w", err)
	}
	return list, nil
}

// 근로자 비교 반영/취소
func (r *Repository) ModifyWorkerCompareState(ctx context.Context, tx Execer, workers entity.WorkerDailys) error {
	agent := utils.GetAgent()

	query := `
		UPDATE IRIS_WORKER_DAILY_SET
		SET
			COMPARE_STATE = :1,
			MOD_DATE = SYSDATE,
			MOD_USER = :2,
			MOD_UNO = :3,
			MOD_AGENT = :4
		WHERE JNO = :5
		AND USER_ID = :6
		AND TRUNC(RECORD_DATE) = TRUNC(:7)`

	for _, worker := range workers {
		if _, err := tx.ExecContext(ctx, query, worker.AfterState, worker.RegUser, worker.RegUno, agent, worker.Jno, worker.UserId, worker.RecordDate); err != nil {
			return fmt.Errorf("ModifyWorkerCompareState: %v", err)
		}
	}
	return nil
}

// 근로자 비교 반영 로그
func (r *Repository) AddCompareLog(ctx context.Context, tx Execer, logs entity.WorkerDailys) error {
	agent := utils.GetAgent()

	query := `
		INSERT INTO IRIS_COMPARE_LOG(JNO, USER_ID, USER_NM, BEFORE_STATE, AFTER_STATE, RECORD_DATE, REG_DATE, REG_USER, REG_UNO, REG_AGENT)
		VALUES(:1, :2, :3, :4, :5, :6, SYSDATE, :7, :8, :9)`

	for _, log := range logs {
		if _, err := tx.ExecContext(ctx, query, log.Jno, log.UserId, log.UserNm, log.BeforeState, log.AfterState, log.RecordDate, log.RegUser, log.RegUno, agent); err != nil {
			return fmt.Errorf("AddCompareLog: %w", err)
		}
	}
	return nil
}
