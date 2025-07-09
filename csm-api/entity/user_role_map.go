package entity

import "github.com/guregu/null"

type UserRoleMap struct {
	UserUno  null.Int    `json:"user_uno" db:"USER_UNO"`
	RoleCode null.String `json:"role_code" db:"ROLE_CODE"`
	Jno      null.Int    `json:"jno" db:"JNO"`
	Base
}
