package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

type ListNotice struct {
	DB   store.Queryer
	Repo store.NoticeAddStore
}

func (l *ListNotice) GetNoticeList(ctx context.Context) ([]entity.Notice, error) {
	notices, err := l.Repo.GetNoticeList(ctx, l.DB)
	if err != nil {
		return nil, fmt.Errorf("fail to list notice: %w", err)
	}

	// IDX      NoticeID  `json:"idx" db:"IDX"`
	// SNO      int64     `json:"sno" db:"SNO"`
	// TITLE    string    `json:"title" db:"TITLE" validate:"required"`
	// CONTENT  string    `json:"content" db:"CONTENT" validate:"required"`
	// SHOW_YN  string    `json:"show_yn" db:"SHOW_YN"`
	// REG_UNO  int64     `json:"reg_uno" db:"REG_UNO" validate:"required"`
	// REG_USER string    `json:"reg_user" db:"REG_USER" validate:"required"`
	// REG_DATE time.Time `json:"reg_date" db:"REG_DATE"`
	// MOD_UNO  int64     `json:"mod_uno" db:"MOD_UNO" validate:"required"`
	// MOD_USER string    `json:"mod_user" db:"MOD_USER" validate:"required"`
	// MOD_DATE

	var rsp []entity.Notice
	for _, n := range notices {
		rsp = append(rsp, entity.Notice{
			IDX:      n.IDX,
			TITLE:    n.TITLE,
			CONTENT:  n.CONTENT,
			REG_UNO:  n.REG_UNO,
			REG_USER: n.REG_USER,
			REG_DATE: n.REG_DATE,
			MOD_USER: n.MOD_USER,
			MOD_DATE: n.MOD_DATE,
		})
	}

	return rsp, nil
}
