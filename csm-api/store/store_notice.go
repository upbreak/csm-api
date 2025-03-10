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
func (r *Repository) GetNoticeList(ctx context.Context, db Queryer, page entity.PageSql, search entity.NoticeSql) (*entity.NoticeSqls, error) {
	sqls := entity.NoticeSqls{}

	// 조건
	condition := "1=1"
	if search.LocName.Valid {
		trimLocName := strings.TrimSpace(search.LocName.String)

		if trimLocName != "" {
			condition += fmt.Sprintf(" AND UPPER(LOC_NAME) LIKE UPPER('%%%s%%')", trimLocName)
		}
	}
	if search.SiteNm.Valid {
		trimSiteNm := strings.TrimSpace(search.SiteNm.String)
		if trimSiteNm != "" {
			condition += fmt.Sprintf(" AND UPPER(SITE_NM) LIKE UPPER('%%%s%%')", trimSiteNm)
		}
	}
	if search.Title.Valid {
		trimTitle := strings.TrimSpace(search.Title.String)
		if trimTitle != "" {
			condition += fmt.Sprintf(" AND UPPER(TITLE) LIKE UPPER('%%%s%%')", trimTitle)
		}
	}
	if search.UserInfo.Valid {
		trimUserInfo := strings.TrimSpace(search.UserInfo.String)
		if trimUserInfo != "" {
			condition += fmt.Sprintf(" AND UPPER(USER_INFO) LIKE UPPER('%%%s%%')", trimUserInfo)
		}
	}

	var order string
	if page.Order.Valid {
		order = page.Order.String
	} else {
		order = "REG_DATE DESC"
	}

	query := fmt.Sprintf(`
				WITH Notice AS (
					SELECT 
						N.IDX,
						N.SNO, 
						S.SITE_NM,
						S.LOC_NAME,
						N.TITLE, 
						N.CONTENT, 
						N.SHOW_YN,
						N.REG_UNO, 
						N.REG_USER, 
						N.REG_DATE,
						U.DUTY_NAME, 
						N.REG_USER || ' ' || U.DUTY_NAME as USER_INFO, 
						N.MOD_USER, 
						N.MOD_DATE,
						N.POSTING_PERIOD AS PERIOD_CODE,
						N.POSTING_DATE,
						C.CODE_NM AS NOTICE_NM
					FROM 
						IRIS_NOTICE_BOARD N 
					INNER JOIN
						S_SYS_USER_SET U ON N.REG_UNO = U.UNO
					LEFT OUTER JOIN 
						IRIS_SITE_SET S ON	N.SNO = S.SNO
					INNER JOIN
						IRIS_CODE_SET C ON N.POSTING_PERIOD = C.CODE AND C.P_CODE = 'NOTICE_PERIOD'
					WHERE
						N.IS_USE = 'Y'
						AND N.POSTING_DATE > SYSDATE
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
							%s
						) sorted_data
					WHERE ROWNUM <= :1
			  	)
			  	WHERE RNUM > :2`,
		condition, order)

	if err := db.SelectContext(ctx, &sqls, query, page.EndNum, page.StartNum); err != nil {
		fmt.Println("store/notice. NoticeList error")
		return nil, err
	}
	return &sqls, nil

}

