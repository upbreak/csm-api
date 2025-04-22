package entity

import (
	"github.com/guregu/null"
)

type ProjectDaily struct {
	Idx        null.Int    `json:"idx" db:"IDX"`
	Jno        null.Int    `json:"jno" db:"JNO"`
	Content    null.String `json:"content" db:"CONTENT"`
	IsUse      null.String `json:"isUse" db:"IS_USE"`
	TargetDate null.Time   `json:"targetDate" db:"TARGET_DATE"`
	Base
}

type ProjectDailys []*ProjectDaily
