package entity

import "github.com/guregu/null"

type Code struct {
	Code      null.String `json:"code" db:"CODE"`
	PCode     null.String `json:"p_code" db:"P_CODE"`
	CodeNm    null.String `json:"code_nm" db:"CODE_NM"`
	CodeColor null.String `json:"code_color" db:"CODE_COLOR"`
}

type Codes []*Code
