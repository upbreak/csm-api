package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct: 현장 데이터 서비스 구조체
type ServiceSite struct {
	DB                  store.Queryer
	Store               store.SiteStore
	ProjectService      ProjectService
	ProjectDailyService ProjectDailyService
	SitePosService      SitePosService
	SiteDateService     SiteDateService
	WhetherApiService   WhetherApiService
}

// func: 현장 관리 리스트 조회
// @param
// - targetDate: 현재시간
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

		// 현장 날짜 조회
		now := time.Now()
		baseDate := now.Format("20060102")
		baseTime := now.Format("1504")
		nx, ny := utils.LatLonToXY(sitePosData.Latitude, sitePosData.Longitude)
		siteWhether, err := s.WhetherApiService.GetWhetherSrtNcst(baseDate, baseTime, nx, ny)
		if err != nil {
			return &entity.Sites{}, fmt.Errorf("service_site/GetWhetherSrt err: %w", err)
		}
		site.Whether = siteWhether

		// 현장 날짜 조회
		siteDateData, err := s.SiteDateService.GetSiteDateData(ctx, sno)
		if err != nil {
			return &entity.Sites{}, fmt.Errorf("service_site/GetSiteDateData err: %w", err)
		}
		site.SiteDate = siteDateData
	}

	return sites, nil
}

// func: 현장 데이터 리스트 조회
// @param
// -
func (s *ServiceSite) GetSiteNmList(ctx context.Context) (*entity.Sites, error) {
	siteSqls, err := s.Store.GetSiteNmList(ctx, s.DB)
	if err != nil {
		return &entity.Sites{}, fmt.Errorf("service_site/GetSiteNmList err: %w", err)
	}
	sites := &entity.Sites{}
	sites.ToSites(siteSqls)

	return sites, nil
}

// func: 현장 상태 조회
// @param
// -
func (s *ServiceSite) GetSiteStatsList(ctx context.Context, targetDate time.Time) (*entity.Sites, error) {
	siteSqls, err := s.Store.GetSiteStatsList(ctx, s.DB, targetDate)
	if err != nil {
		return &entity.Sites{}, fmt.Errorf("service_site/GetSiteStatsList err: %w", err)
	}
	sites := &entity.Sites{}
	sites.ToSites(siteSqls)

	return sites, nil
}
