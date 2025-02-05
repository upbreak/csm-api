package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
	"time"
)

type ServiceSite struct {
	DB                  store.Queryer
	Store               store.SiteStore
	ProjectService      ProjectService
	ProjectDailyService ProjectDailyService
	SitePosService      SitePosService
	SiteDateService     SiteDateService
}

// 현장 관리 리스트 조회
//
// @param targetDate: ???
func (s *ServiceSite) GetSiteList(ctx context.Context, targetDate time.Time) (*entity.Sites, error) {

	//현장관리 테이블 조회
	siteSqls, err := s.Store.GetSiteList(ctx, s.DB, targetDate)
	if err != nil {
		return &entity.Sites{}, fmt.Errorf("service_site/GetSiteList err: %w", err)
	}
	sites := &entity.Sites{}
	sites.ToSites(siteSqls)

	for _, site := range *sites {
		sno := site.Sno

		// 프로젝트 리스트 조회
		projectInfos, err := s.ProjectService.GetProjectList(ctx, sno)
		if err != nil {
			return &entity.Sites{}, fmt.Errorf("service_site/GetProjectList err: %w", err)
		}
		site.ProjectList = projectInfos

		for _, projectInfo := range *site.ProjectList {
			if &projectInfo.Jno != nil {
				projectDailyList, err := s.ProjectDailyService.GetProjectDailyContentList(ctx, projectInfo.Jno, targetDate)
				if err != nil {
					return &entity.Sites{}, fmt.Errorf("service_site/GetProjectDailyContentList err: %w", err)
				}
				projectInfo.DailyContentList = projectDailyList
			}
		}

		// 현장 위치 조회
		sitePosData, err := s.SitePosService.GetSitePosData(ctx, sno)
		if err != nil {
			return &entity.Sites{}, fmt.Errorf("service_site/GetSitePosData err: %w", err)
		}
		site.SitePos = sitePosData

		// 현장 날씨 조회
		siteDateData, err := s.SiteDateService.GetSiteDateData(ctx, sno)
		if err != nil {
			return &entity.Sites{}, fmt.Errorf("service_site/GetSiteDateData err: %w", err)
		}
		site.SiteDate = siteDateData
	}

	return sites, nil
}
