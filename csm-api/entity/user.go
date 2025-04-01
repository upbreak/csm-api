package entity

import (
	"csm-api/utils"
	"database/sql"
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
	Uno    int64  `json:"uno"`
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

type UserPmPeInfos []*UserPmPeInfo

type UserPmPeInfoSql struct {
	Uno    sql.NullInt64  `db:"UNO"`
	UserId sql.NullString `db:"USER_ID"`
	Name   sql.NullString `db:"USER_NAME"`
}

type UserPmPeInfoSqls []*UserPmPeInfoSql

func (u *UserPmPeInfo) ToUserPmPeInfo(sql *UserPmPeInfoSql) *UserPmPeInfo {
	u.Uno = sql.Uno.Int64
	u.UserId = sql.UserId.String
	u.Name = sql.Name.String

	return u
}

func (u *UserPmPeInfos) ToUserPmPeInfos(sqls *UserPmPeInfoSqls) *UserPmPeInfos {
	for _, sql := range *sqls {
		info := UserPmPeInfo{}
		info.ToUserPmPeInfo(sql)
		*u = append(*u, &info)
	}

	return u
}
