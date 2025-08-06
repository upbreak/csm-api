package entity

import (
	"csm-api/utils"
	"github.com/guregu/null"
	"strings"
)

type User struct {
	Uno       int64  `json:"uno" db:"UNO"`
	UserId    string `json:"user_id" db:"USER_ID"`
	UserName  string `json:"user_name" db:"USER_NAME"`
	UserPwd   string `json:"user_pwd" db:"USER_PWD"`
	IsSaved   bool   `json:"is_saved"`
	Agent     string `json:"agent"`
	DeptName  string `json:"dept_name" db:"DEPT_NAME"`
	TeamName  string `json:"team_name" db:"TEAM_NAME"`
	RoleCode  string `json:"role_code" db:"ROLE_CODE"`
	IsCompany bool   `json:"is_company"`
	Admin     bool   `json:"admin"`
}

func (u User) SetUser(uno int64, userName string) User {
	u.Uno = uno
	u.UserName = userName
	u.Agent = utils.GetAgent()
	return u
}

type UserPeInfo struct {
	Uno    null.Int    `json:"uno" db:"UNO"`
	UserId null.String `json:"user_id" db:"USER_ID"`
	Name   null.String `json:"name" db:"USER_NAME"`
}

type UserPeInfos []*UserPeInfo

type Role struct {
	Api        null.String `db:"API"`
	PermitRole null.String `db:"PERMIT_ROLE"`
}

type RoleList []*Role

// func: 권한 테이블(IRIS_LIST_PERMIT_ROLE)에서 role이 있는지 확인 후 있으면 ture값 반환, 없으면 false값 반환
// @param
// - api: 요청 API에 따라 권한 분리
// - roles: 역할 문자열(|로 나열된 것)
func AuthorizationCheck(list RoleList, roles string) bool {

	check := false

	roleList := strings.Split(roles, "|")

	flag := false
	for _, permit := range list {
		for _, role := range roleList {
			if permit.PermitRole.String == role {
				check = true
				flag = true
				break
			}
		}
		if flag {
			break
		}
	}

	return check

}
