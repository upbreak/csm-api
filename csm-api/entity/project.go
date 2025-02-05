package entity

import (
	"database/sql"
	"time"
)

type ProjectInfo struct {
	Sno              int64     `json:"sno"`
	Jno              int64     `json:"jno"`
	IsUse            string    `json:"is_use"`
	IsDefault        string    `json:"is_default"`
	RegDate          time.Time `json:"reg_date"`
	RegUser          string    `json:"reg_user"`
	RegUno           int64     `json:"reg_uno"`
	ModDate          time.Time `json:"mod_date"`
	ModUser          string    `json:"mod_user"`
	ModUno           string    `json:"mod_uno"`
	ProjectType      string    `json:"project_type"`
	ProjectNo        string    `json:"project_no"`
	ProjectNm        string    `json:"project_nm"`
	ProjectYear      int64     `json:"project_year"`
	ProjectLoc       string    `json:"project_loc"`
	ProjectCode      string    `json:"project_code"`
	ProjectCodeName  string    `json:"project_code_name"`
	SiteNm           string    `json:"site_nm"`
	CompCode         string    `json:"comp_code"`
	CompNick         string    `json:"comp_nick"`
	CompName         string    `json:"comp_name"`
	CompEtc          string    `json:"comp_etc"`
	OrderCompCode    string    `json:"order_comp_code"`
	OrderCompNick    string    `json:"order_comp_nick"`
	OrderCompName    string    `json:"order_comp_name"`
	OrderCompJobName string    `json:"order_comp_job_name"`
	ProjectLocName   string    `json:"project_loc_name"`
	JobPm            string    `json:"job_pm"`
	JobPe            string    `json:"job_pe"`
	ProjectStdt      time.Time `json:"project_stdt"`
	ProjectEddt      time.Time `json:"project_eddt"`
	ProjectRegDate   time.Time `json:"project_reg_date"`
	ProjectModDate   time.Time `json:"project_mod_date"`
	ProjectState     string    `json:"project_state"`
	MocNo            string    `json:"moc_no"`
	WoNo             string    `json:"wo_no"`

	ProjectPmList    *UserPmPeInfos `json:"project_pm_list"`
	DailyContentList *ProjectDailys `json:"daily_content_list"`
}

type ProjectInfos []*ProjectInfo

type ProjectInfoSql struct {
	Sno              sql.NullInt64  `db:"SNO"`
	Jno              sql.NullInt64  `db:"JNO"`
	IsUse            sql.NullString `db:"IS_USE"`
	IsDefault        sql.NullString `db:"IS_DEFAULT"`
	RegDate          sql.NullTime   `db:"REG_DATE"`
	RegUser          sql.NullString `db:"REG_USER"`
	RegUno           sql.NullInt64  `db:"REG_UNO"`
	ModDate          sql.NullTime   `db:"MOD_DATE"`
	ModUser          sql.NullString `db:"MOD_USER"`
	ModUno           sql.NullString `db:"MOD_UNO"`
	ProjectType      sql.NullString `db:"PROJECT_TYPE"`
	ProjectNo        sql.NullString `db:"PROJECT_NO"`
	ProjectNm        sql.NullString `db:"PROJECT_NM"`
	ProjectYear      sql.NullInt64  `db:"PROJECT_YEAR"`
	ProjectLoc       sql.NullString `db:"PROJECT_LOC"`
	ProjectCode      sql.NullString `db:"PROJECT_CODE"`
	ProjectCodeName  sql.NullString `db:"PROJECT_CODE_NAME"`
	SiteNm           sql.NullString `db:"SITE_NM"`
	CompCode         sql.NullString `db:"COMP_CODE"`
	CompNick         sql.NullString `db:"COMP_NICK"`
	CompName         sql.NullString `db:"COMP_NAME"`
	CompEtc          sql.NullString `db:"COMP_ETC"`
	OrderCompCode    sql.NullString `db:"ORDER_COMP_CODE"`
	OrderCompNick    sql.NullString `db:"ORDER_COMP_NICK"`
	OrderCompName    sql.NullString `db:"ORDER_COMP_NAME"`
	OrderCompJobName sql.NullString `db:"ORDER_COMP_JOB_NAME"`
	ProjectLocName   sql.NullString `db:"PROJECT_LOC_NAME"`
	JobPm            sql.NullString `db:"JOB_PM"`
	JobPe            sql.NullString `db:"JOB_PE"`
	ProjectStdt      sql.NullTime   `db:"PROJECT_STDT"`
	ProjectEddt      sql.NullTime   `db:"PROJECT_EDDT"`
	ProjectRegDate   sql.NullTime   `db:"PROJECT_REG_DATE"`
	ProjectModDate   sql.NullTime   `db:"PROJECT_MOD_DATE"`
	ProjectState     sql.NullString `db:"PROJECT_STATE"`
	MocNo            sql.NullString `db:"MOC_NO"`
	WoNo             sql.NullString `db:"WO_NO"`
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
	p.ProjectLocName = projectInfoSql.ProjectLoc.String
	p.JobPe = projectInfoSql.JobPe.String
	p.JobPe = projectInfoSql.JobPe.String
	p.ProjectStdt = projectInfoSql.ProjectStdt.Time
	p.ProjectEddt = projectInfoSql.ProjectEddt.Time
	p.ProjectRegDate = projectInfoSql.ProjectRegDate.Time
	p.ProjectModDate = projectInfoSql.ProjectModDate.Time
	p.ProjectState = projectInfoSql.ProjectState.String
	p.MocNo = projectInfoSql.MocNo.String
	p.WoNo = projectInfoSql.WoNo.String

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
