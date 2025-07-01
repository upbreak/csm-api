package service

import (
	"context"
	"crypto/md5"
	"csm-api/auth"
	"csm-api/entity"
	"csm-api/store"
	"encoding/hex"
	"fmt"
)

type UserValid struct {
	DB    store.Queryer
	Store store.GetUserValidStore
}

func (g *UserValid) GetUserValid(ctx context.Context, userId string, userPwd string) (entity.User, error) {
	// 비밀번호 암호화.
	hash := md5.Sum([]byte(userPwd))
	pwMd5 := hex.EncodeToString(hash[:])

	// 유저 db에서 확인
	user, err := g.Store.GetUserValid(ctx, g.DB, userId, pwMd5)
	if err != nil {
		return entity.User{}, fmt.Errorf("service.get user fail: %w", err)
	}

	// 권한
	if user.RoleCode == "" {
		if user.DeptName == "기술연구소" {
			user.RoleCode = string(auth.SystemAdmin)
		} else if user.TeamName == "프로젝트관리팀" {
			user.RoleCode = string(auth.SuperAdmin)
		} else {
			user.RoleCode = string(auth.User)
		}
	}

	return user, nil
}
