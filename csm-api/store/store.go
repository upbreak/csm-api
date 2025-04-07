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
	GetSiteList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.Sites, error)
	GetSiteNmList(ctx context.Context, db Queryer) (*entity.Sites, error)
	GetSiteStatsList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.Sites, error)
	ModifySite(ctx context.Context, db Beginner, site entity.Site) error
	AddSite(ctx context.Context, db Queryer, tdb Beginner, jno int64, user entity.User) error
}

type SitePosStore interface {
	GetSitePosData(ctx context.Context, db Queryer, sno int64) (*entity.SitePos, error)
	ModifySitePosData(ctx context.Context, db Beginner, sno int64, sitePosSql entity.SitePos) error
}

type SiteDateStore interface {
	GetSiteDateData(ctx context.Context, db Queryer, sno int64) (*entity.SiteDate, error)
	ModifySiteDate(ctx context.Context, db Beginner, sno int64, siteDateSql entity.SiteDate) error
}

type ProjectStore interface {
	GetProjectList(ctx context.Context, db Queryer, sno int64, targetDate time.Time) (*entity.ProjectInfos, error)
	GetProjectWorkerCountList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.ProjectInfos, error)
	GetProjectSafeWorkerCountList(ctx context.Context, db Queryer, targetDate time.Time) (*entity.ProjectSafeCounts, error)
	GetProjectNmList(ctx context.Context, db Queryer) (*entity.ProjectInfos, error)
	GetUsedProjectList(ctx context.Context, db Queryer, pageSql entity.PageSql, search entity.JobInfoSql) (*entity.JobInfoSqls, error)
	GetUsedProjectCount(ctx context.Context, db Queryer, search entity.JobInfoSql) (int, error)
	GetAllProjectList(ctx context.Context, db Queryer, pageSql entity.PageSql, search entity.JobInfoSql) (*entity.JobInfoSqls, error)
	GetAllProjectCount(ctx context.Context, db Queryer, search entity.JobInfoSql) (int, error)
	GetStaffProjectList(ctx context.Context, db Queryer, pageSql entity.PageSql, searchSql entity.JobInfoSql, uno sql.NullInt64) (*entity.JobInfoSqls, error)
	GetStaffProjectCount(ctx context.Context, db Queryer, searchSql entity.JobInfoSql, uno sql.NullInt64) (int, error)
	GetFuncNameList(ctx context.Context, db Queryer) (*entity.FuncNameSqls, error)
	GetClientOrganization(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.OrganizationSqls, error)
	GetHitechOrganization(ctx context.Context, db Queryer, jno sql.NullInt64, funcNo sql.NullInt64) (*entity.OrganizationSqls, error)
	GetProjectNmUnoList(ctx context.Context, db Queryer, uno sql.NullInt64, role int) (*entity.ProjectInfos, error)
	GetNonUsedProjectList(ctx context.Context, db Queryer, page entity.PageSql, search entity.NonUsedProject, retry string) (*entity.NonUsedProjects, error)
	GetNonUsedProjectCount(ctx context.Context, db Queryer, search entity.NonUsedProject, retry string) (int, error)
	AddProject(ctx context.Context, db Beginner, project entity.ReqProject) error
	ModifyDefaultProject(ctx context.Context, db Beginner, project entity.ReqProject) error
	ModifyUseProject(ctx context.Context, db Beginner, project entity.ReqProject) error
	RemoveProject(ctx context.Context, db Beginner, sno int64, jno int64) error
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
	GetNoticeList(ctx context.Context, db Queryer, uno sql.NullInt64, role int, pageSql entity.PageSql, search entity.NoticeSql) (*entity.NoticeSqls, error)
	GetNoticeListCount(ctx context.Context, db Queryer, uno sql.NullInt64, role int, search entity.NoticeSql) (int, error)
	AddNotice(ctx context.Context, db Beginner, noticeSql entity.NoticeSql) error
	ModifyNotice(ctx context.Context, db Beginner, noticeSql entity.NoticeSql) error
	RemoveNotice(ctx context.Context, db Beginner, idx entity.NoticeID) error
	GetNoticePeriod(ctx context.Context, db Queryer) (*entity.NoticePeriodSqls, error)
}

type DeviceStore interface {
	GetDeviceList(ctx context.Context, db Queryer, page entity.PageSql, search entity.DeviceSql, retry string) (*entity.DeviceSqls, error)
	GetDeviceListCount(ctx context.Context, db Queryer, search entity.DeviceSql, retry string) (int, error)
	AddDevice(ctx context.Context, db Beginner, device entity.DeviceSql) error
	ModifyDevice(ctx context.Context, db Beginner, device entity.DeviceSql) error
	RemoveDevice(ctx context.Context, db Beginner, dno sql.NullInt64) error
}

type WorkerStore interface {
	GetWorkerTotalList(ctx context.Context, db Queryer, page entity.PageSql, search entity.Worker, retry string) (*entity.Workers, error)
	GetWorkerTotalCount(ctx context.Context, db Queryer, search entity.Worker, retry string) (int, error)
	GetWorkerListByUserId(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerDaily, retry string) (*entity.Workers, error)
	GetWorkerCountByUserId(ctx context.Context, db Queryer, search entity.WorkerDaily, retry string) (int, error)
	AddWorker(ctx context.Context, db Beginner, worker entity.Worker) error
	ModifyWorker(ctx context.Context, db Beginner, worker entity.Worker) error
	GetWorkerSiteBaseList(ctx context.Context, db Queryer, page entity.PageSql, search entity.WorkerDaily, retry string) (*entity.WorkerDailys, error)
	GetWorkerSiteBaseCount(ctx context.Context, db Queryer, search entity.WorkerDaily, retry string) (int, error)
	MergeSiteBaseWorker(ctx context.Context, db Beginner, workers entity.WorkerDailys) error
	ModifyWorkerDeadline(ctx context.Context, db Beginner, workers entity.WorkerDailys) error
	ModifyWorkerProject(ctx context.Context, db Beginner, workers entity.WorkerDailys) error
}

type CompanyStore interface {
	GetJobInfo(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.JobInfoSql, error)
	GetSiteManagerList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.ManagerSqls, error)
	GetSafeManagerList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.ManagerSqls, error)
	GetSupervisorList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.SupervisorSqls, error)
	GetWorkInfoList(ctx context.Context, db Queryer) (*entity.WorkInfosqls, error)
	GetCompanyInfoList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.CompanyInfoSqls, error)
	GetCompanyWorkInfoList(ctx context.Context, db Queryer, jno sql.NullInt64) (*entity.WorkInfosqls, error)
}

type EquipStore interface {
	MergeEquipCnt(ctx context.Context, db Beginner, equips entity.EquipTemps) error
}
