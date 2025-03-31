package entity

import (
	"database/sql"
	"time"
)

type ProjectInfo struct {
	Sno                   int64     `json:"sno"`
	Jno                   int64     `json:"jno"`
	IsUse                 string    `json:"is_use"`
	IsDefault             string    `json:"is_default"`
	RegDate               time.Time `json:"reg_date"`
	RegUser               string    `json:"reg_user"`
	RegUno                int64     `json:"reg_uno"`
	ModDate               time.Time `json:"mod_date"`
	ModUser               string    `json:"mod_user"`
	ModUno                string    `json:"mod_uno"`
	ProjectType           string    `json:"project_type"`
	ProjectTypeNm         string    `json:"project_type_nm"`
	ProjectNo             string    `json:"project_no"`
	ProjectNm             string    `json:"project_nm"`
	ProjectYear           int64     `json:"project_year"`
	ProjectLoc            string    `json:"project_loc"`
	ProjectCode           string    `json:"project_code"`
	ProjectCodeName       string    `json:"project_code_name"`
	SiteNm                string    `json:"site_nm"`
	CompCode              string    `json:"comp_code"`
	CompNick              string    `json:"comp_nick"`
	CompName              string    `json:"comp_name"`
	CompEtc               string    `json:"comp_etc"`
	OrderCompCode         string    `json:"order_comp_code"`
	OrderCompNick         string    `json:"order_comp_nick"`
	OrderCompName         string    `json:"order_comp_name"`
	OrderCompJobName      string    `json:"order_comp_job_name"`
	ProjectLocName        string    `json:"project_loc_name"`
	JobPm                 string    `json:"job_pm"`
	JobPmNm               string    `json:"job_pm_nm"`
	JobPe                 string    `json:"job_pe"`
	ProjectStdt           time.Time `json:"project_stdt"`
	ProjectEddt           time.Time `json:"project_eddt"`
	ProjectRegDate        time.Time `json:"project_reg_date"`
	ProjectModDate        time.Time `json:"project_mod_date"`
	ProjectState          string    `json:"project_state"`
	ProjectStateNm        string    `json:"project_state_nm"`
	MocNo                 string    `json:"moc_no"`
	WoNo                  string    `json:"wo_no"`
	WorkerCountAll        int64     `json:"worker_count_all"`
	WorkerCountDate       int64     `json:"worker_count_date"`
	WorkerCountHtenc      int64     `json:"worker_count_htenc"`
	WorkerCountWork       int64     `json:"worker_count_work"`
	WorkerCountSafe       int64     `json:"worker_count_safe"`
	WorkerCountManager    int64     `json:"worker_count_manager"`
	WorkerCountNotManager int64     `json:"worker_count_not_manager"`

	ProjectPeList    *UserPmPeInfos `json:"project_pe_list"`
	DailyContentList *ProjectDailys `json:"daily_content_list"`
}

type ProjectInfos []*ProjectInfo

