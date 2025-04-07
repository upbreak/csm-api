package entity

import (
	"github.com/guregu/null"
)

type Notice struct {
	RowNum           null.Int    `json:"row_num" db:"RNUM"`
	Idx              null.Int    `json:"idx" db:"IDX"`
	Sno              null.Int    `json:"sno" db:"SNO"`
	Jno              null.Int    `json:"jno" db:"JNO"`
	JobName          null.String `json:"job_name" db:"JOB_NAME"`
	JobLocName       null.String `json:"job_loc_name" db:"JOB_LOC_NAME"`
	Title            null.String `json:"title" db:"TITLE"`
	Content          null.String `json:"content" db:"CONTENT"`
	ShowYN           null.String `json:"show_yn" db:"SHOW_YN"`
	UserDutyName     null.String `json:"user_duty_name" db:"DUTY_NAME"`
	UserInfo         null.String `json:"user_info" db:"USER_INFO"`
	PeriodCode       null.String `json:"period_code" db:"PERIOD_CODE"`
	NoticeNm         null.String `json:"notice_nm" db:"NOTICE_NM"`
	PostingStartDate null.Time   `json:"posting_start_date" db:"POSTING_START_DATE"`
	PostingEndDate   null.Time   `json:"posting_end_date" db:"POSTING_END_DATE"`
	IsImportant      null.String `json:"is_important" db:"IS_IMPORTANT"`
	Base
}

type Notices []*Notice
