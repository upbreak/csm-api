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
	GetSiteStatsList(ctx context.Context, targetDate time.Time) (*entity.Sites, error)
}

type SitePosService interface {
	GetSitePosData(ctx context.Context, sno int64) (*entity.SitePos, error)
}

type SiteDateService interface {
	GetSiteDateData(ctx context.Context, sno int64) (*entity.SiteDate, error)
}

type ProjectService interface {
	GetProjectList(ctx context.Context, sno int64) (*entity.ProjectInfos, error)
	GetProjectNmList(ctx context.Context) (*entity.ProjectInfos, error)
	GetUsedProjectList(ctx context.Context, page entity.Page, search entity.JobInfo) (*entity.JobInfos, error)
	GetUsedProjectCount(ctx context.Context, search entity.JobInfo) (int, error)
	GetAllProjectList(ctx context.Context, page entity.Page, search entity.JobInfo) (*entity.JobInfos, error)
	GetAllProjectCount(ctx context.Context, search entity.JobInfo) (int, error)
	GetStaffProjectList(ctx context.Context, page entity.Page, search entity.JobInfo, uno int64) (*entity.JobInfos, error)
	GetStaffProjectCount(ctx context.Context, search entity.JobInfo, uno int64) (int, error)
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
	GetNoticeList(ctx context.Context, page entity.Page, search entity.Notice) (*entity.Notices, error)
	GetNoticeListCount(ctx context.Context, search entity.Notice) (int, error)
	AddNotice(ctx context.Context, notice entity.Notice) error
	ModifyNotice(ctx context.Context, notice entity.Notice) error
	RemoveNotice(ctx context.Context, idx int64) error
}

type DeviceService interface {
	GetDeviceList(ctx context.Context, page entity.Page, search entity.Device) (*entity.Devices, error)
	GetDeviceListCount(ctx context.Context, search entity.Device) (int, error)
	AddDevice(ctx context.Context, device entity.Device) error
	ModifyDevice(ctx context.Context, device entity.Device) error
	RemoveDevice(ctx context.Context, dno int64) error
}

type WorkerService interface {
	GetWorkerTotalList(ctx context.Context, page entity.Page, search entity.Worker) (*entity.Workers, error)
	GetWorkerTotalCount(ctx context.Context, search entity.Worker) (int, error)
	GetWorkerSiteBaseList(ctx context.Context, page entity.Page, search entity.Worker) (*entity.Workers, error)
	GetWorkerSiteBaseCount(ctx context.Context, search entity.Worker) (int, error)
}

type CompanyService interface {
	GetJobInfo(ctx context.Context, jno int64) (*entity.JobInfo, error)
	GetSiteManagerList(ctx context.Context, jno int64) (*entity.Managers, error)
	GetSafeManagerList(ctx context.Context, jno int64) (*entity.Managers, error)
	GetSupervisorList(ctx context.Context, jno int64) (*entity.Supervisors, error)
	GetWorkInfoList(ctx context.Context) (*entity.WorkInfos, error)
	GetCompanyInfoList(ctx context.Context, jno int64) (*entity.CompanyInfoResList, error)
}

type WhetherApiService interface {
	GetWhetherSrtNcst(date string, time string, nx int, ny int) (entity.WhetherSrtEntityRes, error)
}
