package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
func (r *Repository) GetWorkerTotalList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerSql) (*entity.WorkerSqls, error) {
	sqls := entity.WorkerSqls{}

	whereClause := ""
	if search.SiteNm.Valid {
		trimmedSiteNm := strings.TrimSpace(search.SiteNm.String)
		if trimmedSiteNm != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t2.SITE_NM) LIKE LOWER('%%%s%%')", trimmedSiteNm)
		}
	}

	if search.JobName.Valid {
		trimmedJobName := strings.TrimSpace(search.JobName.String)
		if trimmedJobName != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t4.JOB_NAME) LIKE LOWER('%%%s%%')", trimmedJobName)
		}
	}
	if search.UserNm.Valid {
		trimmedUserNm := strings.TrimSpace(search.UserNm.String)
		if trimmedUserNm != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t1.USER_NM) LIKE LOWER('%%%s%%')", trimmedUserNm)
		}
	}
	if search.Department.Valid {
		trimmedDepartment := strings.TrimSpace(search.Department.String)
		if trimmedDepartment != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t1.DEPARTMENT) LIKE LOWER('%%%s%%')", trimmedDepartment)
		}
	}

	// 기본 WHERE 조건 추가
	if whereClause == "" {
		whereClause = " AND 1=1" // WHERE 절이 비어있을 때 기본 조건 추가
	}

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
								TO_CHAR(t1.RECOG_TIME, 'YYYY-MM-DD') = :1
								%s
							GROUP BY
								t1.USER_ID, t1.USER_GUID
							ORDER BY %s
					) sorted_data
					WHERE ROWNUM <= :2
				)
				WHERE RNUM > :3`, whereClause, order)

	if err := db.SelectContext(ctx, &sqls, query, search.SearchTime, page.EndNum, page.StartNum); err != nil {
		return nil, fmt.Errorf("GetWorkerTotalList err: %v", err)
	}

	return &sqls, nil
}

// func: 전체 근로자 개수 조회
// @param
// - searchTime string: 조회 날짜
func (r *Repository) GetWorkerTotalCount(ctx context.Context, db Queryer, search entity.WorkerSql) (int, error) {
	var count int

	whereClause := ""
	if search.SiteNm.Valid {
		trimmedSiteNm := strings.TrimSpace(search.SiteNm.String)
		if trimmedSiteNm != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t2.SITE_NM) LIKE LOWER('%%%s%%')", trimmedSiteNm)
		}
	}

	if search.JobName.Valid {
		trimmedJobName := strings.TrimSpace(search.JobName.String)
		if trimmedJobName != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t4.JOB_NAME) LIKE LOWER('%%%s%%')", trimmedJobName)
		}
	}
	if search.UserNm.Valid {
		trimmedUserNm := strings.TrimSpace(search.UserNm.String)
		if trimmedUserNm != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t1.USER_NM) LIKE LOWER('%%%s%%')", trimmedUserNm)
		}
	}
	if search.Department.Valid {
		trimmedDepartment := strings.TrimSpace(search.Department.String)
		if trimmedDepartment != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t1.DEPARTMENT) LIKE LOWER('%%%s%%')", trimmedDepartment)
		}
	}

	// 기본 WHERE 조건 추가
	if whereClause == "" {
		whereClause = " AND 1=1" // WHERE 절이 비어있을 때 기본 조건 추가
	}

	query := fmt.Sprintf(`
				SELECT 
				    count (DISTINCT t1.USER_ID || '-' || t1.USER_GUID)
				FROM
					IRIS_RECD_SET t1
					INNER JOIN IRIS_SITE_SET t2 ON t1.SNO = t2.SNO
					INNER JOIN IRIS_SITE_JOB t3 ON t1.JNO = t3.JNO
					INNER JOIN S_JOB_INFO t4 ON t3.JNO = t4.JNO
				WHERE
					TO_CHAR(t1.RECOG_TIME, 'YYYY-MM-DD') = :1
					%s`, whereClause)

	if err := db.GetContext(ctx, &count, query, search.SearchTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("GetWorkerTotalCount fail: %w", err)
	}
	return count, nil
}

// func: 현장 근로자 조회
// @param
// - page entity.PageSql: 정렬, 리스트 수
// - search entity.WorkerSql: 검색 단어
func (r *Repository) GetWorkerSiteBaseList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerSql) (*entity.WorkerSqls, error) {
	sqls := entity.WorkerSqls{}

	whereClause := ""
	if search.SiteNm.Valid {
		trimmedSiteNm := strings.TrimSpace(search.SiteNm.String)
		if trimmedSiteNm != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t2.SITE_NM) LIKE LOWER('%%%s%%')", trimmedSiteNm)
		}
	}

	if search.JobName.Valid {
		trimmedJobName := strings.TrimSpace(search.JobName.String)
		if trimmedJobName != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t4.JOB_NAME) LIKE LOWER('%%%s%%')", trimmedJobName)
		}
	}
	if search.UserNm.Valid {
		trimmedUserNm := strings.TrimSpace(search.UserNm.String)
		if trimmedUserNm != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t1.USER_NM) LIKE LOWER('%%%s%%')", trimmedUserNm)
		}
	}
	if search.Department.Valid {
		trimmedDepartment := strings.TrimSpace(search.Department.String)
		if trimmedDepartment != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t1.DEPARTMENT) LIKE LOWER('%%%s%%')", trimmedDepartment)
		}
	}

	// 기본 WHERE 조건 추가
	if whereClause == "" {
		whereClause = " AND 1=1" // WHERE 절이 비어있을 때 기본 조건 추가
	}

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
								TO_CHAR(t1.RECOG_TIME, 'YYYY-MM-DD') = :1
								AND t1.SNO = :2
								%s
							GROUP BY
								t1.USER_ID, t1.USER_GUID
							ORDER BY %s
					) sorted_data
					WHERE ROWNUM <= :3
				)
				WHERE RNUM > :4`, whereClause, order)

	if err := db.SelectContext(ctx, &sqls, query, search.SearchTime, search.Sno, page.EndNum, page.StartNum); err != nil {
		return nil, fmt.Errorf("GetWorkerSiteBaseList err: %v", err)
	}

	return &sqls, nil
}

