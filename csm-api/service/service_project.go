package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type ServiceProject struct {
	DB          store.Queryer
	Store       store.ProjectStore
	UserService UserService
}

// 현장 고유번호로 현장에 해당하는 프로젝트 리스트 조회 비즈니스
//
// @param sno: 현장 고유번호
func (p *ServiceProject) GetProjectList(ctx context.Context, sno int64) (*entity.ProjectInfos, error) {
	projectInfoSqls, err := p.Store.GetProjectList(ctx, p.DB, sno)
	if err != nil {
		return &entity.ProjectInfos{}, fmt.Errorf("service_project/getProjectList error: %w", err)
	}
	projectInfos := &entity.ProjectInfos{}
	projectInfos.ToProjectInfos(projectInfoSqls)

	// 프로젝트 정보 객체에 pm, pe 정보 삽입
	for _, projectInfo := range *projectInfos {
		var unoList []int
		// pm uno 조회
		if &projectInfo.JobPm != nil && projectInfo.JobPm != "" {
			uno, err := strconv.Atoi(projectInfo.JobPm)
			if err != nil {
				return &entity.ProjectInfos{}, fmt.Errorf("service_project/strconv.Atoi(projectInfo.JobPm) parse err")
			}
			unoList = append(unoList, uno)
		}

		// pe uno 조회
		if &projectInfo.JobPe != nil && projectInfo.JobPe != "" {
			jobPeList := strings.Split(projectInfo.JobPe, ",")
			for _, jonPe := range jobPeList {
				uno, err := strconv.Atoi(jonPe)
				if err != nil {
					return &entity.ProjectInfos{}, fmt.Errorf("service_project/strconv.Atoi(jonPe) parse err")
				}
				unoList = append(unoList, uno)
			}
		}

		// pm, pe 정보 일괄 조회
		userPmPeList, err := p.UserService.GetUserInfoPmPeList(ctx, unoList)
		if err != nil {
			return &entity.ProjectInfos{}, fmt.Errorf("service_project/GetUserInfoPmPeList error: %w", err)
		}
		projectInfo.ProjectPeList = userPmPeList
	}

	return projectInfos, nil
}

// func: 프로젝트 조회(이름)
// @param
// -
func (p *ServiceProject) GetProjectNmList(ctx context.Context) (*entity.ProjectInfos, error) {
	sqlList, err := p.Store.GetProjectNmList(ctx, p.DB)
	if err != nil {
		return &entity.ProjectInfos{}, fmt.Errorf("service_project/getProjectNmList error: %w", err)
	}
	projectInfos := &entity.ProjectInfos{}
	projectInfos.ToProjectInfos(sqlList)

	return projectInfos, nil
}

// func: 프로젝트 전체 조회
// @param
// -
func (p *ServiceProject) GetUsedProjectList(ctx context.Context, page entity.Page, search entity.JobInfo) (*entity.JobInfos, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project/OfPageSql error: %w", err)
	}
	searchSql := &entity.JobInfoSql{}
	if err = entity.ConvertToSQLNulls(search, searchSql); err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project/ConvertToSQLNulls error: %w", err)
	}

	sqlList, err := p.Store.GetUsedProjectList(ctx, p.DB, pageSql, *searchSql)
	if err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project/GetUsedProjectList error: %w", err)
	}

	jobInfos := &entity.JobInfos{}
	if err = entity.ConvertSliceToRegular(*sqlList, jobInfos); err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project;used/ConvertSliceToRegular error: %w", err)
	}

	return jobInfos, nil
}

// func: 프로젝트 전체 조회 개수
// @param
// -
func (p *ServiceProject) GetUsedProjectCount(ctx context.Context, search entity.JobInfo) (int, error) {
	searchSql := &entity.JobInfoSql{}
	if err := entity.ConvertToSQLNulls(search, searchSql); err != nil {
		return 0, fmt.Errorf("service_project/ConvertToSQLNulls error: %w", err)
	}

	count, err := p.Store.GetUsedProjectCount(ctx, p.DB, *searchSql)
	if err != nil {
		return 0, fmt.Errorf("service_project/GetUsedProjectCount error: %w", err)
	}

	return count, nil
}

// func: 프로젝트 전체 조회
// @param
// -
func (p *ServiceProject) GetAllProjectList(ctx context.Context, page entity.Page, search entity.JobInfo) (*entity.JobInfos, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project/OfPageSql error: %w", err)
	}

	searchSql := &entity.JobInfoSql{}
	if err := entity.ConvertToSQLNulls(search, searchSql); err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project/ConvertToSQLNulls error: %w", err)
	}

	jobInfoSqls, err := p.Store.GetAllProjectList(ctx, p.DB, pageSql, *searchSql)
	if err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project/GetUsedProjectList error: %w", err)
	}

	jobInfos := &entity.JobInfos{}
	if err = entity.ConvertSliceToRegular(*jobInfoSqls, jobInfos); err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project;all/ConvertSliceToReqular error: %w", err)
	}

	return jobInfos, nil

}

