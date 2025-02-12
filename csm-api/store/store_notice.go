package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

/*
	func (r *Repository) AddNotice(ctx context.Context, db Queryer) (entity.Notice, error) {

}
*/
func (r *Repository) GetNoticeList(ctx context.Context, db Queryer) (entity.Notices, error) {
	notices := entity.Notices{}
	sql := `SELECT 
				n1.IDX,
				n1.SNO, 
				n1.TITLE, 
				n1.CONTENT, 
				n1.REG_UNO, 
				n1.REG_USER, 
				n1.REG_DATE, 
				n1.MOD_USER, 
				n1.MOD_DATE 
			FROM 
				IRIS_NOTICE_BOARD n1
			WHERE 
				n1.IS_USE = 'Y'
			ORDER BY 
				n1.REG_DATE DESC`

	if err := db.SelectContext(ctx, &notices, sql); err != nil {
		fmt.Println("store/notice. NoticeList error")
		return nil, err
	}
	return notices, nil

}
