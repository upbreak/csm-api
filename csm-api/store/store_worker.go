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
						sorted_data.JNO,
						sorted_data.JOB_NAME,
						sorted_data.USER_ID, 
						sorted_data.USER_NM,
						sorted_data.DEPARTMENT,
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
							t1.JNO,
							t2.JOB_NAME,
							t1.USER_ID, 
							t1.USER_NM,
							t1.DEPARTMENT,
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

// func: 근로자 추가
// @param
// -
func (r *Repository) AddWorker(ctx context.Context, tx Execer, worker entity.Worker) error {
	agent := utils.GetAgent()

	// IRIS_WORKER_SET에 INSERT하는 쿼리
	insertQuery := `
		INSERT INTO IRIS_WORKER_SET(
			SNO, JNO, USER_ID, USER_NM, DEPARTMENT, 
			PHONE, WORKER_TYPE, IS_RETIRE, REG_DATE, REG_AGENT, 
			REG_USER, REG_UNO, REG_NO
		) VALUES (
			:1, :2, :3, :4, :5, 
			:6, :7, :8, SYSDATE, :9, 
			:10, :11, :12
		)`

	_, err := tx.ExecContext(ctx, insertQuery,
		worker.Sno, worker.Jno, worker.UserId, worker.UserNm, worker.Department,
		worker.Phone, worker.WorkerType, worker.IsRetire /*, SYSDATE*/, agent,
		worker.RegUser, worker.RegUno, worker.RegNo,
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
					REG_NO = :10, IS_MANAGE = :11
				WHERE SNO = :12 AND JNO = :13 AND USER_ID = :14`

	result, err := tx.ExecContext(ctx, query,
		worker.UserNm, worker.Department, worker.Phone, worker.WorkerType, worker.IsRetire,
		worker.RetireDate /*, SYSDATE*/, agent, worker.ModUser, worker.ModUno,
		worker.RegNo, worker.IsManage,
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
								t1.REG_USER AS REG_USER,
								t1.REG_DATE AS REG_DATE,
								t1.MOD_USER AS MOD_USER,
								t1.MOD_DATE AS MOD_DATE,
								t1.WORK_STATE AS WORK_STATE
							FROM IRIS_WORKER_DAILY_SET t1
							LEFT JOIN IRIS_WORKER_SET t2 ON t1.USER_ID = t2.USER_ID AND t1.sno = t2.sno AND t1.jno = t2.jno
							WHERE t1.SNO > 100
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
							LEFT JOIN IRIS_WORKER_SET t2 ON t1.SNO = t2.SNO AND t1.JNO  = t2.JNO AND t1.USER_ID = t2.USER_ID 
							WHERE t1.SNO > 100
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
						:11 AS WORK_STATE
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
				    	t1.WORK_STATE = t2.WORK_STATE
					WHERE t1.SNO = t2.SNO
					AND t1.JNO = t2.JNO
					AND t1.USER_ID = t2.USER_ID
				    AND t1.RECORD_DATE   = t2.RECORD_DATE
				WHEN NOT MATCHED THEN
					INSERT (SNO, JNO, USER_ID, RECORD_DATE, IN_RECOG_TIME, OUT_RECOG_TIME, WORK_STATE, REG_DATE, REG_AGENT, REG_USER, REG_UNO, IS_DEADLINE)
					VALUES (t2.SNO, t2.JNO, t2.USER_ID, t2.RECORD_DATE, t2.IN_RECOG_TIME, t2.OUT_RECOG_TIME, t2.WORK_STATE, SYSDATE, t2.REG_AGENT, t2.REG_USER, t2.REG_UNO, t2.IS_DEADLINE)`

	for _, worker := range workers {
		_, err := tx.ExecContext(ctx, query,
			worker.Sno, worker.Jno, worker.UserId, worker.RecordDate, worker.InRecogTime,
			worker.OutRecogTime, agent, worker.ModUser, worker.ModUno, worker.IsDeadline,
			worker.WorkState,
		)
		if err != nil {
			//TODO: 에러 아카이브
			return fmt.Errorf("MergeSiteBaseWorker fail: %w", err)
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
