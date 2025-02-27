package entity

import (
	"database/sql"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-18
 * @modified 최종 수정일: 2025-02-26
 * @modifiedBy 최종 수정자: 정지영
 * @modified description
 * - 현장소장 및 안전관리자 userId, UserInfo 추가
 */

// struct: Begin:job 정보
type JobInfo struct {
	RowNum        int64  `json:"rnum"`
	Jno           int64  `json:"jno"`
	JobName       string `json:"job_name"`
	JobNo         string `json:"job_no"`
	JobSd         string `json:"job_sd"`
	JobEd         string `json:"job_ed"`
	CompName      string `json:"comp_name"`
	OrderCompName string `json:"order_comp_name"`
	JobPmName     string `json:"job_pm_name"`
	JobPmDutyName string `json:"job_pm_duty_name"`
	CdNm          string `json:"cd_nm"`
}
type JobInfos []*JobInfo
type JobInfoSql struct {
	RowNum        sql.NullInt64  `db:"RNUM"`
	Jno           sql.NullInt64  `db:"JNO"`
	JobName       sql.NullString `db:"JOB_NAME"`
	JobNo         sql.NullString `db:"JOB_NO"`
	JobSd         sql.NullString `db:"JOB_SD"`
	JobEd         sql.NullString `db:"JOB_ED"`
	CompName      sql.NullString `db:"COMP_NAME"`
	OrderCompName sql.NullString `db:"ORDER_COMP_NAME"`
	JobPmName     sql.NullString `db:"JOB_PM_NAME"`
	JobPmDutyName sql.NullString `db:"DUTY_NAME"`
	CdNm          sql.NullString `db:"CD_NM"`
}
type JobInfoSqls []*JobInfoSql

// struct: End:job 정보

// struct: Begin::현장소장|안전관리자
type Manager struct {
	Uno        int64  `json:"uno"`
	Jno        int64  `json:"jno"`
	UserName   string `json:"user_name"`
	DutyName   string `json:"duty_name"`
	UserId     string `json:"user_id"`
	TeamLeader string `json:"team_leader"`
	UserInfo   string `json:"user_info"`
}
type Managers []*Manager

type ManagerSql struct {
	Uno        sql.NullInt64  `db:"UNO"`
	Jno        sql.NullInt64  `db:"JNO"`
	UserName   sql.NullString `db:"USER_NAME"`
	DutyName   sql.NullString `db:"DUTY_NAME"`
	UserId     sql.NullString `db:"USER_ID"`
	TeamLeader sql.NullString `db:"TEAM_LEADER"`
	UserInfo   sql.NullString `db:"USER_INFO"`
}
type ManagerSqls []*ManagerSql

// struct: end::현장소장|안전관리자

// struct: Begin::관리감독자
type Supervisor struct {
	Uno       int64     `json:"uno"`
	Jno       int64     `json:"jno"`
	UserName  string    `json:"user_name"`
	UserId    string    `json:"user_id"`
	DutyName  string    `json:"duty_name"`
	DutyCd    string    `json:"duty_cd"`
	JobdutyId string    `json:"jobduty_id"`
	JoinDate  time.Time `json:"join_date"`
	FuncNo    string    `json:"func_no"`
}
type Supervisors []*Supervisor

type SupervisorSql struct {
	Uno       sql.NullInt64  `db:"UNO"`
	Jno       sql.NullInt64  `db:"JNO"`
	UserName  sql.NullString `db:"USER_NAME"`
	UserId    sql.NullString `db:"USER_ID"`
	DutyName  sql.NullString `db:"DUTY_NAME"`
	DutyCd    sql.NullString `db:"DUTY_CD"`
	JobdutyId sql.NullString `db:"JOBDUTY_ID"`
	JoinDate  sql.NullTime   `db:"JOIN_DATE"`
	FuncNo    sql.NullString `db:"FUNC_NO"`
}
type SupervisorSqls []*SupervisorSql

// struct: end::관리감독자

// struct: Begin::협력업체 정보
type CompanyInfo struct {
	Jno       int64  `json:"jno"`
	Cno       int64  `json:"cno"`
	Id        string `json:"id"`
	Pw        string `json:"pw"`
	Cellphone string `json:"cellphone"`
	Email     string `json:"email"`
	UserName  string `json:"username"`
	DutyName  string `json:"duty_name"`
}
type CompanyInfos []*CompanyInfo
type CompanyInfoSql struct {
	Jno       sql.NullInt64  `db:"JNO"`
	Cno       sql.NullInt64  `db:"CNO"`
	Id        sql.NullString `db:"ID"`
	Pw        sql.NullString `db:"PW"`
	Cellphone sql.NullString `db:"CELLPHONE"`
	Email     sql.NullString `db:"EMAIL"`
	UserName  sql.NullString `db:"USER_NAME"`
	DutyName  sql.NullString `db:"DUTY_NAME"`
}
type CompanyInfoSqls []*CompanyInfoSql

type CompanyInfoRes struct {
	Jno        int64   `json:"jno"`
	Cno        int64   `json:"cno"`
	CompNameKr string  `json:"comp_name_kr"`
	WorkerName string  `json:"worker_name"`
	Id         string  `json:"id"`
	Cellphone  string  `json:"cellphone"`
	Email      string  `json:"email"`
	WorkInfo   []int64 `json:"work_infos"`
}
type CompanyInfoResList []*CompanyInfoRes

// struct: End::협력업체 정보

// struct: Begin::공종 정보
type WorkInfo struct {
	Cno      int64  `json:"Cno"`
	Jno      int64  `json:"jno"`
	FuncNo   int64  `json:"func_no"`
	FuncName string `json:"func_name"`
}
type WorkInfos []*WorkInfo
type WorkInfosql struct {
	Cno      sql.NullInt64  `db:"CNO"`
	Jno      sql.NullInt64  `db:"JNO"`
	FuncNo   sql.NullInt64  `db:"FUNC_NO"`
	FuncName sql.NullString `db:"FUNC_NAME"`
}
type WorkInfosqls []*WorkInfosql

// struct: End::공종 정보

// struct: Begin::협력업체 리스트 api
type CompanyApiReq struct {
	ResultType string           `json:"ResultType"`
	ValueType  string           `json:"ValueType"`
	Value      CompanyApiValues `json:"Value"`
}

type CompanyApiValue struct {
	Cno          float64 `json:"CNO"`
	Jno          float64 `json:"JNO"`
	CompCno      float64 `json:"COMP_CNO"`
	ProjectName  string  `json:"PROJECT_NAME"`
	ItemDesc     string  `json:"ITEM_DESC"`
	ContractKind string  `json:"CONTRACT_KIND"`
	ContractDate string  `json:"CONTRACT_DATE"`
	DeliDate     string  `json:"DELI_DATE"`
	RegDate      string  `json:"REG_DATE"`
	ModDate      string  `json:"MOD_DATE"`
	CompNameKr   string  `json:"COMP_NAME_KR"`
	CompRegNo    string  `json:"COMP_REG_NO"`
	CompCeoName  string  `json:"COMP_CEO_NAME"`
	WorkerName   string  `json:"WORKER_NAME"`
}
type CompanyApiValues []*CompanyApiValue

// struct: End::협력업체 리스트 api
