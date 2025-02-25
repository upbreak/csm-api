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
 * -
 */

// func: 공지사항 전체 조회
// @param
// - page entity.PageSql : 현재페이지 번호, 리스트 목록 개수
func (r *Repository) GetNoticeList(ctx context.Context, db Queryer, page entity.PageSql, search entity.NoticeSql) (*entity.NoticeSqls, error) {
	sqls := entity.NoticeSqls{}

	// 조건
	condition := "AND 1=1"
	if search.LocCode.Valid {
		trimLocCode := strings.TrimSpace(search.LocCode.String)

		if trimLocCode != "" {
			condition += fmt.Sprintf(" AND UPPER(S.LOC_CODE) LIKE UPPER('%%%s%%')", trimLocCode)
		}
	}
	if search.SiteNm.Valid {
		trimSiteNm := strings.TrimSpace(search.SiteNm.String)
		if trimSiteNm != "" {
			condition += fmt.Sprintf(" AND UPPER(S.SITE_NM) LIKE UPPER('%%%s%%')", trimSiteNm)
		}
	}
	if search.Title.Valid {
		trimTitle := strings.TrimSpace(search.Title.String)
		if trimTitle != "" {
			condition += fmt.Sprintf(" AND UPPER(N.TITLE) LIKE UPPER('%%%s%%')", trimTitle)
		}
	}
	if search.RegUser.Valid {
		trimRegUser := strings.TrimSpace(search.RegUser.String)
		if trimRegUser != "" {
			condition += fmt.Sprintf(" AND UPPER(N.REG_USER) LIKE UPPER('%%%s%%')", trimRegUser)
		}
	}

	var order string
	if page.Order.Valid {
		order = page.Order.String
	} else {
		order = "N.REG_DATE DESC"
	}

	query := fmt.Sprintf(`
				SELECT * 
			  	FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						SELECT 
							N.IDX,
							N.SNO, 
							S.SITE_NM,
							S.LOC_CODE,
							N.TITLE, 
							N.CONTENT, 
							N.SHOW_YN,
							N.REG_UNO, 
							N.REG_USER, 
							N.REG_DATE,
							U.DUTY_NAME, 
							N.MOD_USER, 
							N.MOD_DATE 
						FROM 
							IRIS_NOTICE_BOARD N 
						INNER JOIN
							S_SYS_USER_SET U ON N.REG_UNO = U.UNO
						LEFT OUTER JOIN 
							IRIS_SITE_SET S ON	N.SNO = S.SNO
						WHERE
							N.IS_USE = 'Y' 
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

	condition := "AND 1=1"
	if search.LocCode.Valid {
		trimLocCode := strings.TrimSpace(search.LocCode.String)

		if trimLocCode != "" {
			condition += fmt.Sprintf(" AND UPPER(n2.LOC_CODE) LIKE UPPER('%%%s%%')", trimLocCode)
		}
	}
	if search.SiteNm.Valid {
		trimSiteNm := strings.TrimSpace(search.SiteNm.String)
		if trimSiteNm != "" {
			condition += fmt.Sprintf(" AND UPPER(n2.SITE_NM) LIKE UPPER('%%%s%%')", trimSiteNm)
		}
	}
	if search.Title.Valid {
		trimTitle := strings.TrimSpace(search.Title.String)
		if trimTitle != "" {
			condition += fmt.Sprintf(" AND UPPER(n1.TITLE) LIKE UPPER('%%%s%%')", trimTitle)
		}
	}
	if search.RegUser.Valid {
		trimRegUser := strings.TrimSpace(search.RegUser.String)
		if trimRegUser != "" {
			condition += fmt.Sprintf(" AND UPPER(n1.REG_USER) LIKE UPPER('%%%s%%')", trimRegUser)
		}
	}

	query := fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM 
				IRIS_NOTICE_BOARD n1 LEFT OUTER JOIN IRIS_SITE_SET n2 
			ON 
				n1.SNO = n2.SNO 
			WHERE
				n1.IS_USE = 'Y' 
				%s`, condition)

	if err := db.GetContext(ctx, &count, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("GetNoticeListCount fail: %w", err)
	}
	return count, nil

}

// func: 공지사항 추가
// @param
// - notice entity.NoticeSql: SNO, TITLE, CONTENT, SHOW_YN, REG_UNO, REG_USER
func (r *Repository) AddNotice(ctx context.Context, db Beginner, noticeSql entity.NoticeSql) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Failed to begin transaction: %w", err)
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
					REG_DATE
				)
				VALUES (
					SEQ_IRIS_NOTICE_BOARD.NEXTVAL,
					:1,
					:2,
					:3,
					:4,
					'Y',
					:5,
					:6,
					SYSDATE
		)`

	_, err = tx.ExecContext(ctx, query, noticeSql.Sno, noticeSql.Title, noticeSql.Content, noticeSql.ShowYN, noticeSql.RegUno, noticeSql.RegUser)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("AddNotice fail %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// func: 공지사항 수정
// @param
// - notice entity.NoticeSql: IDX, SNO, TITLE, CONTENT, SHOW_YN, MOD_UNO, MOD_USER
func (r *Repository) ModifyNotice(ctx context.Context, db Beginner, noticeSql entity.NoticeSql) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Failed to begin transaction: %w", err)
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
					MOD_DATE = SYSDATE
				WHERE 
					IDX = :7
			`

	_, err = tx.ExecContext(ctx, query, noticeSql.Sno, noticeSql.Title, noticeSql.Content, noticeSql.ShowYN, noticeSql.ModUno, noticeSql.ModUser, noticeSql.Idx)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("ModifyNotice fail: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// func: 공지사항 삭제
// @param
// - idx: 공지사항 인덱스
func (r *Repository) RemoveNotice(ctx context.Context, db Beginner, idx entity.NoticeID) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		fmt.Println("Failed to begint transaction: %w", err)
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
		return fmt.Errorf("RemoveNotice fail: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
