package entity

import (
	"github.com/guregu/null"
	"time"
)

type Tbm struct {
	Sno        null.Int    `json:"sno" db:"SNO"`
	Jno        null.Int    `json:"jno" db:"JNO"`
	Department null.String `json:"department" db:"DEPARTMENT"`
	DiscName   null.String `json:"disc_name" db:"DISC_NAME"`
	UserNm     null.String `json:"user_nm" db:"USER_NM"`
	TbmOrder   null.Int    `json:"tbm_order" db:"TBM_ORDER"`
	TbmDate    null.Time   `json:"tbm_date" db:"TBM_DATE"`
	Base
}

type TbmKey struct {
	Jno        int64
	UserNm     string
	Department string
	TbmDate    time.Time
}
