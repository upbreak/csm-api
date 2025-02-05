package store

import (
	"context"
	"csm-api/entity"
	"time"
)

type GetUserValidStore interface {
	GetUserValid(ctx context.Context, db Queryer, userId string, userPwd string) (entity.User, error)
}

type SiteStore interface {
	GetSiteList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.SiteSqls, error)
}

type SitePosStore interface {
	GetSitePosData(ctx context.Context, db Queryer, sno int64) (*entity.SitePosSql, error)
}

type SiteDateStore interface {
	GetSiteDateData(ctx context.Context, db Queryer, sno int64) (*entity.SiteDateSql, error)
}

type ProjectStore interface {
	GetProjectList(ctx context.Context, db Queryer, sno int64) (*entity.ProjectInfoSqls, error)
}

type ProjectDailyStore interface {
	GetProjectDailyContentList(ctx context.Context, db Queryer, jno int64, targetDate time.Time) (*entity.ProjectDailys, error)
}

type UserStore interface {
	GetUserInfoPmPeList(ctx context.Context, db Queryer, unoList []int) (*entity.UserPmPeInfoSqls, error)
}
