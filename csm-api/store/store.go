package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"github.com/guregu/null"
	"time"
)

type GetUserValidStore interface {
	GetUserValid(ctx context.Context, db Queryer, userId string, userPwd string) (entity.User, error)
}

type SiteStore interface {
	GetSiteList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.Sites, error)
	GetSiteNmList(ctx context.Context, db Queryer, page entity.PageSql, search entity.Site, nonSite int) (*entity.Sites, error)
	GetSiteNmCount(ctx context.Context, db Queryer, search entity.Site, nonSite int) (int, error)
	GetSiteStatsList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.Sites, error)
	ModifySite(ctx context.Context, tx Execer, site entity.Site) error
	AddSite(ctx context.Context, db Queryer, tx Execer, jno int64, user entity.User) error
	ModifySiteIsNonUse(ctx context.Context, tx Execer, site entity.ReqSite) error
}

type SitePosStore interface {
	GetSitePosList(ctx context.Context, db Queryer) ([]entity.SitePos, error)
	GetSitePosData(ctx context.Context, db Queryer, sno int64) (*entity.SitePos, error)
	ModifySitePosData(ctx context.Context, tx Execer, sno int64, sitePosSql entity.SitePos) error
	ModifySitePosIsNonUse(ctx context.Context, tx Execer, site entity.ReqSite) error
}

type SiteDateStore interface {
	GetSiteDateData(ctx context.Context, db Queryer, sno int64) (*entity.SiteDate, error)
	ModifySiteDate(ctx context.Context, tx Execer, sno int64, siteDateSql entity.SiteDate) error
	ModifySiteDateIsNonUse(ctx context.Context, tx Execer, site entity.ReqSite) error
}

type ProjectStore interface {
	GetProjectList(ctx context.Context, db Queryer, sno int64, targetDate time.Time) (*entity.ProjectInfos, error)
	GetProjectWorkerCountList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.ProjectInfos, error)
	GetProjectSafeWorkerCountList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.ProjectSafeCounts, error)
	GetProjectNmList(ctx context.Context, db Queryer) (*entity.ProjectInfos, error)
	GetUsedProjectList(ctx context.Context, db Queryer, pageSql entity.PageSql, search entity.JobInfo) (*entity.JobInfos, error)
	GetUsedProjectCount(ctx context.Context, db Queryer, search entity.JobInfo) (int, error)
	GetAllProjectList(ctx context.Context, db Queryer, pageSql entity.PageSql, search entity.JobInfo, isAll int) (*entity.JobInfos, error)
	GetAllProjectCount(ctx context.Context, db Queryer, search entity.JobInfo) (int, error)
	GetStaffProjectList(ctx context.Context, db Queryer, pageSql entity.PageSql, searchSql entity.JobInfo, uno sql.NullInt64) (*entity.JobInfos, error)
	GetStaffProjectCount(ctx context.Context, db Queryer, searchSql entity.JobInfo, uno sql.NullInt64) (int, error)
	GetProjectNmUnoList(ctx context.Context, db Queryer, uno sql.NullInt64, role int) (*entity.ProjectInfos, error)
	GetNonUsedProjectList(ctx context.Context, db Queryer, page entity.PageSql, search entity.NonUsedProject, retry string) (*entity.NonUsedProjects, error)
	GetNonUsedProjectCount(ctx context.Context, db Queryer, search entity.NonUsedProject, retry string) (int, error)
	AddProject(ctx context.Context, tx Execer, project entity.ReqProject) error
	ModifyDefaultProject(ctx context.Context, tx Execer, project entity.ReqProject) error
	ModifyUseProject(ctx context.Context, tx Execer, project entity.ReqProject) error
	RemoveProject(ctx context.Context, tx Execer, sno int64, jno int64) error
	ModifyProjectIsNonUse(ctx context.Context, tx Execer, site entity.ReqSite) error
	ModifyProject(ctx context.Context, tx Execer, project entity.ReqProject) error
}

type OrganizationStore interface {
	GetFuncNameList(ctx context.Context, db Queryer) (*entity.FuncNameSqls, error)
	GetOrganizationClientList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.OrganizationSqls, error)
	GetOrganizationHtencList(ctx context.Context, db Queryer, jno sql.NullInt64, funcNo sql.NullInt64) (*entity.OrganizationSqls, error)
}

type ProjectDailyStore interface {
	GetProjectDailyContentList(ctx context.Context, db Queryer, jno int64, targetDate time.Time) (*entity.ProjectDailys, error)
	GetDailyJobList(ctx context.Context, db Queryer, jno int64, targetDate string) (entity.ProjectDailys, error)
	AddDailyJob(ctx context.Context, tx Execer, project entity.ProjectDailys) error
	ModifyDailyJob(ctx context.Context, tx Execer, project entity.ProjectDaily) error
	RemoveDailyJob(ctx context.Context, tx Execer, idx int64) error
}

