package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
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
		projectInfo.ProjectPmList = userPmPeList
	}

	return projectInfos, nil
}

func (p *ServiceProject) GetProjectNmList(ctx context.Context) (*entity.ProjectInfos, error) {
	sqlList, err := p.Store.GetProjectNmList(ctx, p.DB)
	if err != nil {
		return &entity.ProjectInfos{}, fmt.Errorf("service_project/getProjectNmList error: %w", err)
	}
	projectInfos := &entity.ProjectInfos{}
	projectInfos.ToProjectInfos(sqlList)

	return projectInfos, nil
}