type ProjectInfoSql struct {
	Sno                   sql.NullInt64  `db:"SNO"`
	Jno                   sql.NullInt64  `db:"JNO"`
	IsUse                 sql.NullString `db:"IS_USE"`
	IsDefault             sql.NullString `db:"IS_DEFAULT"`
	RegDate               sql.NullTime   `db:"REG_DATE"`
	RegUser               sql.NullString `db:"REG_USER"`
	RegUno                sql.NullInt64  `db:"REG_UNO"`
	ModDate               sql.NullTime   `db:"MOD_DATE"`
	ModUser               sql.NullString `db:"MOD_USER"`
	ModUno                sql.NullString `db:"MOD_UNO"`
	ProjectType           sql.NullString `db:"PROJECT_TYPE"`
	ProjectTypeNm         sql.NullString `db:"PROJECT_TYPE_NM"`
	ProjectNo             sql.NullString `db:"PROJECT_NO"`
	ProjectNm             sql.NullString `db:"PROJECT_NM"`
	ProjectYear           sql.NullInt64  `db:"PROJECT_YEAR"`
	ProjectLoc            sql.NullString `db:"PROJECT_LOC"`
	ProjectCode           sql.NullString `db:"PROJECT_CODE"`
	ProjectCodeName       sql.NullString `db:"PROJECT_CODE_NAME"`
	SiteNm                sql.NullString `db:"SITE_NM"`
	CompCode              sql.NullString `db:"COMP_CODE"`
	CompNick              sql.NullString `db:"COMP_NICK"`
	CompName              sql.NullString `db:"COMP_NAME"`
	CompEtc               sql.NullString `db:"COMP_ETC"`
	OrderCompCode         sql.NullString `db:"ORDER_COMP_CODE"`
	OrderCompNick         sql.NullString `db:"ORDER_COMP_NICK"`
	OrderCompName         sql.NullString `db:"ORDER_COMP_NAME"`
	OrderCompJobName      sql.NullString `db:"ORDER_COMP_JOB_NAME"`
	ProjectLocName        sql.NullString `db:"PROJECT_LOC_NAME"`
	JobPm                 sql.NullString `db:"JOB_PM"`
	JobPmNm               sql.NullString `db:"JOB_PM_NAME"`
	JobPe                 sql.NullString `db:"JOB_PE"`
	ProjectStdt           sql.NullTime   `db:"PROJECT_STDT"`
	ProjectEddt           sql.NullTime   `db:"PROJECT_EDDT"`
	ProjectRegDate        sql.NullTime   `db:"PROJECT_REG_DATE"`
	ProjectModDate        sql.NullTime   `db:"PROJECT_MOD_DATE"`
	ProjectState          sql.NullString `db:"PROJECT_STATE"`
	ProjectStateNm        sql.NullString `db:"PROJECT_STATE_NM"`
	MocNo                 sql.NullString `db:"MOC_NO"`
	WoNo                  sql.NullString `db:"WO_NO"`
	WorkerCountAll        sql.NullInt64  `db:"WORKER_COUNT_ALL"`
	WorkerCountDate       sql.NullInt64  `db:"WORKER_COUNT_DATE"`
	WorkerCountHtenc      sql.NullInt64  `db:"WORKER_COUNT_HTENC"`
	WorkerCountWork       sql.NullInt64  `db:"WORKER_COUNT_WORK"`
	WorkerCountSafe       sql.NullInt64  `db:"WORKER_COUNT_SAFE"`
	WorkerCountManager    sql.NullInt64  `db:"WORKER_COUNT_MANAGER"`
	WorkerCountNotManager sql.NullInt64  `db:"WORKER_COUNT_NOT_MANAGER"`
}

type ProjectInfoSqls []*ProjectInfoSql

