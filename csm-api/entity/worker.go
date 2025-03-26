package entity

import (
	"database/sql"
	"time"
)

type Worker struct {
	RowNum     int64     `json:"rnum"`
	Sno        int64     `json:"sno"` //현장 고유번호
	Jno        int64     `json:"jno"` //프로젝트 고유번호
	JobName    string    `json:"job_name"`
	UserId     string    `json:"user_id"`    //근로자 아이디
	UserNm     string    `json:"user_nm"`    //근로자명
	Department string    `json:"department"` //부서or조직
	Phone      string    `json:"phone"`
	WorkerType string    `json:"worker_type"`
	IsUse      string    `json:"is_use"`
	IsRetire   string    `json:"is_retire"`
	RetireDate time.Time `json:"retire_date"`
	RegUser    string    `json:"reg_user"`
	RegDate    time.Time `json:"reg_date"`
	RegUno     int64     `json:"reg_uno"`
	ModUser    string    `json:"mod_user"`
	ModDate    time.Time `json:"mod_date"`
	ModUno     int64     `json:"mod_uno"`
	RecordDate string    `json:"record_date"`
}
type Workers []*Worker

type WorkerSql struct {
	RowNum     sql.NullInt64  `db:"RNUM"`
	Sno        sql.NullInt64  `db:"SNO"`
	Jno        sql.NullInt64  `db:"JNO"`
	JobName    sql.NullString `db:"JOB_NAME"`
	UserId     sql.NullString `db:"USER_ID"`
	UserNm     sql.NullString `db:"USER_NM"`
	Department sql.NullString `db:"DEPARTMENT"`
	Phone      sql.NullString `db:"PHONE"`
	WorkerType sql.NullString `db:"WORKER_TYPE"`
	IsUse      sql.NullString `db:"IS_USE"`
	IsRetire   sql.NullString `db:"IS_RETIRE"`
	RetireDate sql.NullTime   `db:"RETIRE_DATE"`
	RegUser    sql.NullString `db:"REG_USER"`
	RegDate    sql.NullTime   `db:"REG_DATE"`
	RegUno     sql.NullInt64  `db:"REG_UNO"`
	ModUser    sql.NullString `db:"MOD_USER"`
	ModDate    sql.NullTime   `db:"MOD_DATE"`
	ModUno     sql.NullInt64  `db:"MOD_UNO"`
	RecordDate sql.NullString `db:"RECORD_DATE"`
}
type WorkerSqls []*WorkerSql

type WorkerDaily struct {
	RowNum          int64     `json:"rnum"`
	Sno             int64     `json:"sno"`     //현장 고유번호
	Jno             int64     `json:"jno"`     //프로젝트 고유번호
	UserId          string    `json:"user_id"` //근로자 아이디
	UserNm          string    `json:"user_nm"`
	Department      string    `json:"department"`
	RecordDate      time.Time `json:"record_date"`
	InRecogTime     time.Time `json:"in_recog_time"`  //출근시간
	OutRecogTime    time.Time `json:"out_recog_time"` //퇴근시간
	Commute         string    `json:"commute"`
	IsDeadline      string    `json:"is_deadline"`
	SearchStartTime string    `json:"search_start_time"`
	SearchEndTime   string    `json:"search_end_time"`
	RegUser         string    `json:"reg_user"`
	RegDate         time.Time `json:"reg_date"`
	RegUno          int64     `json:"reg_uno"`
	ModUser         string    `json:"mod_user"`
	ModDate         time.Time `json:"mod_date"`
	ModUno          int64     `json:"mod_uno"`
}
type WorkerDailys []*WorkerDaily

type WorkerDailySql struct {
	RowNum          sql.NullInt64  `db:"RNUM"`
	Sno             sql.NullInt64  `db:"SNO"`
	Jno             sql.NullInt64  `db:"JNO"`
	UserId          sql.NullString `db:"USER_ID"`
	UserNm          sql.NullString `db:"USER_NM"`
	Department      sql.NullString `db:"DEPARTMENT"`
	RecordDate      sql.NullTime   `db:"RECORD_DATE"`
	InRecogTime     sql.NullTime   `db:"IN_RECOG_TIME"`
	OutRecogTime    sql.NullTime   `db:"OUT_RECOG_TIME"`
	Commute         sql.NullString `db:"COMMUTE"`
	IsDeadline      sql.NullString `db:"IS_DEADLINE"`
	SearchStartTime sql.NullString `db:"SEARCH_START_TIME"`
	SearchEndTime   sql.NullString `db:"SEARCH_END_TIME"`
	RegUser         sql.NullString `db:"REG_USER"`
	RegDate         sql.NullTime   `db:"REG_DATE"`
	RegUno          sql.NullInt64  `db:"REG_UNO"`
	ModUser         sql.NullString `db:"MOD_USER"`
	ModDate         sql.NullTime   `db:"MOD_DATE"`
	ModUno          sql.NullInt64  `db:"MOD_UNO"`
}
type WorkerDailySqls []*WorkerDailySql
