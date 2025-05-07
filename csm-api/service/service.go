package service

import (
	"context"
	"csm-api/entity"
	"github.com/guregu/null"
	"github.com/xuri/excelize/v2"
	"time"
)

type GetUserValidService interface {
	GetUserValid(ctx context.Context, userId string, userPwd string) (entity.User, error)
}

type SiteService interface {
	GetSiteList(ctx context.Context, targetDate time.Time) (*entity.Sites, error)
	GetSiteNmList(ctx context.Context, page entity.Page, search entity.Site, nonSite int) (*entity.Sites, error)
	GetSiteNmCount(ctx context.Context, search entity.Site, nonSite int) (int, error)
	GetSiteStatsList(ctx context.Context, targetDate time.Time) (*entity.Sites, error)
	ModifySite(ctx context.Context, site entity.Site) error
	AddSite(ctx context.Context, jno int64, user entity.User) error
	ModifySiteIsNonUse(ctx context.Context, site entity.ReqSite) error
}

type SitePosService interface {
	GetSitePosList(ctx context.Context) ([]entity.SitePos, error)
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
	GetAllProjectList(ctx context.Context, page entity.Page, search entity.JobInfo, isAll int) (*entity.JobInfos, error)
	GetAllProjectCount(ctx context.Context, search entity.JobInfo, isAll int) (int, error)
	GetStaffProjectList(ctx context.Context, page entity.Page, search entity.JobInfo, uno int64) (*entity.JobInfos, error)
	GetStaffProjectCount(ctx context.Context, search entity.JobInfo, uno int64) (int, error)
	GetProjectNmUnoList(ctx context.Context, uno int64, role string) (*entity.ProjectInfos, error)
	GetNonUsedProjectList(ctx context.Context, page entity.Page, search entity.NonUsedProject, retry string) (*entity.NonUsedProjects, error)
	GetNonUsedProjectCount(ctx context.Context, search entity.NonUsedProject, retry string) (int, error)
	AddProject(ctx context.Context, project entity.ReqProject) error
	ModifyDefaultProject(ctx context.Context, project entity.ReqProject) error
	ModifyUseProject(ctx context.Context, project entity.ReqProject) error
	RemoveProject(ctx context.Context, sno int64, jno int64) error
}

type OrganizationService interface {
	GetOrganizationClientList(ctx context.Context, jno int64) (*entity.OrganizationPartitions, error)
	GetOrganizationHtencList(ctx context.Context, jno int64) (*entity.OrganizationPartitions, error)
}

type ProjectDailyService interface {
	GetDailyJobList(ctx context.Context, jno int64, targetDate string) (entity.ProjectDailys, error)
	AddDailyJob(ctx context.Context, project entity.ProjectDailys) error
	ModifyDailyJob(ctx context.Context, project entity.ProjectDaily) error
	RemoveDailyJob(ctx context.Context, idx int64) error
}

type UserService interface {
	GetUserInfoPmPeList(ctx context.Context, unoList []int) (*entity.UserPmPeInfos, error)
}

type CodeService interface {
	GetCodeList(ctx context.Context, pCode string) (*entity.Codes, error)
	GetCodeTree(ctx context.Context) (*entity.CodeTrees, error)
	MergeCode(ctx context.Context, code entity.Code) error
	RemoveCode(ctx context.Context, idx int64) error
	ModifySortNo(ctx context.Context, codeSorts entity.CodeSorts) error
	DuplicateCheckCode(ctx context.Context, code string) (bool, error)
}

type NoticeService interface {
	GetNoticeList(ctx context.Context, uno null.Int, role null.String, page entity.Page, search entity.Notice) (*entity.Notices, error)
	GetNoticeListCount(ctx context.Context, uno null.Int, role null.String, search entity.Notice) (int, error)
	AddNotice(ctx context.Context, notice entity.Notice) error
	ModifyNotice(ctx context.Context, notice entity.Notice) error
	RemoveNotice(ctx context.Context, idx null.Int) error
}

type DeviceService interface {
	GetDeviceList(ctx context.Context, page entity.Page, search entity.Device, retry string) (*entity.Devices, error)
	GetDeviceListCount(ctx context.Context, search entity.Device, retry string) (int, error)
	AddDevice(ctx context.Context, device entity.Device) error
	ModifyDevice(ctx context.Context, device entity.Device) error
	RemoveDevice(ctx context.Context, dno int64) error
	GetCheckRegisteredDevices(ctx context.Context) ([]string, error)
}

type WorkerService interface {
	GetWorkerTotalList(ctx context.Context, page entity.Page, search entity.Worker, retry string) (*entity.Workers, error)
	GetWorkerTotalCount(ctx context.Context, search entity.Worker, retry string) (int, error)
	GetAbsentWorkerList(ctx context.Context, page entity.Page, search entity.WorkerDaily, retry string) (*entity.Workers, error)
	GetAbsentWorkerCount(ctx context.Context, search entity.WorkerDaily, retry string) (int, error)
	AddWorker(ctx context.Context, worker entity.Worker) error
	ModifyWorker(ctx context.Context, worker entity.Worker) error
	GetWorkerSiteBaseList(ctx context.Context, page entity.Page, search entity.WorkerDaily, retry string) (*entity.WorkerDailys, error)
	GetWorkerSiteBaseCount(ctx context.Context, search entity.WorkerDaily, retry string) (int, error)
	MergeSiteBaseWorker(ctx context.Context, workers entity.WorkerDailys) error
	ModifyWorkerDeadline(ctx context.Context, workers entity.WorkerDailys) error
	ModifyWorkerProject(ctx context.Context, workers entity.WorkerDailys) error
	ModifyWorkerDeadlineInit(ctx context.Context) error
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

type RestDateApiService interface {
	GetRestDelDates(year string, month string) (entity.RestDates, error)
}

type EquipService interface {
	GetEquipList(ctx context.Context) (entity.EquipTemps, error)
	MergeEquipCnt(ctx context.Context, equips entity.EquipTemps) error
}

type ScheduleService interface {
	GetRestScheduleList(ctx context.Context, jno int64, year string, month string) (entity.RestSchedules, error)
	AddRestSchedule(ctx context.Context, schedule entity.RestSchedules) error
	ModifyRestSchedule(ctx context.Context, schedule entity.RestSchedule) error
	RemoveRestSchedule(ctx context.Context, cno int64) error
}

type ExcelService interface {
	ExportDailyDeduction(rows []entity.DailyDeduction) (*excelize.File, error)
	ImportDeduction(path string) error
}
