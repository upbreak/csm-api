package entity

import "time"

type ProjectDaily struct {
	Jno     int64     `json:"jno" db:"JNO"`
	Content string    `json:"content" db:"CONTENT"`
	IsUse   string    `json:"isUse" db:"IS_USE"`
	RegDate time.Time `json:"regDate" db:"REG_DATE"`
	RegUno  int64     `json:"regUno" db:"REG_UNO"`
	RegUser string    `json:"regUser" db:"REG_USER"`
	ModDate time.Time `json:"modDate" db:"MOD_DATE"`
	ModUno  int64     `json:"modUno" db:"MOD_UNO"`
	ModUser string    `json:"modUser" db:"MOD_USER"`
}

type ProjectDailys []*ProjectDaily
