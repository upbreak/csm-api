package store

import "C"
import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/godror/godror"
	"github.com/guregu/null"
	"strings"
)

/**
 * @author 작성자: 정지영
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * - 검색기능 추가
 * - UserInfo 추가
 */

// func: 공지사항 전체 조회
// @param
// - page entity.PageSql : 현재페이지 번호, 리스트 목록 개수
func (r *Repository) GetNoticeList(ctx context.Context, db Queryer, uno null.Int, role int, page entity.PageSql, search entity.Notice) (*entity.Notices, error) {
	notices := entity.Notices{}

	// 조건
	condition := "1=1"
	condition = utils.Int64WhereConvert2(condition, search.Jno, "JNO")
	condition = utils.StringWhereConvert2(condition, search.JobLocName, "JOB_LOC_NAME")
	condition = utils.StringWhereConvert2(condition, search.JobName, "JOB_NAME")
	condition = utils.StringWhereConvert2(condition, search.Title, "TITLE")
	condition = utils.StringWhereConvert2(condition, search.UserInfo, "USER_INFO")

	var order string
	if page.Order.Valid {
		order = page.Order.String
	} else {
		order = "NULL"
	}

	query := fmt.Sprintf(`
				WITH Notice AS (
					SELECT 
						N.IDX,
						N.JNO, 
						DECODE(J.JOB_NAME, 'NONE', '전체', J.JOB_NAME) AS JOB_NAME,
						J.JOB_LOC_NAME,
						N.TITLE, 
						N.CONTENT, 
						N.SHOW_YN,
						N.REG_UNO, 
						N.REG_USER, 
						N.REG_DATE,
						U.DUTY_NAME, 
						N.REG_USER || DECODE(N.REG_USER, '관리자', '',  ' ' || U.DUTY_NAME) AS USER_INFO, 
						N.MOD_USER, 
						N.MOD_DATE,
						N.POSTING_PERIOD AS PERIOD_CODE,
						N.POSTING_START_DATE,
						N.POSTING_END_DATE,
						C.CODE_NM AS NOTICE_NM,
						N.IS_IMPORTANT
					FROM 
						IRIS_NOTICE_BOARD N 
					INNER JOIN
						S_SYS_USER_SET U ON N.REG_UNO = U.UNO
					LEFT OUTER JOIN 
						S_JOB_INFO J ON J.JNO = N.JNO
					LEFT OUTER JOIN
						IRIS_CODE_SET C ON N.POSTING_PERIOD = C.CODE AND C.P_CODE = 'NOTICE_PERIOD'
					WHERE
						N.IS_USE = 'Y'
						AND N.POSTING_START_DATE <= SYSDATE
						AND N.POSTING_END_DATE > SYSDATE
						AND (N.JNO IN (SELECT DISTINCT(JNO) FROM TIMESHEET.JOB_MEMBER_LIST WHERE 1 = :1 OR UNO = :2) OR N.JNO = 0 )
				)
				SELECT * 
			  	FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						SELECT *
						FROM Notice
						WHERE
							%s
						ORDER BY
							%s,
							CASE WHEN 
									IS_IMPORTANT= 'Y' 
								THEN 0 
								ELSE 1
							END,
							CASE WHEN
									JNO = 0 
								THEN 0
								ELSE 1 
							END,
							POSTING_START_DATE DESC
						) sorted_data
					WHERE ROWNUM <= :3
			  	)
			  	WHERE RNUM > :4`,
		condition, order)

	if err := db.SelectContext(ctx, &notices, query, role, uno, page.EndNum, page.StartNum); err != nil {
		//TODO: 에러 아카이브 처리
		fmt.Printf("store/notice. NoticeList error %s", err)
		return nil, err
	}
	return &notices, nil
}

