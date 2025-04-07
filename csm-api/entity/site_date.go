package entity

import (
	"github.com/guregu/null"
)

type SiteDate struct {
	OpeningDate         null.Time `json:"opening_date" db:"OPENING_DATE"`
	ClosingPlanDate     null.Time `json:"closing_plan_date" db:"CLOSING_PLAN_DATE"`
	ClosingForecastDate null.Time `json:"closing_forecast_date" db:"CLOSING_FORECAST_DATE"`
	ClosingActualDate   null.Time `json:"closing_actual_date" db:"CLOSING_ACTUAL_DATE"`
	Base
}

//type SiteDateSql struct {
//	OpeningDate         sql.NullTime   `db:"OPENING_DATE"`
//	ClosingPlanDate     sql.NullTime   `db:"CLOSING_PLAN_DATE"`
//	ClosingForecastDate sql.NullTime   `db:"CLOSING_FORECAST_DATE"`
//	ClosingActualDate   sql.NullTime   `db:"CLOSING_ACTUAL_DATE"`
//	RegUno              sql.NullInt64  `db:"REG_UNO"`
//	RegUser             sql.NullString `db:"REG_USER"`
//	RegDate             sql.NullTime   `db:"REG_DATE"`
//}
//
//func (s *SiteDate) ToSiteDate(sql *SiteDateSql) *SiteDate {
//	s.OpeningDate = sql.OpeningDate.Time
//	s.ClosingPlanDate = sql.ClosingPlanDate.Time
//	s.ClosingForecastDate = sql.ClosingForecastDate.Time
//	s.ClosingActualDate = sql.ClosingActualDate.Time
//	s.RegUno = sql.RegUno.Int64
//	s.RegUser = sql.RegUser.String
//	s.RegDate = sql.RegDate.Time
//
//	return s
//}
