package entity

import "github.com/guregu/null"

type RestSchedule struct {
	Cno         null.Int    `json:"cno" db:"CNO"`
	Jno         null.Int    `json:"jno" db:"JNO"`
	IsEveryYear null.String `json:"is_every_year" db:"IS_EVERY_YEAR"`
	RestYear    null.Int    `json:"rest_year" db:"REST_YEAR"`
	RestMonth   null.Int    `json:"rest_month" db:"REST_MONTH"`
	RestDay     null.Int    `json:"rest_day" db:"REST_DAY"`
	Reason      null.String `json:"reason" db:"REASON"`
	Base
}

type RestSchedules []RestSchedule
