package entity

import (
	"github.com/guregu/null"
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

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Site struct {
	Sno                null.Int    `json:"sno" db:"SNO"`
	SiteNm             null.String `json:"site_nm" db:"SITE_NM"`
	Etc                null.String `json:"etc" db:"ETC"`
	LocCode            null.String `json:"loc_code" db:"LOC_CODE"`
	LocName            null.String `json:"loc_name" db:"LOC_NAME"`
	IsUse              null.String `json:"is_use" db:"IS_USE"`
	DefaultJno         null.Int    `json:"default_jno" db:"DEFAULT_JNO"`
	DefaultProjectName null.String `json:"default_project_name" db:"DEFAULT_PROJECT_NAME"`
	DefaultProjectNo   null.String `json:"default_project_no" db:"DEFAULT_PROJECT_NO"`
	CurrentSiteStats   null.String `json:"current_site_stats" db:"CURRENT_SITE_STATS"`
	Base

	ProjectList *ProjectInfos       `json:"project_list"`
	SitePos     *SitePos            `json:"site_pos"`
	SiteDate    *SiteDate           `json:"site_date"`
	Whether     WhetherSrtEntityRes `json:"whether"`
}

// struct: 현장 데이터 json용 구조체
//type SiteObj struct {
//	Site
//	ProjectList *ProjectInfos       `json:"project_list"`
//	SitePos     *SitePos            `json:"site_pos"`
//	SiteDate    *SiteDate           `json:"site_date"`
//	Whether     WhetherSrtEntityRes `json:"whether"`
//}

// struct: 현장 데이터 json 배열 구조체
type Sites []*Site

// struct: 현장 데이터 db용 구조체
//type SiteSql struct {
//	Sno                sql.NullInt64  `db:"SNO"`
//	SiteNm             sql.NullString `db:"SITE_NM"`
//	Etc                sql.NullString `db:"ETC"`
//	LocCode            sql.NullString `db:"LOC_CODE"`
//	LocName            sql.NullString `db:"LOC_NAME"`
//	IsUse              sql.NullString `db:"IS_USE"`
//	RegDate            sql.NullTime   `db:"REG_DATE"`
//	RegUser            sql.NullString `db:"REG_USER"`
//	RegUno             sql.NullInt64  `db:"REG_UNO"`
//	ModDate            sql.NullTime   `db:"MOD_DATE"`
//	ModUser            sql.NullString `db:"MOD_USER"`
//	ModUno             sql.NullInt64  `db:"MOD_UNO"`
//	DefaultJno         sql.NullInt64  `db:"DEFAULT_JNO"`
//	DefaultProjectName sql.NullString `db:"DEFAULT_PROJECT_NAME"`
//	DefaultProjectNo   sql.NullString `db:"DEFAULT_PROJECT_NO"`
//	CurrentSiteStats   sql.NullString `db:"CURRENT_SITE_STATS"`
//}
//
//// struct: 현장 데이터 db 배열 구조체
//type SiteSqls []*SiteSql
//
//// func: db -> json 구조체 변환
//// @param
//// - SiteSql: 현장 데이터 db 구조체
//func (s *Site) ToSite(siteSql *SiteSql) *Site {
//	s.Sno = siteSql.Sno.Int64
//	s.SiteNm = siteSql.SiteNm.String
//	s.Etc = siteSql.Etc.String
//	s.LocCode = siteSql.LocCode.String
//	s.LocName = siteSql.LocName.String
//	s.IsUse = siteSql.IsUse.String
//	s.RegDate = siteSql.RegDate.Time
//	s.RegUser = siteSql.RegUser.String
//	s.RegUno = siteSql.RegUno.Int64
//	s.ModDate = siteSql.ModDate.Time
//	s.ModUser = siteSql.ModUser.String
//	s.ModUno = siteSql.ModUno.Int64
//	s.DefaultJno = siteSql.DefaultJno.Int64
//	s.DefaultProjectName = siteSql.DefaultProjectName.String
//	s.DefaultProjectNo = siteSql.DefaultProjectNo.String
//	s.CurrentSiteStats = siteSql.CurrentSiteStats.String
//
//	return s
//}
//
//// func: db -> json 배열 구조체 변환
//// @param
//// - SiteSql: 현장 데이터 db 배열 구조체
//func (s *Sites) ToSites(siteSqls *SiteSqls) *Sites {
//	for _, siteSql := range *siteSqls {
//		site := Site{}
//		site.ToSite(siteSql)
//		*s = append(*s, &site)
//	}
//
//	return s
//}
