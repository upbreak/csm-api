package entity

import (
	"csm-api/utils"
	"github.com/guregu/null"
)

type User struct {
	Uno      int64  `json:"uno" db:"UNO"`
	UserId   string `json:"user_id" db:"USER_ID"`
	UserName string `json:"user_name" db:"USER_NAME"`
	UserPwd  string `json:"user_pwd" db:"USER_PWD"`
	IsSaved  bool   `json:"is_saved"`
	Agent    string `json:"agent"`
}

func (u User) SetUser(uno int64, userName string) User {
	u.Uno = uno
	u.UserName = userName
	u.Agent = utils.GetAgent()
	return u
}

type UserPmPeInfo struct {
	Uno    null.Int    `json:"uno" db:"UNO"`
	UserId null.String `json:"user_id" db:"USER_ID"`
	Name   null.String `json:"name" db:"USER_NAME"`
}

type UserPmPeInfos []*UserPmPeInfo