type UserStore interface {
	GetUserInfoPmPeList(ctx context.Context, db Queryer, unoList []int) (*entity.UserPmPeInfos, error)
}

type CodeStore interface {
	GetCodeList(ctx context.Context, db Queryer, pCode string) (*entity.Codes, error)
	GetCodeTree(ctx context.Context, db Queryer) (*entity.Codes, error)
	MergeCode(ctx context.Context, tx Execer, code entity.Code) error
	RemoveCode(ctx context.Context, tx Execer, idx int64) error
	ModifySortNo(ctx context.Context, tx Execer, codeSort entity.CodeSort) error
	DuplicateCheckCode(ctx context.Context, db Queryer, code string) (int, error)
}

type NoticeStore interface {
	GetNoticeList(ctx context.Context, db Queryer, uno null.Int, role int, pageSql entity.PageSql, search entity.Notice) (*entity.Notices, error)
	GetNoticeListCount(ctx context.Context, db Queryer, uno null.Int, role int, search entity.Notice) (int, error)
	AddNotice(ctx context.Context, tx Execer, notice entity.Notice) error
	ModifyNotice(ctx context.Context, tx Execer, notice entity.Notice) error
	RemoveNotice(ctx context.Context, tx Execer, idx null.Int) error
}

type DeviceStore interface {
	GetDeviceList(ctx context.Context, db Queryer, page entity.PageSql, search entity.Device, retry string) (*entity.Devices, error)
	GetDeviceListCount(ctx context.Context, db Queryer, search entity.Device, retry string) (int, error)
	AddDevice(ctx context.Context, tx Execer, device entity.Device) error
	ModifyDevice(ctx context.Context, tx Execer, device entity.Device) error
	RemoveDevice(ctx context.Context, tx Execer, dno sql.NullInt64) error
	GetDeviceLog(ctx context.Context, db Queryer) (*entity.RecdLogOrigins, error)
	GetCheckRegistered(ctx context.Context, db Queryer, deviceName string) (int, error)
}

type WorkerStore interface {
	GetWorkerTotalList(ctx context.Context, db Queryer, page entity.PageSql, search entity.Worker, retry string) (*entity.Workers, error)
	GetWorkerTotalCount(ctx context.Context, db Queryer, search entity.Worker, retry string) (int, error)
	GetAbsentWorkerList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerDaily, retry string) (*entity.Workers, error)
	GetAbsentWorkerCount(ctx context.Context, db Queryer, search entity.WorkerDaily, retry string) (int, error)
	AddWorker(ctx context.Context, tx Execer, worker entity.Worker) error
	ModifyWorker(ctx context.Context, tx Execer, worker entity.Worker) error
	GetWorkerSiteBaseList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerDaily, retry string) (*entity.WorkerDailys, error)
	GetWorkerSiteBaseCount(ctx context.Context, db Queryer, search entity.WorkerDaily, retry string) (int, error)
	MergeSiteBaseWorker(ctx context.Context, tx Execer, workers entity.WorkerDailys) error
	ModifyWorkerDeadline(ctx context.Context, tx Execer, workers entity.WorkerDailys) error
	ModifyWorkerProject(ctx context.Context, tx Execer, workers entity.WorkerDailys) error
	ModifyWorkerDeadlineInit(ctx context.Context, tx Execer) error
}

type CompanyStore interface {
	GetJobInfo(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.JobInfo, error)
	GetSiteManagerList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.Managers, error)
	GetSafeManagerList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.Managers, error)
	GetSupervisorList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.Supervisors, error)
	GetWorkInfoList(ctx context.Context, db Queryer) (*entity.WorkInfos, error)
	GetCompanyInfoList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.CompanyInfos, error)
	GetCompanyWorkInfoList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.WorkInfos, error)
}

type EquipStore interface {
	GetEquipList(ctx context.Context, db Queryer) (entity.EquipTemps, error)
	MergeEquipCnt(ctx context.Context, tx Execer, equips entity.EquipTemps) error
}

type ScheduleStore interface {
	GetRestScheduleList(ctx context.Context, db Queryer, jno int64, year string, month string) (entity.RestSchedules, error)
	AddRestSchedule(ctx context.Context, tx Execer, schedule entity.RestSchedules) error
	ModifyRestSchedule(ctx context.Context, tx Execer, schedule entity.RestSchedule) error
	RemoveRestSchedule(ctx context.Context, tx Execer, cno int64) error
}
