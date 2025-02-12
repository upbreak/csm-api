package entity

import (
	"database/sql"
	"fmt"
)

type Page struct {
	PageNum int64 `json:"page_num"`
	RowSize int64 `json:"row_size"`
}

type PageSql struct {
	StartNum sql.NullInt64 `db:"page_num"`
	EndNum   sql.NullInt64 `db:"row_size"`
}

func (s PageSql) OfPageSql(p Page) (PageSql, error) {

	if p.PageNum != 0 && p.RowSize != 0 {
		s.StartNum = sql.NullInt64{Valid: true, Int64: (p.PageNum - 1) * p.RowSize}
		s.EndNum = sql.NullInt64{Valid: true, Int64: p.PageNum * p.RowSize}
	} else {
		return PageSql{}, fmt.Errorf("PageNum or RowSize is zero")
	}

	return s, nil
}
