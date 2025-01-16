package entity

type User struct {
	UserId   string `json:"user_id" db:"USER_ID"`
	UserName string `json:"user_name" db:"USER_NAME"`
	UserPwd  string `json:"user_pwd" db:"USER_PWD"`
}