// func: 공지사항 전체 개수 조회
// @param
// -
func (r *Repository) GetNoticeListCount(ctx context.Context, db Queryer, search entity.NoticeSql) (int, error) {
	var count int

	condition := "1=1"
	if search.LocName.Valid {
		trimLocName := strings.TrimSpace(search.LocName.String)

		if trimLocName != "" {
			condition += fmt.Sprintf(" AND UPPER(LOC_NAME) LIKE UPPER('%%%s%%')", trimLocName)
		}
	}
	if search.SiteNm.Valid {
		trimSiteNm := strings.TrimSpace(search.SiteNm.String)
		if trimSiteNm != "" {
			condition += fmt.Sprintf(" AND UPPER(SITE_NM) LIKE UPPER('%%%s%%')", trimSiteNm)
		}
	}
	if search.Title.Valid {
		trimTitle := strings.TrimSpace(search.Title.String)
		if trimTitle != "" {
			condition += fmt.Sprintf(" AND UPPER(TITLE) LIKE UPPER('%%%s%%')", trimTitle)
		}
	}
	if search.UserInfo.Valid {
		trimUserInfo := strings.TrimSpace(search.UserInfo.String)
		if trimUserInfo != "" {
			condition += fmt.Sprintf(" AND UPPER(USER_INFO) LIKE UPPER('%%%s%%')", trimUserInfo)
		}
	}

	query := fmt.Sprintf(`
			WITH Notice AS (
				SELECT 
					N.IDX,
					N.SNO, 
					S.SITE_NM,
					S.LOC_NAME,
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
					IRIS_SITE_SET S ON	N.SNO = S.SNO
				INNER JOIN
					IRIS_CODE_SET C ON N.POSTING_PERIOD = C.CODE AND C.P_CODE = 'NOTICE_PERIOD'
				WHERE
					N.IS_USE = 'Y'
					AND N.POSTING_DATE > SYSDATE
			)
			SELECT COUNT(*) 
			FROM  Notice
			WHERE
				%s`, condition)

	if err := db.GetContext(ctx, &count, query); err != nil {
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
func (r *Repository) AddNotice(ctx context.Context, db Beginner, noticeSql entity.NoticeSql) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("store/notice. Failed to begin transaction: %w", err)
	}

	query := `
				INSERT INTO IRIS_NOTICE_BOARD(
					IDX,
					SNO,
					TITLE,
					CONTENT,
					SHOW_YN,
					IS_USE,
					REG_UNO,
					REG_USER,
					REG_DATE,
					POSTING_PERIOD,
					POSTING_DATE
				)
				SELECT
					SEQ_IRIS_NOTICE_BOARD.NEXTVAL,
					:1,
					:2,
					:3,
					:4,
					'Y',
					:5,
					:6,
					SYSDATE,
					C.CODE,
					ADD_MONTHS(SYSDATE, C.UDF_VAL_03) + C.UDF_VAL_04
				FROM IRIS_CODE_SET C
				WHERE C.P_CODE = 'NOTICE_PERIOD' AND C.CODE = :7 
		`

	_, err = tx.ExecContext(ctx, query, noticeSql.Sno, noticeSql.Title, noticeSql.Content, noticeSql.ShowYN, noticeSql.RegUno, noticeSql.RegUser, noticeSql.PeriodCode)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("store/notice. AddNotice fail %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store/notice. failed to commit transaction: %v", err)
	}

	return nil
}

// func: 공지사항 수정
// @param
// - notice entity.NoticeSql: IDX, SNO, TITLE, CONTENT, SHOW_YN, MOD_UNO, MOD_USER
func (r *Repository) ModifyNotice(ctx context.Context, db Beginner, noticeSql entity.NoticeSql) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("store/notice. Failed to begin transaction: %w", err)
	}

	query := `
				UPDATE IRIS_NOTICE_BOARD
				SET
					SNO = :1,
					TITLE = :2,
					CONTENT = :3,
					SHOW_YN = :4,
					IS_USE = 'Y',
					MOD_UNO = :5,	
					MOD_USER = :6,
					MOD_DATE = SYSDATE,
					(POSTING_PERIOD,
					POSTING_DATE) = (
						SELECT
							C.CODE, ADD_MONTHS(N.REG_DATE, C.UDF_VAL_03) + C.UDF_VAL_04
						FROM 
							IRIS_CODE_SET C
						INNER JOIN
							IRIS_NOTICE_BOARD N ON N.IDX = :7
						WHERE C.CODE = :8 AND C.P_CODE = 'NOTICE_PERIOD'
					)
				WHERE 
					IDX = :9
			`

	_, err = tx.ExecContext(ctx, query, noticeSql.Sno, noticeSql.Title, noticeSql.Content, noticeSql.ShowYN, noticeSql.ModUno, noticeSql.ModUser, noticeSql.Idx, noticeSql.PeriodCode, noticeSql.Idx)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("store/notice. ModifyNotice fail: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store/notice. failed to commit transaction: %v", err)
	}

	return nil
}

// func: 공지사항 삭제
// @param
// - idx: 공지사항 인덱스
func (r *Repository) RemoveNotice(ctx context.Context, db Beginner, idx entity.NoticeID) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		fmt.Println("store/notice. Failed to begint transaction: %w", err)
	}

	query := `
		UPDATE IRIS_NOTICE_BOARD 
		SET 
			IS_USE = 'N'
		WHERE 
			IDX = :1
			`

	_, err = tx.ExecContext(ctx, query, idx)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("store/notice. RemoveNotice fail: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store/notice. failed to commit transaction: %v", err)
	}

	return nil
}

// func: 공지기간 조회
// @param
// -
func (r *Repository) GetNoticePeriod(ctx context.Context, db Queryer) (*entity.NoticePeriodSqls, error) {
	periodSqls := entity.NoticePeriodSqls{}

	query := fmt.Sprintf(`
		SELECT
			CODE AS PERIOD_CODE,
			CODE_NM AS NOTICE_NM
		FROM
			IRIS_CODE_SET
		WHERE
			P_CODE = 'NOTICE_PERIOD'
	`)

	if err := db.SelectContext(ctx, &periodSqls, query); err != nil {
		return &entity.NoticePeriodSqls{}, fmt.Errorf("store/notice. GetNoticePeriod %w", err)
	}

	return &periodSqls, nil

}
