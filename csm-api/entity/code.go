package entity

import "github.com/guregu/null"

type Code struct {
	Level     null.Int    `json:"level" db:"LEVEL"`
	IDX       null.Int    `json:"idx" db:"IDX"`
	Code      null.String `json:"code" db:"CODE"`
	PCode     null.String `json:"p_code" db:"P_CODE"`
	CodeNm    null.String `json:"code_nm" db:"CODE_NM"`
	CodeColor null.String `json:"code_color" db:"CODE_COLOR"`
	UdfVal03  null.String `json:"udf_val_03" db:"UDF_VAL_03"`
	UdfVal04  null.String `json:"udf_val_04" db:"UDF_VAL_04"`
	UdfVal05  null.String `json:"udf_val_05" db:"UDF_VAL_05"`
	UdfVal06  null.String `json:"udf_val_06" db:"UDF_VAL_06"`
	UdfVal07  null.String `json:"udf_val_07" db:"UDF_VAL_07"`
	SortNo    null.Int    `json:"sort_no" db:"SORT_NO"`
	IsUse     null.String `json:"is_use" db:"IS_USE"`
	Etc       null.String `json:"etc" db:"ETC"`
}

type Codes []*Code

type CodeTree struct {
	IDX      null.Int    `json:"idx" db:"IDX"` // level이 1이 아니면 쌓아야함.
	Code     null.String `json:"code" db:"CODE"`
	Level    null.Int    `json:"level" db:"LEVEL"`
	PCode    null.String `json:"p_code" db:"P_CODE"`
	Expand   null.Bool   `json:"expand" db:"EXPAND"`
	Children *CodeTrees  `json:"code_trees" db:"CODE_TREES"`
	CodeSet  *Code       `json:"code_set" db:"CODE_SET"`
}

type CodeTrees []*CodeTree
