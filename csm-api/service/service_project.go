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

// func: 진행중 프로젝트 전체 조회
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

// func: 진행중 프로젝트 개수 조회
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
		return &entity.JobInfos{}, fmt.Errorf("seravice_project/GetStaffProjectList: %w", err)
	}

	jobInfos := &entity.JobInfos{}
	if err := entity.ConvertSliceToRegular(*jobInfoSqls, jobInfos); err != nil {
		return &entity.JobInfos{}, fmt.Errorf("seravice_project/:staff/ConvertSliceToReqular error %w", err)
	}

	return jobInfos, nil

}
