package entity

import "github.com/guregu/null"

type ProjectSetting struct {
	Jno         null.Int    `json:"jno" db:"JNO"`
	InTime      null.Time   `json:"in_time" db:"IN_TIME"`           // 출근시간
	OutTime     null.Time   `json:"out_time" db:"OUT_TIME"`         // 퇴근시간
	RespiteTime null.Int    `json:"respite_time" db:"RESPITE_TIME"` // 출/퇴근 유예시간(분)
	CancelCode  null.String `json:"cancel_code" db:"CANCEL_CODE"`   // 마감취소가능기한 CODE
	ManHours    *ManHours   `json:"man_hours"`
	Base
}

type ProjectSettings []*ProjectSetting

type ManHour struct {
	Mhno     null.Int    `json:"mhno" db:"MHNO"`
	WorkHour null.Int    `json:"work_hour" db:"WORK_HOUR"`
	ManHour  null.Float  `json:"man_hour" db:"MAN_HOUR"`
	Jno      null.Int    `json:"jno" db:"JNO"`
	Etc      null.String `json:"etc" db:"ETC"`
	Base
}
type ManHours []*ManHour
