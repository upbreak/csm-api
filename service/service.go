package service

import "context"

type GetUserValidService interface {
	GetUserValid(ctx context.Context, userId string, userPwd string) (string, error)
}
