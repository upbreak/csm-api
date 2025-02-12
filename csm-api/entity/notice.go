package entity

import (
	"database/sql"
)

type NoticeID int64

type Notice struct {
	IDX      NoticeID       `json:"idx" db:"IDX"`
	SNO      sql.NullInt64  `json:"sno" db:"SNO"`
	TITLE    sql.NullString `json:"title" db:"TITLE" validate:"required"`
	CONTENT  sql.NullString `json:"content" db:"CONTENT" validate:"required"`
	SHOW_YN  sql.NullString `json:"show_yn" db:"SHOW_YN"`
	REG_UNO  sql.NullInt64  `json:"reg_uno" db:"REG_UNO" validate:"required"`
	REG_USER sql.NullString `json:"reg_user" db:"REG_USER" validate:"required"`
	REG_DATE sql.NullTime   `json:"reg_date" db:"REG_DATE"`
	MOD_UNO  sql.NullInt64  `json:"mod_uno" db:"MOD_UNO" validate:"required"`
	MOD_USER sql.NullString `json:"mod_user" db:"MOD_USER" validate:"required"`
	MOD_DATE sql.NullTime   `json:"mod_date" db:"MOD_DATE"`
}

type Notices []*Notice
