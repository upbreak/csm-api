package entity

import (
	"database/sql"
	"time"
)

type SiteDate struct {
	OpeningDate         time.Time `json:"opening_date"`
	ClosingPlanDate     time.Time `json:"closing_plan_date"`
	ClosingForecastDate time.Time `json:"closing_forecast_date"`
	ClosingActualDate   time.Time `json:"closing_actual_date"`
	RegUno              int64     `json:"reg_uno"`
	RegUser             string    `json:"reg_user"`
	RegDate             time.Time `json:"reg_date"`
}

type SiteDateSql struct {
	OpeningDate         sql.NullTime   `db:"OPENING_DATE"`
	ClosingPlanDate     sql.NullTime   `db:"CLOSING_PLAN_DATE"`
	ClosingForecastDate sql.NullTime   `db:"CLOSING_FORECAST_DATE"`
	ClosingActualDate   sql.NullTime   `db:"CLOSING_ACTUAL_DATE"`
	RegUno              sql.NullInt64  `db:"REG_UNO"`
	RegUser             sql.NullString `db:"REG_USER"`
	RegDate             sql.NullTime   `db:"REG_DATE"`
}

func (s *SiteDate) ToSiteDate(sql *SiteDateSql) *SiteDate {
	s.OpeningDate = sql.OpeningDate.Time
	s.ClosingPlanDate = sql.ClosingPlanDate.Time
	s.ClosingForecastDate = sql.ClosingForecastDate.Time
	s.ClosingActualDate = sql.ClosingActualDate.Time
	s.RegUno = sql.RegUno.Int64
	s.RegUser = sql.RegUser.String
	s.RegDate = sql.RegDate.Time

	return s
}
