package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"time"
)

type GetUserValidStore interface {
	GetUserValid(ctx context.Context, db Queryer, userId string, userPwd string) (entity.User, error)
}

type SiteStore interface {
	GetSiteList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.SiteSqls, error)
	GetSiteNmList(ctx context.Context, db Queryer) (*entity.SiteSqls, error)
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

type CodeStore interface {
	GetCodeList(ctx context.Context, db Queryer, pCode string) (*entity.CodeSqls, error)
}

type NoticeStore interface {
	GetNoticeList(ctx context.Context, db Queryer, pageSql entity.PageSql) (*entity.NoticeSqls, error)
	GetNoticeListCount(ctx context.Context, db Queryer) (int, error)
	AddNotice(ctx context.Context, db Beginner, noticeSql entity.NoticeSql) error
	ModifyNotice(ctx context.Context, db Beginner, noticeSql entity.NoticeSql) error
	RemoveNotice(ctx context.Context, db Beginner, idx entity.NoticeID) error
}

type DeviceStore interface {
	GetDeviceList(ctx context.Context, db Queryer, page entity.PageSql) (*entity.DeviceSqls, error)
	GetDeviceListCount(ctx context.Context, db Queryer) (int, error)
	AddDevice(ctx context.Context, db Beginner, device entity.DeviceSql) error
	ModifyDevice(ctx context.Context, db Beginner, device entity.DeviceSql) error
	RemoveDevice(ctx context.Context, db Beginner, dno sql.NullInt64) error
}

type WorkerStore interface {
	GetWorkerTotalList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerSql) (*entity.WorkerSqls, error)
	GetWorkerTotalCount(ctx context.Context, db Queryer, searchTime string) (int, error)
}