// func: 현장 근로자 개수 조회
// @param
// - searchTime string: 조회 날짜
func (r *Repository) GetWorkerSiteBaseCount(ctx context.Context, db Queryer, search entity.WorkerSql) (int, error) {
	var count int

	whereClause := ""
	if search.SiteNm.Valid {
		trimmedSiteNm := strings.TrimSpace(search.SiteNm.String)
		if trimmedSiteNm != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t2.SITE_NM) LIKE LOWER('%%%s%%')", trimmedSiteNm)
		}
	}

	if search.JobName.Valid {
		trimmedJobName := strings.TrimSpace(search.JobName.String)
		if trimmedJobName != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t4.JOB_NAME) LIKE LOWER('%%%s%%')", trimmedJobName)
		}
	}
	if search.UserNm.Valid {
		trimmedUserNm := strings.TrimSpace(search.UserNm.String)
		if trimmedUserNm != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t1.USER_NM) LIKE LOWER('%%%s%%')", trimmedUserNm)
		}
	}
	if search.Department.Valid {
		trimmedDepartment := strings.TrimSpace(search.Department.String)
		if trimmedDepartment != "" {
			whereClause += fmt.Sprintf(" AND LOWER(t1.DEPARTMENT) LIKE LOWER('%%%s%%')", trimmedDepartment)
		}
	}

	// 기본 WHERE 조건 추가
	if whereClause == "" {
		whereClause = " AND 1=1" // WHERE 절이 비어있을 때 기본 조건 추가
	}

	query := fmt.Sprintf(`
				SELECT 
				    count (DISTINCT t1.USER_ID || '-' || t1.USER_GUID)
				FROM
					IRIS_RECD_SET t1
					INNER JOIN IRIS_SITE_SET t2 ON t1.SNO = t2.SNO
					INNER JOIN IRIS_SITE_JOB t3 ON t1.JNO = t3.JNO
					INNER JOIN S_JOB_INFO t4 ON t3.JNO = t4.JNO
				WHERE
					TO_CHAR(t1.RECOG_TIME, 'YYYY-MM-DD') = :1
					AND t1.SNO = :2
					%s`, whereClause)

	if err := db.GetContext(ctx, &count, query, search.SearchTime, search.Sno); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("GetWorkerSiteBaseCount fail: %w", err)
	}
	return count, nil
}
