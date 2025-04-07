package entity

import (
	"github.com/guregu/null"
)

type Worker struct {
	RowNum      null.Int    `json:"rnum" db:"RNUM"`
	Sno         null.Int    `json:"sno" db:"SNO"` //현장 고유번호
	Jno         null.Int    `json:"jno" db:"JNO"` //프로젝트 고유번호
	JobName     null.String `json:"job_name" db:"JOB_NAME"`
	UserId      null.String `json:"user_id" db:"USER_ID"` //근로자 아이디
	AfterUserId null.String `json:"after_user_id" db:"AFTER_USER_ID"`
	UserNm      null.String `json:"user_nm" db:"USER_NM"`       //근로자명
	Department  null.String `json:"department" db:"DEPARTMENT"` //부서or조직
	Phone       null.String `json:"phone" db:"PHONE"`
	WorkerType  null.String `json:"worker_type" db:"WORKER_TYPE"`
	IsUse       null.String `json:"is_use" db:"IS_USE"`
	IsRetire    null.String `json:"is_retire" db:"IS_RETIRE"`
	RetireDate  null.Time   `json:"retire_date" db:"RETIRE_DATE"`
	RecordDate  null.String `json:"record_date" db:"RECORD_DATE"`
	RegNo       null.String `json:"reg_no" db:"REG_NO"`
	Base
}
type Workers []*Worker

type WorkerDaily struct {
	RowNum          null.Int    `json:"rnum" db:"RNUM"`
	Sno             null.Int    `json:"sno" db:"SNO"`         //현장 고유번호
	Jno             null.Int    `json:"jno" db:"JNO"`         //프로젝트 고유번호
	UserId          null.String `json:"user_id" db:"USER_ID"` //근로자 아이디
	UserNm          null.String `json:"user_nm" db:"USER_NM"`
	Department      null.String `json:"department" db:"DEPARTMENT"`
	RecordDate      null.Time   `json:"record_date" db:"RECORD_DATE"`
	InRecogTime     null.Time   `json:"in_recog_time" db:"IN_RECOG_TIME"`   //출근시간
	OutRecogTime    null.Time   `json:"out_recog_time" db:"OUT_RECOG_TIME"` //퇴근시간
	WorkState       null.String `json:"work_state" db:"WORK_STATE"`
	IsDeadline      null.String `json:"is_deadline" db:"IS_DEADLINE"`
	SearchStartTime null.String `json:"search_start_time" db:"SEARCH_START_TIME"`
	SearchEndTime   null.String `json:"search_end_time" db:"SEARCH_END_TIME"`
	AfterJno        null.Int    `json:"after_jno" db:"AFTER_JNO"`
	Base
}
type WorkerDailys []*WorkerDaily
