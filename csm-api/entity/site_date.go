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
