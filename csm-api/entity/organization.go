package entity

import "database/sql"

type Organization struct {
	Jno          int64  `json:"jno"`
	FuncName     string `json:"func_name"`
	ChargeDetail string `json:"charge_detail"`
	UserName     string `json:"user_name"`
	DutyName     string `json:"duty_name"`
	DeptName     string `json:"dept_name"`
	Cell         string `json:"cell"`
	Tel          string `json:"tel"`
	Email        string `json:"email"`
	IsUse        string `json:"is_use"`
	CoId         string `json:"co_id"`
	CdNm         string `json:"cd_nm"`
	Uno          int64  `json:"uno"`
}

type Organizations []*Organization

type OrganizationSql struct {
	Jno          sql.NullInt64  `db:"JNO"`
	FuncName     sql.NullString `db:"FUNC_NAME"`
	ChargeDetail sql.NullString `db:"CHARGE_DETAIL"`
	UserName     sql.NullString `db:"USER_NAME"`
	DutyName     sql.NullString `db:"DUTY_NAME"`
	DeptName     sql.NullString `db:"DEPT_NAME"`
	Cell         sql.NullString `db:"CELL"`
	Tel          sql.NullString `db:"TEL"`
	Email        sql.NullString `db:"EMAIL"`
	IsUse        sql.NullString `db:"IS_USE"`
	CoId         sql.NullString `db:"CO_ID"`
	CdNm         sql.NullString `db:"CD_NM"`
	Uno          sql.NullInt64  `db:"UNO"`
}

type OrganizationSqls []*OrganizationSql
