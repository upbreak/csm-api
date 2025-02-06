package entity

import "database/sql"

type Code struct {
	Code   string `json:"code"`
	PCode  string `json:"p_code"`
	CodeNm string `json:"code_nm"`
}

type Codes []*Code

type CodeSql struct {
	Code   sql.NullString `db:"CODE"`
	PCode  sql.NullString `db:"P_CODE"`
	CodeNm sql.NullString `db:"CODE_NM"`
}

type CodeSqls []*CodeSql

func (c *Code) ToCode(sql *CodeSql) *Code {
	c.Code = sql.Code.String
	c.PCode = sql.PCode.String
	c.CodeNm = sql.CodeNm.String

	return c
}

func (c *Codes) ToCodes(sqls *CodeSqls) *Codes {
	for _, sql := range *sqls {
		code := Code{}
		code.ToCode(sql)
		*c = append(*c, &code)
	}
	return c
}
