package service

import (
	"context"
	"csm-api/entity"
)

type GetUserValidService interface {
	GetUserValid(ctx context.Context, userId string, userPwd string) (entity.User, error)
}
