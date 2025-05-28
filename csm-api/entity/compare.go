package entity

import "github.com/guregu/null"

type Compare struct {
	Jno              null.Int    `json:"jno"`
	UserId           null.String `json:"user_id"`
	UserNm           null.String `json:"user_nm"`
	Department       null.String `json:"department"`
	DiscName         null.String `json:"disc_name"`
	IsTbm            null.String `json:"is_tbm"`
	RecordDate       null.Time   `json:"record_date"`
	WorkerInTime     null.Time   `json:"worker_in_time"`
	WorkerOutTime    null.Time   `json:"worker_out_time"`
	CompareState     null.String `json:"compare_state"`
	IsDeadline       null.String `json:"is_deadline"`
	DeductionInTime  null.Time   `json:"deduction_in_time"`
	DeductionOutTime null.Time   `json:"deduction_out_time"`
}
