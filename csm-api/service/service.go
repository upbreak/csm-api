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
	ModifySite(ctx context.Context, site entity.Site) error
	AddSite(ctx context.Context, jno int64, user entity.User) error
}

type SitePosService interface {
	GetSitePosData(ctx context.Context, sno int64) (*entity.SitePos, error)
	ModifySitePos(ctx context.Context, sno int64, sitePos entity.SitePos) error
}

type SiteDateService interface {
	GetSiteDateData(ctx context.Context, sno int64) (*entity.SiteDate, error)
	ModifySiteDate(ctx context.Context, sno int64, siteDate entity.SiteDate) error
}

type ProjectService interface {
	GetProjectList(ctx context.Context, sno int64, targetDate time.Time) (*entity.ProjectInfos, error)
	GetProjectWorkerCountList(ctx context.Context, targetDate time.Time) (*entity.ProjectInfos, error)
	GetProjectNmList(ctx context.Context) (*entity.ProjectInfos, error)
	GetUsedProjectList(ctx context.Context, page entity.Page, search entity.JobInfo) (*entity.JobInfos, error)
	GetUsedProjectCount(ctx context.Context, search entity.JobInfo) (int, error)
	GetAllProjectList(ctx context.Context, page entity.Page, search entity.JobInfo) (*entity.JobInfos, error)
	GetAllProjectCount(ctx context.Context, search entity.JobInfo) (int, error)
	GetStaffProjectList(ctx context.Context, page entity.Page, search entity.JobInfo, uno int64) (*entity.JobInfos, error)
	GetStaffProjectCount(ctx context.Context, search entity.JobInfo, uno int64) (int, error)
	GetClientOrganization(ctx context.Context, jno int64) (*entity.OrganizationPartition, error)
	GetHitechOrganization(ctx context.Context, jno int64) (*entity.OrganizationPartitions, error)
	GetProjectNmUnoList(ctx context.Context, uno int64, role string) (*entity.ProjectInfos, error)
	GetNonUsedProjectList(ctx context.Context, page entity.Page, search entity.NonUsedProject, retry string) (*entity.NonUsedProjects, error)
	GetNonUsedProjectCount(ctx context.Context, search entity.NonUsedProject, retry string) (int, error)
	AddProject(ctx context.Context, project entity.ReqProject) error
	ModifyDefaultProject(ctx context.Context, project entity.ReqProject) error
	ModifyUseProject(ctx context.Context, project entity.ReqProject) error
	RemoveProject(ctx context.Context, sno int64, jno int64) error
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
	GetNoticeList(ctx context.Context, uno int64, role string, page entity.Page, search entity.Notice) (*entity.Notices, error)
	GetNoticeListCount(ctx context.Context, uno int64, role string, search entity.Notice) (int, error)
	AddNotice(ctx context.Context, notice entity.Notice) error
	ModifyNotice(ctx context.Context, notice entity.Notice) error
	RemoveNotice(ctx context.Context, idx int64) error
	GetNoticePeriod(ctx context.Context) (*entity.NoticePeriods, error)
}

type DeviceService interface {
	GetDeviceList(ctx context.Context, page entity.Page, search entity.Device, retry string) (*entity.Devices, error)
	GetDeviceListCount(ctx context.Context, search entity.Device, retry string) (int, error)
	AddDevice(ctx context.Context, device entity.Device) error
	ModifyDevice(ctx context.Context, device entity.Device) error
	RemoveDevice(ctx context.Context, dno int64) error
}

type WorkerService interface {
	GetWorkerTotalList(ctx context.Context, page entity.Page, search entity.Worker, retry string) (*entity.Workers, error)
	GetWorkerTotalCount(ctx context.Context, search entity.Worker, retry string) (int, error)
	GetWorkerListByUserId(ctx context.Context, page entity.Page, search entity.WorkerDaily, retry string) (*entity.Workers, error)
	GetWorkerCountByUserId(ctx context.Context, search entity.WorkerDaily, retry string) (int, error)
	AddWorker(ctx context.Context, worker entity.Worker) error
	ModifyWorker(ctx context.Context, worker entity.Worker) error
	GetWorkerSiteBaseList(ctx context.Context, page entity.Page, search entity.WorkerDaily, retry string) (*entity.WorkerDailys, error)
	GetWorkerSiteBaseCount(ctx context.Context, search entity.WorkerDaily, retry string) (int, error)
	MergeSiteBaseWorker(ctx context.Context, workers entity.WorkerDailys) error
	ModifyWorkerDeadline(ctx context.Context, workers entity.WorkerDailys) error
	ModifyWorkerProject(ctx context.Context, workers entity.WorkerDailys) error
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
	GetWhetherWrnMsg() (entity.WhetherWrnMsgList, error)
}
type AddressSearchAPIService interface {
	GetAPILatitudeLongtitude(roadAddress string) (*entity.Point, error)
	GetAPISiteMapPoint(roadAddress string) (*entity.MapPoint, error)
}
