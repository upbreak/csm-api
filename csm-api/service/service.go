package service

import (
	"context"
	"csm-api/entity"
	"time"
)

type GetUserValidService interface {
	GetUserValid(ctx context.Context, userId string, userPwd string) (entity.User, error)
}

type SiteService interface {
	GetSiteList(ctx context.Context, targetDate time.Time) (*entity.Sites, error)
	GetSiteNmList(ctx context.Context) (*entity.Sites, error)
}

type SitePosService interface {
	GetSitePosData(ctx context.Context, sno int64) (*entity.SitePos, error)
}

type SiteDateService interface {
	GetSiteDateData(ctx context.Context, sno int64) (*entity.SiteDate, error)
}

type ProjectService interface {
	GetProjectList(ctx context.Context, sno int64) (*entity.ProjectInfos, error)
}

type ProjectDailyService interface {
	GetProjectDailyContentList(ctx context.Context, jno int64, targetDate time.Time) (*entity.ProjectDailys, error)
}

type UserService interface {
	GetUserInfoPmPeList(ctx context.Context, unoList []int) (*entity.UserPmPeInfos, error)
}

type CodeService interface {
	GetCodeList(ctx context.Context, pCode string) (*entity.Codes, error)
}

type NoticeService interface {
	GetNoticeList(ctx context.Context) ([]entity.Notice, error)
}

type DeviceService interface {
	GetDeviceList(ctx context.Context, page entity.Page) (*entity.Devices, error)
	GetDeviceListCount(ctx context.Context) (int, error)
}
