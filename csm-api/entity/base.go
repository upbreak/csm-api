package entity

import "github.com/guregu/null"

type Base struct {
	RegUser  null.String `json:"reg_user" db:"REG_USER"`
	RegAgent null.String `json:"reg_agent" db:"REG_AGENT"`
	RegUno   null.Int    `json:"reg_uno" db:"REG_UNO"`
	RegDate  null.Time   `json:"reg_date" db:"REG_DATE"`
	ModUser  null.String `json:"mod_user" db:"MOD_USER"`
	ModAgent null.String `json:"mod_agent" db:"MOD_AGENT"`
	ModUno   null.Int    `json:"mod_uno" db:"MOD_UNO"`
	ModDate  null.Time   `json:"mod_date" db:"MOD_DATE"`
}