// func: 공지사항 전체 개수 조회
// @param
// -
func (r *Repository) GetNoticeListCount(ctx context.Context, db Queryer, uno null.Int, role int, search entity.Notice) (int, error) {
	var count int

	condition := "1=1"
	condition = utils.Int64WhereConvert2(condition, search.Jno, "JNO")
	condition = utils.StringWhereConvert2(condition, search.JobLocName, "JOB_LOC_NAME")
	condition = utils.StringWhereConvert2(condition, search.JobName, "JOB_NAME")
	condition = utils.StringWhereConvert2(condition, search.Title, "TITLE")
	condition = utils.StringWhereConvert2(condition, search.UserInfo, "USER_INFO")

	query := fmt.Sprintf(`
			WITH Notice AS (
				SELECT 
					N.IDX,
					N.JNO, 
					J.JOB_NAME,
					J.JOB_LOC_NAME,
					N.TITLE, 
					N.CONTENT, 
					N.SHOW_YN,
					N.REG_UNO, 
					N.REG_USER, 
					N.REG_DATE,
					U.DUTY_NAME, 
					N.REG_USER || ' ' || U.DUTY_NAME as USER_INFO, 
					N.MOD_USER, 
					N.MOD_DATE 
				FROM 
					IRIS_NOTICE_BOARD N 
				INNER JOIN
					S_SYS_USER_SET U ON N.REG_UNO = U.UNO
				LEFT OUTER JOIN 
					S_JOB_INFO J ON J.JNO = N.JNO
				LEFT OUTER JOIN
					IRIS_CODE_SET C ON N.POSTING_PERIOD = C.CODE AND C.P_CODE = 'NOTICE_PERIOD'
				WHERE
					N.IS_USE = 'Y'
					AND N.POSTING_START_DATE <= SYSDATE
					AND N.POSTING_END_DATE > SYSDATE
					AND (N.JNO IN (SELECT DISTINCT(JNO) FROM TIMESHEET.JOB_MEMBER_LIST WHERE 1 = :1 OR UNO = :2) OR N.JNO = 0 )
			)
			SELECT COUNT(*) 
			FROM  Notice
			WHERE
				%s`, condition)

	if err := db.GetContext(ctx, &count, query, role, uno); err != nil {
		//TODO: 에러 아카이브 처리
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("store/notice. GetNoticeListCount fail: %w", err)
	}
	return count, nil

}

// func: 공지사항 추가
// @param
// - notice entity.NoticeSql: SNO, TITLE, CONTENT, SHOW_YN, REG_UNO, REG_USER
func (r *Repository) AddNotice(ctx context.Context, db Beginner, notice entity.Notice) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		//TODO: 에러 아카이브 처리
		fmt.Println("AddNotice Begin transaction fail: %w", err)
	}

	// 공지사항
	contentCLOB := godror.Lob{
		IsClob: true,
		Reader: strings.NewReader(notice.Content.String),
	}

	query := `
				INSERT INTO IRIS_NOTICE_BOARD(
					IDX,
					SNO,
					JNO,
					TITLE,
					CONTENT,
					SHOW_YN,
					IS_USE,
				    IS_IMPORTANT,
					REG_UNO,
					REG_USER,
					REG_DATE,
					POSTING_START_DATE,
					POSTING_END_DATE,
				    REG_USER_DUTY_NAME
				) VALUES (
					SEQ_IRIS_NOTICE_BOARD.NEXTVAL,
					(SELECT S.SNO FROM IRIS_SITE_JOB S RIGHT JOIN S_JOB_INFO J ON S.JNO = J.JNO WHERE J.JNO = :1 AND S.IS_USE = 'Y'),
					:2,
					:3,
					:4,
					:5,
					'Y',
					:6,
					:7,
					:8,
					SYSDATE,
--					C.CODE,
					:9,
					:10,
--					ADD_MONTHS(SYSDATE, C.UDF_VAL_03) + C.UDF_VAL_04,
					(SELECT U.DUTY_NAME FROM S_SYS_USER_SET U WHERE U.UNO = :11)
				)
--				FROM IRIS_CODE_SET C
--				WHERE C.P_CODE = 'NOTICE_PERIOD' AND C.CODE = :9
		`

	_, origErr := tx.ExecContext(ctx, query, notice.Jno, notice.Jno, notice.Title, contentCLOB, notice.ShowYN, notice.IsImportant, notice.RegUno, notice.RegUser, notice.PostingStartDate, notice.PostingEndDate, notice.RegUno)

	if origErr != nil {
		//TODO: 에러 아카이브 처리
		if err = tx.Rollback(); err != nil {

			//TODO: 에러 아카이브 처리
			return fmt.Errorf("AddNotice Rollback transaction fail: %w", err)
		}
		return fmt.Errorf("store/notice. AddNotice fail %v", origErr)
	}

	if err = tx.Commit(); err != nil {
		//TODO: 에러 아카이브 처리
		return fmt.Errorf("AddNotice Commit transaction fail: %v", err)
	}

	return nil
}

