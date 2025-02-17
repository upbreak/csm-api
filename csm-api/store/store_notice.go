package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"errors"
	"fmt"
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
func (r *Repository) GetNoticeList(ctx context.Context, db Queryer, page entity.PageSql) (*entity.NoticeSqls, error) {
	sqls := entity.NoticeSqls{}
	query := `SELECT * 
			  	FROM (
					SELECT ROWNUM AS RNUM, sorted_data.*
					FROM (
						SELECT 
							n1.IDX,
							n1.SNO, 
							n2.SITE_NM,
							n2.LOC_CODE,
							n1.TITLE, 
							n1.CONTENT, 
							n1.REG_UNO, 
							n1.REG_USER, 
							n1.REG_DATE, 
							n1.MOD_USER, 
							n1.MOD_DATE 
						FROM 
							IRIS_NOTICE_BOARD n1 LEFT OUTER JOIN IRIS_SITE_SET n2 
						ON 
							n1.SNO = n2.SNO
						WHERE 
							n1.IS_USE = 'Y'
						ORDER BY 
							n1.REG_DATE DESC
						) sorted_data
					WHERE ROWNUM <= :1
			  	)
			  	WHERE RNUM > :2`
	if err := db.SelectContext(ctx, &sqls, query, page.EndNum, page.StartNum); err != nil {
		fmt.Println("store/notice. NoticeList error")
		return nil, err
	}
	return &sqls, nil

}

// func: 공지사항 전체 개수 조회
// @param
// -
func (r *Repository) GetNoticeListCount(ctx context.Context, db Queryer) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM IRIS_NOTICE_BOARD WHERE IS_USE = 'Y'`

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
func (r *Repository) AddNotice(ctx context.Context, db Beginner, noticeSql *entity.NoticeSql) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Failed to begin transaction: %v", err)
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
