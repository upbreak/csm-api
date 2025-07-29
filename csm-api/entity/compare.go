package entity

import "github.com/guregu/null"

type Compare struct {
	Sno              null.Int    `json:"sno" db:"SNO"`
	Jno              null.Int    `json:"jno" db:"JNO"`
	UserKey          null.String `json:"user_key" db:"USER_KEY"`
	UserId           null.String `json:"user_id" db:"USER_ID"`
	UserNm           null.String `json:"user_nm" db:"USER_NM"`
	Department       null.String `json:"department" db:"DEPARTMENT"`
	DiscName         null.String `json:"disc_name" db:"DISC_NAME"`
	Phone            null.String `json:"phone" db:"PHONE"`
	Gender           null.String `json:"gender" db:"GENDER"`
	IsTbm            null.String `json:"is_tbm" db:"IS_TBM"`
	DeviceNm         null.String `json:"device_nm" db:"DEVICE_NM"`
	RecordDate       null.Time   `json:"record_date" db:"RECORD_DATE"`
	WorkerInTime     null.Time   `json:"worker_in_time" db:"WORKER_IN_TIME"`
	WorkerOutTime    null.Time   `json:"worker_out_time" db:"WORKER_OUT_TIME"`
	CompareState     null.String `json:"compare_state" db:"COMPARE_STATE"`
	IsDeadline       null.String `json:"is_deadline" db:"IS_DEADLINE"`
	DeductionInTime  null.Time   `json:"deduction_in_time" db:"DEDUCTION_IN_TIME"`
	DeductionOutTime null.Time   `json:"deduction_out_time" db:"DEDUCTION_OUT_TIME"`
	DeductionBirth   null.String `json:"deduction_birth" db:"DEDUCTION_BIRTH"`
}
