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
func (r *Repository) GetWorkerTotalList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerSql, retry string) (*entity.WorkerSqls, error) {
	sqls := entity.WorkerSqls{}

	condition := ""
	condition = utils.StringWhereConvert(condition, search.JobName, "t2.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.UserId, "t1.USER_ID")
	condition = utils.StringWhereConvert(condition, search.UserNm, "t1.USER_NM")
	condition = utils.StringWhereConvert(condition, search.Department, "t1.DEPARTMENT")
	condition = utils.StringWhereConvert(condition, search.Phone, "t1.PHONE")
	condition = utils.StringWhereConvert(condition, search.WorkerType, "t1.WORKER_TYPE")

	var columns []string
	columns = append(columns, "t2.JOB_NAME")
	columns = append(columns, "t1.USER_NM")
	columns = append(columns, "t1.DEPARTMENT")
	retryCondition := utils.RetrySearchTextConvert(retry, columns)

	var order string
	if page.Order.Valid {
		order = page.Order.String
	} else {
		order = "WNO DESC"
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
						sorted_data.REG_USER,
						sorted_data.REG_DATE,
						sorted_data.MOD_USER,
						sorted_data.MOD_DATE
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
							t1.REG_USER,
							t1.REG_DATE,
							t1.MOD_USER,
							t1.MOD_DATE
						FROM IRIS_WORKER_SET t1
						LEFT JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
						WHERE t1.SNO > 100
						%s %s
					) sorted_data
					WHERE ROWNUM <= :1
					ORDER BY %s
				)
				WHERE RNUM > :2`, condition, retryCondition, order)

	if err := db.SelectContext(ctx, &sqls, query, page.EndNum, page.StartNum); err != nil {
		return nil, fmt.Errorf("GetWorkerTotalList err: %v", err)
	}

	return &sqls, nil
}

// func: 전체 근로자 개수 조회
// @param
// - searchTime string: 조회 날짜
// - retry string: 통합검색 텍스트
func (r *Repository) GetWorkerTotalCount(ctx context.Context, db Queryer, search entity.WorkerSql, retry string) (int, error) {
	var count int

	condition := ""
	condition = utils.StringWhereConvert(condition, search.JobName, "t2.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.UserId, "t1.USER_ID")
	condition = utils.StringWhereConvert(condition, search.UserNm, "t1.USER_NM")
	condition = utils.StringWhereConvert(condition, search.Department, "t1.DEPARTMENT")
	condition = utils.StringWhereConvert(condition, search.Phone, "t1.PHONE")
	condition = utils.StringWhereConvert(condition, search.WorkerType, "t1.WORKER_TYPE")

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
			return 0, nil
		}
		return 0, fmt.Errorf("GetWorkerTotalCount fail: %w", err)
	}
	return count, nil
}

// func: 근로자 추가
// @param
// -
func (r *Repository) AddWorker(ctx context.Context, db Beginner, worker entity.WorkerSql) error {
	// 트랜잭션 시작
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	agent := utils.GetAgent()

	// IRIS_WORKER_SET에 INSERT하는 쿼리
	insertQuery := `
		INSERT INTO IRIS_WORKER_SET(
			SNO, JNO, USER_ID, USER_NM, DEPARTMENT, 
			PHONE, WORKER_TYPE, IS_RETIRE, REG_DATE, REG_AGENT, 
			REG_USER, REG_UNO
		) VALUES (
			:1, :2, :3, :4, :5, 
			:6, :7, :8, SYSDATE, :9, 
			:10, :12
		)`

	_, err = tx.ExecContext(ctx, insertQuery,
		worker.Sno, worker.Jno, worker.UserId, worker.UserNm, worker.Department,
		worker.Phone, worker.WorkerType, worker.IsRetire /*, SYSDATE*/, agent,
		worker.RegUser, worker.RegUno,
	)

	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("AddWorker; IRIS_WORKER_SET INSERT fail: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

// func: 근로자 수정
// @param
// -
func (r *Repository) ModifyWorker(ctx context.Context, db Beginner, worker entity.WorkerSql) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Failed to begin transaction: %v", err)
	}

	agent := utils.GetAgent()

	query := `
				UPDATE IRIS_WORKER_SET 
				SET 
					USER_NM = :1, DEPARTMENT = :2, PHONE = :3, WORKER_TYPE = :4, IS_RETIRE = :5,
					RETIRE_DATE = :6, MOD_DATE = SYSDATE, MOD_AGENT = :7, MOD_USER = :8, MOD_UNO = :9, TRG_EDITABLE_YN = 'N' 
				WHERE SNO = :10 AND JNO = :11 AND USER_ID = :12`

	_, err = tx.ExecContext(ctx, query,
		worker.UserNm, worker.Department, worker.Phone, worker.WorkerType, worker.IsRetire,
		worker.RetireDate /*, SYSDATE*/, agent, worker.ModUser, worker.ModUno,
		worker.Sno, worker.Jno, worker.UserId,
	)

	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("ModifyWorker fail: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

// func: 현장 근로자 조회
// @param
// - page entity.PageSql: 정렬, 리스트 수
// - search entity.WorkerSql: 검색 단어
func (r *Repository) GetWorkerSiteBaseList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerSql) (*entity.WorkerSqls, error) {
	sqls := entity.WorkerSqls{}

	condition := ""

	condition = utils.StringWhereConvert(condition, search.JobName, "t4.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.UserNm, "t1.USER_NM")
	condition = utils.StringWhereConvert(condition, search.Department, "t1.DEPARTMENT")

	var order string
	if page.Order.Valid {
		order = page.Order.String
	} else {
		order = "DEPARTMENT ASC, USER_ID DESC"
	}

	query := fmt.Sprintf(`
				SELECT *
				FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						   SELECT
								MAX(t1.DNO) DNO
								,MAX(t1.SNO) SNO
								,MAX(t2.SITE_NM) SITE_NM
								,MAX(t1.JNO) JNO
								,MAX(t4.JOB_NAME) JOB_NAME
								,t1.USER_ID USER_ID
								,MAX(t1.USER_NM) USER_NM
								,MAX(t1.DEPARTMENT) DEPARTMENT
								,MIN(CASE WHEN t1.TRANS_TYPE = 'Clock in' THEN t1.RECOG_TIME END) AS IN_RECOG_TIME
								,MAX(CASE WHEN t1.TRANS_TYPE = 'Clock out' THEN t1.RECOG_TIME END) AS OUT_RECOG_TIME
							FROM
								IRIS_RECD_SET t1
								INNER JOIN IRIS_SITE_SET t2 ON t1.SNO = t2.SNO
								INNER JOIN IRIS_SITE_JOB t3 ON t1.JNO = t3.JNO
								INNER JOIN S_JOB_INFO t4 ON t3.JNO = t4.JNO
							WHERE
								t1.sno > 100
							AND t1.SNO = :1
								%s
							GROUP BY
								t1.USER_ID, t1.USER_GUID
							ORDER BY %s
					) sorted_data
					WHERE ROWNUM <= :2
				)
				WHERE RNUM > :3`, condition, order)

	if err := db.SelectContext(ctx, &sqls, query, search.Sno, page.EndNum, page.StartNum); err != nil {
		return nil, fmt.Errorf("GetWorkerSiteBaseList err: %v", err)
	}

	return &sqls, nil
}

// func: 현장 근로자 개수 조회
// @param
// - searchTime string: 조회 날짜
func (r *Repository) GetWorkerSiteBaseCount(ctx context.Context, db Queryer, search entity.WorkerSql) (int, error) {
	var count int

	condition := ""

	condition = utils.StringWhereConvert(condition, search.JobName, "t4.JOB_NAME")
	condition = utils.StringWhereConvert(condition, search.UserNm, "t1.USER_NM")
	condition = utils.StringWhereConvert(condition, search.Department, "t1.DEPARTMENT")

	query := fmt.Sprintf(`
				SELECT 
				    count (DISTINCT t1.USER_ID || '-' || t1.USER_GUID)
				FROM
					IRIS_RECD_SET t1
					INNER JOIN IRIS_SITE_SET t2 ON t1.SNO = t2.SNO
					INNER JOIN IRIS_SITE_JOB t3 ON t1.JNO = t3.JNO
					INNER JOIN S_JOB_INFO t4 ON t3.JNO = t4.JNO
				WHERE
					t1.sno > 100
				AND t1.SNO = :1
					%s`, condition)

	if err := db.GetContext(ctx, &count, query, search.Sno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("GetWorkerSiteBaseCount fail: %w", err)
	}
	return count, nil
}
