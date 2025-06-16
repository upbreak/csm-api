package entity

import "github.com/guregu/null"

type ManHour struct {
	Mhno     null.Int   `json:"mhno" db:"MHNO"`
	WorkHour null.Int   `json:"work_hour" db:"WORK_HOUR"`
	ManHour  null.Float `json:"man_hour" db:"MAN_HOUR"`
	Jno      null.Int   `json:"jno" db:"JNO"`
	Base
}
type ManHours []*ManHour