// func: 프로젝트 개수 조회
// @param
// -
func (p *ServiceProject) GetAllProjectCount(ctx context.Context, search entity.JobInfo) (int, error) {
	searchSql := &entity.JobInfoSql{}
	if err := entity.ConvertToSQLNulls(search, searchSql); err != nil {
		return 0, fmt.Errorf("service_project/ConvertToSQLNulls error: %w", err)
	}

	count, err := p.Store.GetAllProjectCount(ctx, p.DB, *searchSql)
	if err != nil {
		return 0, fmt.Errorf("service_project/GetAllProjectCount error: %w", err)
	}

	return count, nil
}

// func: 조직도 확인
// @param
// - UNO
func (p *ServiceProject) GetStaffProjectList(ctx context.Context, page entity.Page, search entity.JobInfo, uno int64) (*entity.JobInfos, error) {
	var unoSql sql.NullInt64

	if uno != 0 {
		unoSql = sql.NullInt64{Valid: true, Int64: uno}
	} else {
		unoSql = sql.NullInt64{Valid: false}
	}

	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project/OfPageSql error: %w", err)
	}

	searchSql := &entity.JobInfoSql{}
	if err := entity.ConvertToSQLNulls(search, searchSql); err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project/ConvertToSQLNulls error: %w", err)

	}

	jobInfoSqls, err := p.Store.GetStaffProjectList(ctx, p.DB, pageSql, *searchSql, unoSql)
	if err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project/GetStaffProjectList: %w", err)
	}

	jobInfos := &entity.JobInfos{}
	if err := entity.ConvertSliceToRegular(*jobInfoSqls, jobInfos); err != nil {
		return &entity.JobInfos{}, fmt.Errorf("service_project;staff/ConvertSliceToReqular error %w", err)
	}

	return jobInfos, nil

}

// func: 조직도 확인 개수
// @param
// - UNO
func (p *ServiceProject) GetStaffProjectCount(ctx context.Context, search entity.JobInfo, uno int64) (int, error) {
	var unoSql sql.NullInt64

	if uno != 0 {
		unoSql = sql.NullInt64{Valid: true, Int64: uno}
	} else {
		unoSql = sql.NullInt64{Valid: false}
	}

	searchSql := &entity.JobInfoSql{}
	if err := entity.ConvertToSQLNulls(search, searchSql); err != nil {
		return 0, fmt.Errorf("service_project/ConvertToSQLNulls error: %w", err)

	}

	count, err := p.Store.GetStaffProjectCount(ctx, p.DB, *searchSql, unoSql)
	if err != nil {
		return 0, fmt.Errorf("service_project/GetStaffProjectCount error: %w", err)
	}

	return count, nil

}

// func: 조직도 조회: 고객사
// @param
// - JNO
func (p *ServiceProject) GetClientOrganization(ctx context.Context, jno int64) (*entity.Organizations, error) {
	var jnoSql sql.NullInt64

	if jno != 0 {
		jnoSql = sql.NullInt64{Valid: true, Int64: jno}
	} else {
		jnoSql = sql.NullInt64{Valid: false}
	}

	clientSql := &entity.OrganizationSqls{}
	clientSql, err := p.Store.GetClientOrganization(ctx, p.DB, jnoSql)
	if err != nil {
		return &entity.Organizations{}, fmt.Errorf("service_project/GetClientOrganization: %w", err)
	}

	client := &entity.Organizations{}
	if err := entity.ConvertSliceToRegular(*clientSql, client); err != nil {
		return &entity.Organizations{}, fmt.Errorf("service_project/CovertSliceToRegular: %w", err)
	}

	return client, nil
}

// func: 조직도 조회: 계약자(외부직원, 내부직원, 협력사)
// @param
// - JNO
func (p *ServiceProject) GetHitechOrganization(ctx context.Context, jno int64) (*entity.Organizations, error) {
	var jnoSql sql.NullInt64

	if jno != 0 {
		jnoSql = sql.NullInt64{Valid: true, Int64: jno}
	} else {
		jnoSql = sql.NullInt64{Valid: false}
	}

	hitechSql := &entity.OrganizationSqls{}
	hitechSql, err := p.Store.GetHitechOrganization(ctx, p.DB, jnoSql)
	if err != nil {
		return &entity.Organizations{}, fmt.Errorf("service_project/GetHitechOrganization: %w", err)
	}

	hitech := &entity.Organizations{}
	if err := entity.ConvertSliceToRegular(*hitechSql, hitech); err != nil {
		return &entity.Organizations{}, fmt.Errorf("service_project/ConvertSliceToRegular: %w", err)
	}

	return hitech, nil
}
