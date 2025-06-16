package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/guregu/null"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-17
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// func: 전체 근로자 조회
// @param
// - page entity.PageSql: 정렬, 리스트 수
// - search entity.WorkerSql: 검색 단어
// - retry string: 통합검색 텍스트
func (r *Repository) GetWorkerTotalList(ctx context.Context, db Queryer, page entity.PageSql, search entity.Worker, retry string) (*entity.Workers, error) {
	workers := entity.Workers{}

	condition := ""
	condition = utils.StringWhereConvert(condition, search.JobName.NullString, "t2.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.UserId.NullString, "t1.USER_ID")
	condition = utils.StringWhereConvert(condition, search.UserNm.NullString, "t1.USER_NM")
	condition = utils.StringWhereConvert(condition, search.Department.NullString, "t1.DEPARTMENT")
	condition = utils.StringWhereConvert(condition, search.Phone.NullString, "t1.PHONE")
	condition = utils.StringWhereConvert(condition, search.WorkerType.NullString, "t1.WORKER_TYPE")
	condition = utils.StringWhereConvert(condition, search.DiscName.NullString, "t1.DISC_NAME")
	var columns []string
	columns = append(columns, "t2.JOB_NAME")
	columns = append(columns, "t1.USER_NM")
	columns = append(columns, "t1.DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	var order string
	if page.Order.Valid {
		order = page.Order.String
	} else {
		order = `
				(
					CASE 
						WHEN REG_DATE IS NULL THEN MOD_DATE 
						WHEN MOD_DATE IS NULL THEN REG_DATE 
						ELSE GREATEST(REG_DATE, MOD_DATE) 
					END
				) DESC NULLS LAST`
	}

	query := fmt.Sprintf(`
				SELECT *
				FROM (
					SELECT 
					    ROWNUM AS RNUM,
						sorted_data.SNO,
						sorted_data.SITE_NM,
						sorted_data.JNO,
						sorted_data.JOB_NAME,
						sorted_data.USER_ID, 
						sorted_data.USER_NM,
						sorted_data.DEPARTMENT,
						sorted_data.DISC_NAME,
						sorted_data.PHONE, 
						sorted_data.WORKER_TYPE,
						sorted_data.IS_RETIRE,
						sorted_data.RETIRE_DATE,
						sorted_data.IS_MANAGE,
						sorted_data.REG_USER,
						sorted_data.REG_DATE,
						sorted_data.MOD_USER,
						sorted_data.MOD_DATE,
						sorted_data.REG_NO
					FROM (
						SELECT 
						    t1.WNO,
							t1.SNO,
							t3.SITE_NM,
							t1.JNO,
							t2.JOB_NAME,
							t1.USER_ID, 
							t1.USER_NM,
							t1.DEPARTMENT,
							t1.DISC_NAME,
							t1.PHONE, 
							t1.WORKER_TYPE,
							t1.IS_RETIRE,
							t1.RETIRE_DATE,
							t1.IS_MANAGE,
							t1.REG_USER,
							t1.REG_DATE,
							t1.MOD_USER,
							t1.MOD_DATE,
							t1.REG_NO
						FROM IRIS_WORKER_SET t1
						LEFT JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
						LEFT JOIN IRIS_SITE_SET t3 ON t1.SNO = t3.SNO
						WHERE t1.SNO > 100
						%s %s
						ORDER BY %s
					) sorted_data
					WHERE ROWNUM <= :1
					ORDER BY RNUM %s
				)
				WHERE RNUM > :2`, condition, retryCondition, order, page.RnumOrder)

	if err := db.SelectContext(ctx, &workers, query, page.EndNum, page.StartNum); err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("GetWorkerTotalList err: %v", err)
	}

	return &workers, nil
}

// func: 전체 근로자 개수 조회
// @param
// - searchTime string: 조회 날짜
// - retry string: 통합검색 텍스트
func (r *Repository) GetWorkerTotalCount(ctx context.Context, db Queryer, search entity.Worker, retry string) (int, error) {
	var count int

	condition := ""
	condition = utils.StringWhereConvert(condition, search.JobName.NullString, "t2.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.UserId.NullString, "t1.USER_ID")
	condition = utils.StringWhereConvert(condition, search.UserNm.NullString, "t1.USER_NM")
	condition = utils.StringWhereConvert(condition, search.Department.NullString, "t1.DEPARTMENT")
	condition = utils.StringWhereConvert(condition, search.Phone.NullString, "t1.PHONE")
	condition = utils.StringWhereConvert(condition, search.WorkerType.NullString, "t1.WORKER_TYPE")

	var columns []string
	columns = append(columns, "t2.JOB_NAME")
	columns = append(columns, "t1.USER_NM")
	columns = append(columns, "t1.DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	query := fmt.Sprintf(`
						SELECT 
							COUNT(*)
						FROM IRIS_WORKER_SET t1
						LEFT JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
						WHERE t1.SNO > 100
						%s %s`, condition, retryCondition)

	if err := db.GetContext(ctx, &count, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//TODO: 에러 아카이브
			return 0, nil
		}
		//TODO: 에러 아카이브
		return 0, fmt.Errorf("GetWorkerTotalCount fail: %w", err)
	}
	return count, nil
}

// func: 근로자 검색(현장근로자 추가시 사용)
// @param
// - userId string
func (r *Repository) GetAbsentWorkerList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerDaily, retry string) (*entity.Workers, error) {
	workers := entity.Workers{}

	var columns []string
	columns = append(columns, "USER_ID")
	columns = append(columns, "USER_NM")
	columns = append(columns, "DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	query := fmt.Sprintf(`
				SELECT *
				FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						SELECT USER_ID, USER_NM, DEPARTMENT, :1 as RECORD_DATE
						FROM IRIS_WORKER_SET
						WHERE JNO = :2
						AND USER_ID NOT IN (
							SELECT USER_ID
							FROM IRIS_WORKER_DAILY_SET
							WHERE JNO = :3
							AND TO_CHAR(RECORD_DATE, 'YYYY-MM-DD') = :4
						)
						%s
					) sorted_data
					WHERE ROWNUM <= :5
				)
				WHERE RNUM > :6`, retryCondition)

	if err := db.SelectContext(ctx, &workers, query, search.SearchStartTime, search.Jno, search.Jno, search.SearchStartTime, page.EndNum, page.StartNum); err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("GetAbsentWorkerList fail: %v", err)
	}

	return &workers, nil
}

// func: 근로자 개수 검색(현장근로자 추가시 사용)
// @param
// - userId string
func (r *Repository) GetAbsentWorkerCount(ctx context.Context, db Queryer, search entity.WorkerDaily, retry string) (int, error) {
	var count int

	var columns []string
	columns = append(columns, "USER_ID")
	columns = append(columns, "USER_NM")
	columns = append(columns, "DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	query := fmt.Sprintf(`
				SELECT COUNT(*)
				FROM IRIS_WORKER_SET
				WHERE JNO = :1
				AND USER_ID NOT IN (
					SELECT USER_ID
					FROM IRIS_WORKER_DAILY_SET
					WHERE JNO = :2
					AND TO_CHAR(RECORD_DATE, 'YYYY-MM-DD') = :3
				)
				%s`, retryCondition)

	if err := db.GetContext(ctx, &count, query, search.Jno, search.Jno, search.SearchStartTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//TODO: 에러 아카이브
			return 0, nil
		}
		//TODO: 에러 아카이브
		return 0, fmt.Errorf("GetAbsentWorkerCount fail: %w", err)
	}
	return count, nil
}

// 프로젝트에 참여한 회사명 리스트
func (r *Repository) GetWorkerDepartList(ctx context.Context, db Queryer, jno int64) ([]string, error) {
	var list []string

	query := `
		SELECT DISTINCT
		  CASE
			WHEN INSTR(DEPARTMENT, ' ', -1) > 0 THEN SUBSTR(DEPARTMENT, 1, INSTR(DEPARTMENT, ' ', -1) - 1)
			ELSE DEPARTMENT
		  END AS COMPANY_NAME
		FROM IRIS_WORKER_SET
		WHERE JNO = :1
		  AND DEPARTMENT IS NOT NULL`

	if err := db.SelectContext(ctx, &list, query, jno); err != nil {
		return nil, fmt.Errorf("GetWorkerDepartList fail: %w", err)
	}
	return list, nil
}

// func: 근로자 추가
// @param
// -
func (r *Repository) AddWorker(ctx context.Context, tx Execer, worker entity.Worker) error {
	agent := utils.GetAgent()

	// IRIS_WORKER_SET에 INSERT하는 쿼리
	insertQuery := `
		INSERT INTO IRIS_WORKER_SET(
			SNO, JNO, USER_ID, USER_NM, DEPARTMENT, 
			DISC_NAME, PHONE, WORKER_TYPE, IS_RETIRE, REG_DATE, 
			REG_AGENT, REG_USER, REG_UNO, REG_NO
		) VALUES (
			:1, :2, :3, :4, :5, 
			:6, :7, :8, :9, SYSDATE, 
			:10, :11, :12, :13
		)`

	_, err := tx.ExecContext(ctx, insertQuery,
		worker.Sno, worker.Jno, worker.UserId, worker.UserNm, worker.Department,
		worker.DiscName, worker.Phone, worker.WorkerType, worker.IsRetire, /*, SYSDATE*/
		agent, worker.RegUser, worker.RegUno, worker.RegNo,
	)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("AddWorker; IRIS_WORKER_SET INSERT fail: %v", err)
	}

	return nil
}

// func: 근로자 수정
// @param
// -
func (r *Repository) ModifyWorker(ctx context.Context, tx Execer, worker entity.Worker) error {
	agent := utils.GetAgent()

	query := `
				UPDATE IRIS_WORKER_SET 
				SET 
					USER_NM = :1, DEPARTMENT = :2, PHONE = :3, WORKER_TYPE = :4, IS_RETIRE = :5,
					RETIRE_DATE = :6, MOD_DATE = SYSDATE, MOD_AGENT = :7, MOD_USER = :8, MOD_UNO = :9, TRG_EDITABLE_YN = 'N',
					REG_NO = :10, IS_MANAGE = :11, DISC_NAME=:12
				WHERE SNO = :13 AND JNO = :14 AND USER_ID = :15`

	result, err := tx.ExecContext(ctx, query,
		worker.UserNm, worker.Department, worker.Phone, worker.WorkerType, worker.IsRetire,
		worker.RetireDate /*, SYSDATE*/, agent, worker.ModUser, worker.ModUno,
		worker.RegNo, worker.IsManage, worker.DiscName,
		worker.Sno, worker.Jno, worker.UserId,
	)

	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("ModifyWorker fail: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("ModifyWorker RowsAffected fail: %v", err)
	}
	if rowsAffected == 0 {
		//TODO: 에러 아카이브
		return fmt.Errorf("Rows add/update cnt: %d\n", rowsAffected)
	}

	return nil
}

// func: 현장 근로자 조회
// @param
// - page entity.PageSql: 정렬, 리스트 수
// - search entity.WorkerSql: 검색 단어
func (r *Repository) GetWorkerSiteBaseList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerDaily, retry string) (*entity.WorkerDailys, error) {
	list := entity.WorkerDailys{}

	condition := ""

	condition = utils.StringWhereConvert(condition, search.UserId.NullString, "t1.USER_ID")
	condition = utils.StringWhereConvert(condition, search.UserNm.NullString, "t2.USER_NM")
	condition = utils.StringWhereConvert(condition, search.Department.NullString, "t2.DEPARTMENT")

	var columns []string
	columns = append(columns, "t1.USER_ID")
	columns = append(columns, "t2.USER_NM")
	columns = append(columns, "t2.DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	var order string
	if page.Order.Valid {
		order = page.Order.String
	} else {
		//order = "RECORD_DATE DESC, OUT_RECOG_TIME DESC NULLS LAST"
		order = `
				RECORD_DATE DESC, (
					CASE 
						WHEN REG_DATE IS NULL THEN MOD_DATE 
						WHEN MOD_DATE IS NULL THEN REG_DATE 
						ELSE GREATEST(REG_DATE, MOD_DATE) 
					END
				) DESC NULLS LAST`
	}

	query := fmt.Sprintf(`
				SELECT *
				FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						   	SELECT 
								t1.SNO AS SNO,
								t1.JNO AS JNO,
								t1.USER_ID AS USER_ID,
								t2.USER_NM AS USER_NM,
								t2.DEPARTMENT AS DEPARTMENT,
								t1.RECORD_DATE AS RECORD_DATE,
								t1.IN_RECOG_TIME AS IN_RECOG_TIME,
								t1.OUT_RECOG_TIME AS OUT_RECOG_TIME,
								t1.IS_DEADLINE AS IS_DEADLINE,
								t1.IS_OVERTIME AS IS_OVERTIME,
								t1.REG_USER AS REG_USER,
								t1.REG_DATE AS REG_DATE,
								t1.MOD_USER AS MOD_USER,
								t1.MOD_DATE AS MOD_DATE,
								t1.WORK_STATE AS WORK_STATE,
								t1.COMPARE_STATE AS COMPARE_STATE
							FROM IRIS_WORKER_DAILY_SET t1
							LEFT JOIN IRIS_WORKER_SET t2 ON t1.USER_ID = t2.USER_ID AND t1.sno = t2.sno
							WHERE t1.SNO > 100
							AND t1.COMPARE_STATE in ('S', 'X')
							AND t1.JNO = :1
							AND TO_CHAR(t1.RECORD_DATE, 'yyyy-mm-dd') BETWEEN :2 AND :3
							%s %s
							ORDER BY %s
					) sorted_data
					WHERE ROWNUM <= :4
					ORDER BY RNUM %s
				)
				WHERE RNUM > :5`, condition, retryCondition, order, page.RnumOrder)

	if err := db.SelectContext(ctx, &list, query, search.Jno, search.SearchStartTime, search.SearchEndTime, page.EndNum, page.StartNum); err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("GetWorkerSiteBaseList err: %v", err)
	}

	return &list, nil
}

// func: 현장 근로자 개수 조회
// @param
// - searchTime string: 조회 날짜
func (r *Repository) GetWorkerSiteBaseCount(ctx context.Context, db Queryer, search entity.WorkerDaily, retry string) (int, error) {
	var count int

	condition := ""

	condition = utils.StringWhereConvert(condition, search.UserId.NullString, "t1.USER_ID")
	condition = utils.StringWhereConvert(condition, search.UserNm.NullString, "t2.USER_NM")
	condition = utils.StringWhereConvert(condition, search.Department.NullString, "t2.DEPARTMENT")

	var columns []string
	columns = append(columns, "t1.USER_ID")
	columns = append(columns, "t2.USER_NM")
	columns = append(columns, "t2.DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	query := fmt.Sprintf(`
							SELECT 
								count(*)
							FROM IRIS_WORKER_DAILY_SET t1
							LEFT JOIN IRIS_WORKER_SET t2 ON t1.SNO = t2.SNO AND t1.USER_ID = t2.USER_ID 
							WHERE t1.SNO > 100
							AND t1.COMPARE_STATE in ('S', 'X')
							AND t1.JNO = :1
							AND TO_CHAR(t1.RECORD_DATE, 'yyyy-mm-dd') BETWEEN :2 AND :3
							%s %s`, condition, retryCondition)

	if err := db.GetContext(ctx, &count, query, search.Jno, search.SearchStartTime, search.SearchEndTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//TODO: 에러 아카이브
			return 0, nil
		}
		//TODO: 에러 아카이브
		return 0, fmt.Errorf("GetWorkerSiteBaseCount fail: %w", err)
	}
	return count, nil
}

// func: 현장 근로자 추가/수정
// @param
// -
func (r *Repository) MergeSiteBaseWorker(ctx context.Context, tx Execer, workers entity.WorkerDailys) error {
	agent := utils.GetAgent()

	query := `
				MERGE INTO IRIS_WORKER_DAILY_SET t1
				USING (
					SELECT 
						:1 AS SNO,
						:2 AS JNO,
						:3 AS USER_ID,
						:4 AS RECORD_DATE,
						:5 AS IN_RECOG_TIME,
						:6 AS OUT_RECOG_TIME,
						:7 AS REG_AGENT,
						:8 AS REG_USER,
						:9 AS REG_UNO,
						:10 AS IS_DEADLINE,
						:11 AS WORK_STATE,
						:12 AS IS_OVERTIME
					FROM DUAL
				) t2
				ON (
					t1.SNO = t2.SNO 
					AND t1.JNO = t2.JNO 
					AND t1.USER_ID = t2.USER_ID
				    AND t1.RECORD_DATE   = t2.RECORD_DATE
				) WHEN MATCHED THEN
					UPDATE SET
						t1.IN_RECOG_TIME = t2.IN_RECOG_TIME,
						t1.OUT_RECOG_TIME = t2.OUT_RECOG_TIME,
						t1.MOD_DATE      = SYSDATE,
						t1.MOD_AGENT     = t2.REG_AGENT,
						t1.MOD_USER      = t2.REG_USER,
						t1.MOD_UNO       = t2.REG_UNO,
				    	t1.IS_DEADLINE   = t2.IS_DEADLINE,
				    	t1.WORK_STATE = t2.WORK_STATE,
						t1.IS_OVERTIME   = t2.IS_OVERTIME
					WHERE t1.SNO = t2.SNO
					AND t1.JNO = t2.JNO
					AND t1.USER_ID = t2.USER_ID
				    AND t1.RECORD_DATE   = t2.RECORD_DATE
				WHEN NOT MATCHED THEN
					INSERT (SNO, JNO, USER_ID, RECORD_DATE, IN_RECOG_TIME, OUT_RECOG_TIME, WORK_STATE, COMPARE_STATE, REG_DATE, REG_AGENT, REG_USER, REG_UNO, IS_DEADLINE, IS_OVERTIME)
					VALUES (t2.SNO, t2.JNO, t2.USER_ID, t2.RECORD_DATE, t2.IN_RECOG_TIME, t2.OUT_RECOG_TIME, t2.WORK_STATE, 'X', SYSDATE, t2.REG_AGENT, t2.REG_USER, t2.REG_UNO, t2.IS_DEADLINE, t2.IS_OVERTIME)`

	for _, worker := range workers {
		_, err := tx.ExecContext(ctx, query,
			worker.Sno, worker.Jno, worker.UserId, worker.RecordDate, worker.InRecogTime,
			worker.OutRecogTime, agent, worker.ModUser, worker.ModUno, worker.IsDeadline,
			worker.WorkState, worker.IsOvertime,
		)
		if err != nil {
			//TODO: 에러 아카이브
			return fmt.Errorf("MergeSiteBaseWorker fail: %w", err)
		}
	}

	return nil
}

// 현장 근로자 변경사항 로그 저장
func (r *Repository) MergeSiteBaseWorkerLog(ctx context.Context, tx Execer, workers entity.WorkerDailys) error {
	agent := utils.GetAgent()

	query := `
		INSERT INTO IRIS_WORKER_DAILY_LOG(SNO, JNO, USER_ID, RECOG_TIME, TRANS_TYPE, MESSAGE, REG_DATE, REG_USER, REG_UNO, REG_AGENT)
		VALUES(:1, :2, :3, :4, :5, :6, SYSDATE, :7, :8, :9)`

	for _, worker := range workers {
		if _, err := tx.ExecContext(ctx, query, worker.Sno, worker.Jno, worker.UserId, worker.RecordDate, worker.WorkState, worker.Message, worker.ModUser, worker.ModUno, agent); err != nil {
			return fmt.Errorf("MergeSiteBaseWorkerLog fail: %w", err)
		}
	}
	return nil
}

// func: 현장 근로자 일괄마감
// @param
// -
func (r *Repository) ModifyWorkerDeadline(ctx context.Context, tx Execer, workers entity.WorkerDailys) error {
	agent := utils.GetAgent()

	query := `
				UPDATE IRIS_WORKER_DAILY_SET 
				SET 
					IS_DEADLINE = 'Y',
					MOD_DATE = SYSDATE,
					MOD_AGENT = :1,
					MOD_USER = :2,
					MOD_UNO = :3
				WHERE SNO = :4
				AND JNO = :5
				AND USER_ID = :6
				AND RECORD_DATE = :7`

	for _, worker := range workers {
		_, err := tx.ExecContext(ctx, query,
			agent, worker.ModUser, worker.ModUno, worker.Sno, worker.Jno,
			worker.UserId, worker.RecordDate,
		)
		if err != nil {
			//TODO: 에러 아카이브
			return fmt.Errorf("ModifyWorkerDeadline fail: %w", err)
		}
	}

	return nil
}

// func: 현장 근로자 프로젝트 변경
// @param
// -
func (r *Repository) ModifyWorkerProject(ctx context.Context, tx Execer, workers entity.WorkerDailys) error {
	agent := utils.GetAgent()

	query := `
				UPDATE IRIS_WORKER_DAILY_SET 
				SET 
				    JNO = :1,
					MOD_DATE = SYSDATE,
					MOD_AGENT = :2,
					MOD_USER = :3,
					MOD_UNO = :4
				WHERE SNO = :5
				AND JNO = :6
				AND USER_ID = :7
				AND RECORD_DATE = :8`

	for _, worker := range workers {
		_, err := tx.ExecContext(ctx, query,
			worker.AfterJno, agent, worker.ModUser, worker.ModUno, worker.Sno,
			worker.Jno, worker.UserId, worker.RecordDate,
		)
		if err != nil {
			//TODO: 에러 아카이브
			return fmt.Errorf("ModifyWorkerProject fail: %w", err)
		}
	}

	return nil
}

// 현장 근로자 프로젝트 변경시 같은 현장내 프로젝트일 경우 전체 근로자 프로젝트 변경
func (r *Repository) ModifyWorkerDefaultProject(ctx context.Context, tx Execer, workers entity.WorkerDailys) error {
	agent := utils.GetAgent()

	query := `
			UPDATE IRIS_WORKER_SET
			SET 
				JNO = :1,
				MOD_DATE = SYSDATE,
				MOD_USER = :2,
				MOD_UNO = :3,
				MOD_AGENT = :4
			WHERE SNO = :5
			AND USER_ID = :6
			AND EXISTS (
				SELECT 1
				FROM IRIS_SITE_JOB
				WHERE SNO = :7 AND JNO = :8
			)`

	for _, worker := range workers {
		if _, err := tx.ExecContext(ctx, query, worker.AfterJno, worker.ModUser, worker.ModUno, agent, worker.Sno, worker.UserId, worker.Sno, worker.Jno); err != nil {
			return fmt.Errorf("ModifyWorkerDefaultProject fail: %w", err)
		}
	}
	return nil
}

// func: 현장 근로자 일일 마감처리
// @param
// -
func (r *Repository) ModifyWorkerDeadlineInit(ctx context.Context, tx Execer) error {
	agent := utils.GetAgent()

	query := `
			UPDATE IRIS_WORKER_DAILY_SET 
			SET 
				IS_DEADLINE = 'Y',
				MOD_DATE = SYSDATE,
				MOD_AGENT = :1,
				MOD_USER = 'Scheduled'
			WHERE TRUNC(RECORD_DATE) >= TRUNC(SYSDATE) - 7
			AND TRUNC(RECORD_DATE) < TRUNC(SYSDATE)
			AND WORK_STATE = '02'
			AND IS_DEADLINE = 'N'
			AND COMPARE_STATE = 'S'`

	if _, err := tx.ExecContext(ctx, query, agent); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("ModifyWorkerDeadlineInit fail: %w", err)
	}

	return nil
}

// func: 철야 근로자 조회
// @param
// -
func (r *Repository) GetWorkerOverTime(ctx context.Context, db Queryer) (*entity.WorkerOverTimes, error) {

	workerOverTimes := entity.WorkerOverTimes{}
	query := `
			SELECT 
				w1.CNO AS BEFORE_CNO, 
				w2.OUT_RECOG_TIME AS OUT_RECOG_TIME, 
				w2.CNO AS AFTER_CNO 
			FROM iris_worker_daily_set w1 
			INNER JOIN iris_worker_daily_set w2 
			ON w1.user_id = w2.user_id AND w1.jno = w2.jno 
			WHERE to_date(w2.record_date) = TRUNC(SYSDATE) 
			  AND w2.IN_RECOG_TIME IS NULL 
			  AND w2.OUT_RECOG_TIME IS NOT NULL 
			  AND TO_DATE(w1.RECORD_DATE) = TRUNC(SYSDATE - 1) 
			  AND w1.IN_RECOG_TIME IS NOT NULL 
			  AND w1.OUT_RECOG_TIME IS NULL
			  AND W2.COMPARE_STATE = 'S'
		`

	if err := db.SelectContext(ctx, &workerOverTimes, query); err != nil {
		return nil, fmt.Errorf("GetWorkerOverTime fail: %w", err)
	}

	return &workerOverTimes, nil

}

// func: 현장 근로자 철야 처리
// @param
// - workerOverTime entity.WorkerOverTime: BeforeCno, AfterCno, OutRecogTime
func (r *Repository) ModifyWorkerOverTime(ctx context.Context, tx Execer, workerOverTime entity.WorkerOverTime) error {
	agent := utils.GetAgent()

	query := `
		UPDATE 
		    IRIS_WORKER_DAILY_SET 
		SET 
		    OUT_RECOG_TIME = :1,
		    IS_OVERTIME = 'Y',
		    WORK_STATE = '02',
			MOD_DATE = SYSDATE,
			MOD_AGENT = :2,
			MOD_USER = 'Scheduled'
		WHERE 
		    CNO = :3
			
	`

	if _, err := tx.ExecContext(ctx, query, workerOverTime.OutRecogTime, agent, workerOverTime.BeforeCno); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("ModifyWorkerOverTime fail: %w", err)
	}
	return nil

}

// func: 현장 근로자 철야 처리 후 삭제
// @param
// - cno: 근로자 PK
func (r *Repository) DeleteWorkerOverTime(ctx context.Context, tx Execer, cno null.Int) error {

	query := `
		DELETE FROM iris_worker_daily_set
		WHERE  CNO = :1
		`
	if _, err := tx.ExecContext(ctx, query, cno); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("DeleteWorkerOverTime fail: %w", err)
	}
	return nil
}
