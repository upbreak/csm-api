package entity

import (
	"database/sql"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct: 현장 관리 응답 구조체
type SiteRes struct {
	Site Sites `json:"site"`
	Code Codes `json:"code"`
}

// struct: 현장 데이터 json용 구조체
type Site struct {
	Sno                int64     `json:"sno"`
	SiteNm             string    `json:"site_nm"`
	Etc                string    `json:"etc"`
	LocCode            string    `json:"loc_code"`
	LocName            string    `json:"loc_name"`
	IsUse              string    `json:"is_use"`
	RegDate            time.Time `json:"reg_date"`
	RegUser            string    `json:"reg_user"`
	RegUno             int64     `json:"reg_uno"`
	ModDate            time.Time `json:"mod_date"`
	ModUser            string    `json:"mod_user"`
	ModUno             int64     `json:"mod_uno"`
	DefaultJno         int64     `json:"default_jno"`
	DefaultProjectName string    `json:"default_project_name"`
	DefaultProjectNo   string    `json:"default_project_no"`
	CurrentSiteStats   string    `json:"current_site_stats"`

	ProjectList *ProjectInfos `json:"project_list"`
	SitePos     *SitePos      `json:"site_pos"`
	SiteDate    *SiteDate     `json:"site_date"`
}

// struct: 현장 데이터 json 배열 구조체
type Sites []*Site

// struct: 현장 데이터 db용 구조체
type SiteSql struct {
	Sno                sql.NullInt64  `db:"SNO"`
	SiteNm             sql.NullString `db:"SITE_NM"`
	Etc                sql.NullString `db:"ETC"`
	LocCode            sql.NullString `db:"LOC_CODE"`
	LocName            sql.NullString `db:"LOC_NAME"`
	IsUse              sql.NullString `db:"IS_USE"`
	RegDate            sql.NullTime   `db:"REG_DATE"`
	RegUser            sql.NullString `db:"REG_USER"`
	RegUno             sql.NullInt64  `db:"REG_UNO"`
	ModDate            sql.NullTime   `db:"MOD_DATE"`
	ModUser            sql.NullString `db:"MOD_USER"`
	ModUno             sql.NullInt64  `db:"MOD_UNO"`
	DefaultJno         sql.NullInt64  `db:"DEFAULT_JNO"`
	DefaultProjectName sql.NullString `db:"DEFAULT_PROJECT_NAME"`
	DefaultProjectNo   sql.NullString `db:"DEFAULT_PROJECT_NO"`
	CurrentSiteStats   sql.NullString `db:"CURRENT_SITE_STATS"`
}

// struct: 현장 데이터 db 배열 구조체
type SiteSqls []*SiteSql

// func: db -> json 구조체 변환
// @param
// - SiteSql: 현장 데이터 db 구조체
func (s *Site) ToSite(siteSql *SiteSql) *Site {
	s.Sno = siteSql.Sno.Int64
	s.SiteNm = siteSql.SiteNm.String
	s.Etc = siteSql.Etc.String
	s.LocCode = siteSql.LocCode.String
	s.LocName = siteSql.LocName.String
	s.IsUse = siteSql.IsUse.String
	s.RegDate = siteSql.RegDate.Time
	s.RegUser = siteSql.RegUser.String
	s.RegUno = siteSql.RegUno.Int64
	s.ModDate = siteSql.ModDate.Time
	s.ModUser = siteSql.ModUser.String
	s.ModUno = siteSql.ModUno.Int64
	s.DefaultJno = siteSql.DefaultJno.Int64
	s.DefaultProjectName = siteSql.DefaultProjectName.String
	s.DefaultProjectNo = siteSql.DefaultProjectNo.String
	s.CurrentSiteStats = siteSql.CurrentSiteStats.String

	return s
}

// func: db -> json 배열 구조체 변환
// @param
// - SiteSql: 현장 데이터 db 배열 구조체
func (s *Sites) ToSites(siteSqls *SiteSqls) *Sites {
	for _, siteSql := range *siteSqls {
		site := Site{}
		site.ToSite(siteSql)
		*s = append(*s, &site)
	}

	return s
}
