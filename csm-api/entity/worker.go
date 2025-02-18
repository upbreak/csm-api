package entity

import (
	"database/sql"
	"time"
)

type Worker struct {
	RowNum       int64     `json:"row_num"`
	Dno          int64     `json:"dno"`     //홍채인식기 고유번호
	Sno          int64     `json:"sno"`     //현장 고유번호
	SiteNm       string    `json:"site_nm"` //현장 이름
	Jno          int64     `json:"jno"`     //프로젝트 고유번호
	JobName      string    `json:"job_name"`
	UserId       string    `json:"user_id"`        //근로자 아이디
	UserNm       string    `json:"user_nm"`        //근로자명
	Department   string    `json:"department"`     //부서or조직
	InRecogTime  time.Time `json:"in_recog_time"`  //출근시간
	OutRecogTime time.Time `json:"out_recog_time"` //퇴근시간
	SearchTime   string    `json:"search_time"`
}

type Workers []*Worker

type WorkerSql struct {
	RowNum       sql.NullInt64  `db:"RNUM"`
	Dno          sql.NullInt64  `db:"DNO"`
	Sno          sql.NullInt64  `db:"SNO"`
	SiteNm       sql.NullString `db:"SITE_NM"`
	Jno          sql.NullInt64  `db:"JNO"`
	JobName      sql.NullString `db:"JOB_NAME"`
	UserId       sql.NullString `db:"USER_ID"`
	UserNm       sql.NullString `db:"USER_NM"`
	Department   sql.NullString `db:"DEPARTMENT"`
	InRecogTime  sql.NullTime   `db:"IN_RECOG_TIME"`
	OutRecogTime sql.NullTime   `db:"OUT_RECOG_TIME"`
	SearchTime   sql.NullString `db:"SEARCH_TIME"`
}

type WorkerSqls []*WorkerSql