// func: 공지사항 수정
// @param
// - notice entity.NoticeSql: IDX, SNO, TITLE, CONTENT, SHOW_YN, MOD_UNO, MOD_USER
func (r *Repository) ModifyNotice(ctx context.Context, db Beginner, notice entity.Notice) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(
			//TODO: 에러 아카이브 처리
			"ModifyNotice Begin transaction fail: %w", err)
	}

	query := `
				UPDATE IRIS_NOTICE_BOARD
				SET
				    SNO = (SELECT S.SNO FROM IRIS_SITE_JOB S RIGHT JOIN S_JOB_INFO J ON S.JNO = J.JNO WHERE J.JNO = :1 AND S.IS_USE = 'Y'),
					JNO = :2,
					TITLE = :3,
					CONTENT = :4,
					SHOW_YN = :5,
					IS_USE = 'Y',
				    IS_IMPORTANT = :6,
					MOD_UNO = :7,	
					MOD_USER = :8,
					MOD_DATE = SYSDATE,
					POSTING_START_DATE = :9,
					POSTING_END_DATE = :10
				WHERE 
					IDX = :11
			`

	_, origErr := tx.ExecContext(ctx, query, notice.Jno, notice.Jno, notice.Title, notice.Content, notice.ShowYN, notice.IsImportant, notice.ModUno, notice.ModUser, notice.PostingStartDate, notice.PostingEndDate, notice.Idx)

	if origErr != nil {
		if err = tx.Rollback(); err != nil {
			//TODO: 에러 아카이브 처리
			return fmt.Errorf("ModifyNotice Rollback transaction fail: %w", err)
		}
		//TODO: 에러 아카이브 처리
		return fmt.Errorf("store/notice. ModifyNotice fail: %v", origErr)
	}

	if err = tx.Commit(); err != nil {
		//TODO: 에러 아카이브 처리
		return fmt.Errorf("ModifyNotice Commit transaction fail: %v", err)
	}

	return nil
}

// func: 공지사항 삭제
// @param
// - idx: 공지사항 인덱스
func (r *Repository) RemoveNotice(ctx context.Context, db Beginner, idx null.Int) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		//TODO: 에러 아카이브 처리
		fmt.Println("store/notice. Failed to begint transaction: %w", err)
	}

	query := `
		UPDATE IRIS_NOTICE_BOARD 
		SET 
			IS_USE = 'N'
		WHERE 
			IDX = :1
			`

	_, origErr := tx.ExecContext(ctx, query, idx)

	if origErr != nil {
		if err = tx.Rollback(); err != nil {
			//TODO: 에러 아카이브 처리
			return fmt.Errorf("RemoveNotice Rollback transaction fail: %w", err)
		}
		//TODO: 에러 아카이브 처리
		return fmt.Errorf("store/notice. RemoveNotice fail: %v", origErr)
	}

	if err = tx.Commit(); err != nil {
		//TODO: 에러 아카이브 처리
		return fmt.Errorf("RemoveNotice Commit transaction fail: %v", err)
	}

	return nil
}
