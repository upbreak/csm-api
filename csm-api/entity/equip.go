package entity

import "github.com/guregu/null"

type EquipTemp struct {
	Sno        null.Int  `json:"sno" db:"SNO"`
	Jno        null.Int  `json:"jno" db:"JNO"`
	Cnt        null.Int  `json:"cnt" db:"CNT"`
	RecordDate null.Time `json:"record_date" db:"RECORD_DATE"`
	Base
}

type EquipTemps []*EquipTemp
