package entity

import (
	"github.com/guregu/null"
	"time"
)

type Deduction struct {
	Sno          null.Int    `json:"sno" db:"SNO"`
	Jno          null.Int    `json:"jno" db:"JNO"`
	UserNm       null.String `json:"user_nm" db:"USER_NM"`
	Department   null.String `json:"department" db:"DEPARTMENT"`
	Gender       null.String `json:"gender" db:"GENDER"`
	RegNo        null.String `json:"reg_no" db:"REG_NO"`
	Phone        null.String `json:"phone" db:"PHONE"`
	InRecogTime  null.Time   `json:"in_recog_time" db:"IN_RECOG_TIME"`
	OutRecogTime null.Time   `json:"out_recog_time" db:"OUT_RECOG_TIME"`
	RecordDate   null.Time   `json:"record_date" db:"RECORD_DATE"`
	DeductOrder  null.String `json:"deduct_order" db:"DEDUCT_ORDER"`
	Base
}

type DeductionKey struct {
	Jno        int64
	UserNm     string
	Department string
	RecordDate time.Time
}

type DeductionRegKey struct {
	Jno        int64
	RegNo      string // ' ', '-' 제거 주민번호
	RecordDate time.Time
}
