package entity

import "github.com/guregu/null"

type ProjectSetting struct {
	Jno           null.Int    `json:"jno" db:"JNO"`
	InTime        null.Time   `json:"in_time" db:"IN_TIME"`           // 출근시간
	OutTime       null.Time   `json:"out_time" db:"OUT_TIME"`         // 퇴근시간
	RespiteTime   null.Int    `json:"respite_time" db:"RESPITE_TIME"` // 출/퇴근 유예시간(분)
	CancelCode    null.String `json:"cancel_code" db:"CANCEL_CODE"`   // 마감취소가능기한 CODE
	CancelDay     null.Int    `json:"cancel_day" db:"CANCEL_DAY"`
	ManHours      *ManHours   `json:"man_hours"`                          // 공수 정보
	Message       null.String `json:"message" db:"MESSAGE"`               // 로그 message
	ChangeSetting null.String `json:"change_setting" db:"CHANGE_SETTING"` // 변경된 테이블
	Base
}

type ProjectSettings []*ProjectSetting

type ManHour struct {
	Mhno          null.Int    `json:"mhno" db:"MHNO"`
	WorkHour      null.Int    `json:"work_hour" db:"WORK_HOUR"`
	ManHour       null.Float  `json:"man_hour" db:"MAN_HOUR"`
	Jno           null.Int    `json:"jno" db:"JNO"`
	Etc           null.String `json:"etc" db:"ETC"`
	Message       null.String `json:"message" db:"MESSAGE"`               // 로그 message
	ChangeSetting null.String `json:"change_setting" db:"CHANGE_SETTING"` // 변경된 테이블
	Base
}

type ManHours []*ManHour