func (p *ProjectInfo) ToProjectInfo(projectInfoSql *ProjectInfoSql) *ProjectInfo {
	p.Sno = projectInfoSql.Sno.Int64
	p.Jno = projectInfoSql.Jno.Int64
	p.IsUse = projectInfoSql.IsUse.String
	p.IsDefault = projectInfoSql.IsDefault.String
	p.RegDate = projectInfoSql.RegDate.Time
	p.RegUser = projectInfoSql.RegUser.String
	p.RegUno = projectInfoSql.RegUno.Int64
	p.ModDate = projectInfoSql.ModDate.Time
	p.ModUser = projectInfoSql.ModUser.String
	p.ModUno = projectInfoSql.ModUno.String
	p.ProjectType = projectInfoSql.ProjectType.String
	p.ProjectTypeNm = projectInfoSql.ProjectTypeNm.String
	p.ProjectNo = projectInfoSql.ProjectNo.String
	p.ProjectNm = projectInfoSql.ProjectNm.String
	p.ProjectYear = projectInfoSql.ProjectYear.Int64
	p.ProjectLoc = projectInfoSql.ProjectLoc.String
	p.ProjectCode = projectInfoSql.ProjectCode.String
	p.ProjectCodeName = projectInfoSql.ProjectCodeName.String
	p.SiteNm = projectInfoSql.SiteNm.String
	p.CompCode = projectInfoSql.CompCode.String
	p.CompNick = projectInfoSql.CompNick.String
	p.CompName = projectInfoSql.CompName.String
	p.CompEtc = projectInfoSql.CompEtc.String
	p.OrderCompCode = projectInfoSql.OrderCompCode.String
	p.OrderCompNick = projectInfoSql.OrderCompNick.String
	p.OrderCompName = projectInfoSql.OrderCompName.String
	p.OrderCompJobName = projectInfoSql.OrderCompJobName.String
	p.ProjectLocName = projectInfoSql.ProjectLocName.String
	p.JobPe = projectInfoSql.JobPe.String
	p.JobPmNm = projectInfoSql.JobPmNm.String
	p.JobPe = projectInfoSql.JobPe.String
	p.ProjectStdt = projectInfoSql.ProjectStdt.Time
	p.ProjectEddt = projectInfoSql.ProjectEddt.Time
	p.ProjectRegDate = projectInfoSql.ProjectRegDate.Time
	p.ProjectModDate = projectInfoSql.ProjectModDate.Time
	p.ProjectState = projectInfoSql.ProjectState.String
	p.ProjectStateNm = projectInfoSql.ProjectStateNm.String
	p.MocNo = projectInfoSql.MocNo.String
	p.WoNo = projectInfoSql.WoNo.String
	p.WorkerCountAll = projectInfoSql.WorkerCountAll.Int64
	p.WorkerCountDate = projectInfoSql.WorkerCountDate.Int64
	p.WorkerCountHtenc = projectInfoSql.WorkerCountHtenc.Int64
	p.WorkerCountWork = projectInfoSql.WorkerCountWork.Int64
	p.WorkerCountSafe = projectInfoSql.WorkerCountSafe.Int64
	p.WorkerCountManager = projectInfoSql.WorkerCountManager.Int64
	p.WorkerCountNotManager = projectInfoSql.WorkerCountNotManager.Int64

	return p
}

func (p *ProjectInfos) ToProjectInfos(projectInfoSqls *ProjectInfoSqls) *ProjectInfos {
	for _, projectInfoSql := range *projectInfoSqls {
		projectInfo := ProjectInfo{}
		projectInfo.ToProjectInfo(projectInfoSql)
		*p = append(*p, &projectInfo)
	}

	return p
}

type ProjectSafeCount struct {
	Sno       int64 `json:"sno"`
	Jno       int64 `json:"jno"`
	SafeCount int64 `json:"safe_count"`
}
type ProjectSafeCounts []*ProjectSafeCount

type ProjectSafeCountSql struct {
	Sno       sql.NullInt64 `db:"SNO"`
	Jno       sql.NullInt64 `db:"JNO"`
	SafeCount sql.NullInt64 `db:"SAFE_COUNT"`
}
type ProjectSafeCountSqls []*ProjectSafeCountSql

// nonUsedProject
type NonUsedProject struct {
	Rnum     int64  `json:"rnum"`
	Jno      int64  `json:"jno"`
	JobNo    string `json:"job_no"`
	JobName  string `json:"job_name"`
	JobYear  int64  `json:"job_year"`
	JobSd    string `json:"job_sd"`
	JobEd    string `json:"job_ed"`
	JobPmNm  string `json:"job_pm_nm"`
	DutyName string `json:"duty_name"`
}

type NonUsedProjects []*NonUsedProject

type NonUsedProjectSql struct {
	Rnum     sql.NullInt64  `db:"RNUM"`
	Jno      sql.NullInt64  `db:"JNO"`
	JobNo    sql.NullString `db:"JOB_NO"`
	JobName  sql.NullString `db:"JOB_NAME"`
	JobYear  sql.NullInt64  `db:"JOB_YEAR"`
	JobSd    sql.NullString `db:"JOB_SD"`
	JobEd    sql.NullString `db:"JOB_ED"`
	JobPmNm  sql.NullString `db:"JOB_PM_NM"`
	DutyName sql.NullString `db:"DUTY_NAME"`
}
type NonUsedProjectSqls []*NonUsedProjectSql
