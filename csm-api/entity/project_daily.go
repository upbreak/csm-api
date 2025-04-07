package entity

import (
	"github.com/guregu/null"
)

type ProjectDaily struct {
	Jno     null.Int    `json:"jno" db:"JNO"`
	Content null.String `json:"content" db:"CONTENT"`
	IsUse   null.String `json:"isUse" db:"IS_USE"`
	Base
}

type ProjectDailys []*ProjectDaily
