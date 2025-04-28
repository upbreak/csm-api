package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"strings"
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
	SafeDB                  store.Queryer
	SafeTDB                 store.Beginner
	Store                   store.SiteStore
	ProjectStore            store.ProjectStore
	ProjectDailyStore       store.ProjectDailyStore
	SitePosStore            store.SitePosStore
	SiteDateStore           store.SiteDateStore
	ProjectService          ProjectService
	WhetherApiService       WhetherApiService
	AddressSearchAPIService AddressSearchAPIService
}

// func: 현장 관리 리스트 조회
// @param
// - targetDate: 현재시간
func (s *ServiceSite) GetSiteList(ctx context.Context, targetDate time.Time) (*entity.Sites, error) {

	//현장관리 테이블 조회
	sites, err := s.Store.GetSiteList(ctx, s.SafeDB, targetDate)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.Sites{}, fmt.Errorf("service_site/GetSiteList err: %w", err)
	}

	for _, site := range *sites {
		sno := site.Sno.Int64

		// 프로젝트 리스트 조회
		projectInfos, err := s.ProjectService.GetProjectList(ctx, sno, targetDate)
		if err != nil {
			//TODO: 에러 아카이브
			return &entity.Sites{}, fmt.Errorf("service_site/GetProjectList err: %w", err)
		}
		site.ProjectList = projectInfos

		for _, projectInfo := range *site.ProjectList {
			if &projectInfo.Jno != nil {
				projectDailyList, err := s.ProjectDailyStore.GetProjectDailyContentList(ctx, s.SafeDB, projectInfo.Jno.Int64, targetDate)
				if err != nil {
					//TODO: 에러 아카이브
					return &entity.Sites{}, fmt.Errorf("service_site/GetProjectDailyContentList err: %w", err)
				}
				projectInfo.DailyContentList = projectDailyList
			}
		}

		// 현장 위치 조회
		sitePos, err := s.SitePosStore.GetSitePosData(ctx, s.SafeDB, sno)
		if err != nil {
			//TODO: 에러 아카이브
			return &entity.Sites{}, fmt.Errorf("service_site/GetSitePosData err: %w", err)
		}
		if sitePos.RoadAddress.String == "" {
			depthArray := []string{sitePos.RoadAddressNameDepth1.String, sitePos.RoadAddressNameDepth2.String, sitePos.RoadAddressNameDepth3.String, sitePos.RoadAddressNameDepth4.String, sitePos.RoadAddressNameDepth5.String}
			roadAddress := ""
			for _, depth := range depthArray {
				if depth != "" {
					roadAddress = roadAddress + " " + depth
				}
			}
			sitePos.RoadAddress = utils.ParseNullString(strings.Trim(roadAddress, " "))
		}
		site.SitePos = sitePos

		// 현장 날씨 조회
		//now := time.Now()
		//baseDate := now.Format("20060102")
		//baseTime := now.Add(time.Minute * -30).Format("1504") // 기상청에서 30분 단위로 발표하기 때문에 30분 전의 데이터 요청
		//nx, ny := utils.LatLonToXY(sitePos.Latitude.Float64, sitePos.Longitude.Float64)
		//
		//siteWhether, err := s.WhetherApiService.GetWhetherSrtNcst(baseDate, baseTime, nx, ny)
		//if err != nil {
		//	//TODO: 에러 아카이브
		//	return &entity.Sites{}, fmt.Errorf("service_site/GetWhetherSrt err: %w", err)
		//}
		//site.Whether = siteWhether

		// 현장 날짜 조회
		siteDateData, err := s.SiteDateStore.GetSiteDateData(ctx, s.SafeDB, sno)
		if err != nil {
			//TODO: 에러 아카이브
			return &entity.Sites{}, fmt.Errorf("service_site/GetSiteDateData err: %w", err)
		}
		site.SiteDate = siteDateData
	}

	return sites, nil
}

// func: 현장 데이터 리스트 조회
// @param
// -
func (s *ServiceSite) GetSiteNmList(ctx context.Context, page entity.Page, search entity.Site, nonSite int) (*entity.Sites, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_site/GetSiteNmList OfPageSql err : %w", err)
	}

	sites, err := s.Store.GetSiteNmList(ctx, s.SafeDB, pageSql, search, nonSite)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.Sites{}, fmt.Errorf("service_site/GetSiteNmList err: %w", err)
	}

	return sites, nil
}

// func: 현장 데이터 리스트 개수 조회
// @param
// -
func (s *ServiceSite) GetSiteNmCount(ctx context.Context, search entity.Site, nonSite int) (int, error) {

	count, err := s.Store.GetSiteNmCount(ctx, s.SafeDB, search, nonSite)
	if err != nil {
		//TODO: 에러 아카이브
		return 0, fmt.Errorf("service_site/GetSiteNmList err: %w", err)
	}

	return count, nil
}

// func: 현장 상태 조회
// @param
// -
func (s *ServiceSite) GetSiteStatsList(ctx context.Context, targetDate time.Time) (*entity.Sites, error) {
	sites, err := s.Store.GetSiteStatsList(ctx, s.SafeDB, targetDate)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.Sites{}, fmt.Errorf("service_site/GetSiteStatsList err: %w", err)
	}

	return sites, nil
}

// func: 현장 수정
// @param
// -
func (s *ServiceSite) ModifySite(ctx context.Context, site entity.Site) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_site/ModifySite err: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_site/ModifySite rollback err: %v", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_site/ModifySite commit err: %v", commitErr)
			}
		}
	}()

	if site.Sno.Int64 == 0 {
		//TODO: 에러 아카이브
		return fmt.Errorf("sno parameter is missing")
	}
	// 비고 정보 수정
	if err = s.Store.ModifySite(ctx, tx, site); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_site/ModifySite err: %w", err)
	}

	// 기본 프로젝트 수정
	project := entity.ReqProject{
		Jno: site.DefaultJno,
		Sno: site.Sno,
		Base: entity.Base{
			ModUno:  site.ModUno,
			ModUser: site.ModUser,
		},
	}
	if err = s.ProjectStore.ModifyDefaultProject(ctx, tx, project); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_site/ModifyDefaultProject err: %w", err)
	}

	// 날짜 수정 정보가 있는 경우만 실행
	if site.SiteDate != nil {
		siteDate := *site.SiteDate
		if err = s.SiteDateStore.ModifySiteDate(ctx, tx, site.Sno.Int64, siteDate); err != nil {
			//TODO: 에러 아카이브
			return fmt.Errorf("service_site/ModifySiteDate err: %v\n", err)
		}
	}

	// 장소 수정 할 정보가 있는 경우만 실행
	if site.SitePos != nil && site.SitePos.RoadAddress.String != "" {
		sitePos := *site.SitePos
		point, err := s.AddressSearchAPIService.GetAPILatitudeLongtitude(site.SitePos.RoadAddress.String)
		if err != nil {
			//TODO: 에러 아카이브
			return fmt.Errorf("service_site/GetApiLatitudeLongtitude err: %w", err)
		}
		if point.Latitude != 0.0 {
			sitePos.Latitude.Float64 = point.Latitude
			sitePos.Latitude.Valid = true
		}
		if point.Longitude != 0.0 {
			sitePos.Longitude.Float64 = point.Longitude
			sitePos.Longitude.Valid = true
		}

		if err = s.SitePosStore.ModifySitePosData(ctx, tx, site.Sno.Int64, sitePos); err != nil {
			//TODO: 에러 아카이브
			return fmt.Errorf("service_site/ModifySitePos err: %v\n", err)
		}
	}

	return

}

// func: 현장 생성
// @param
// -
func (s *ServiceSite) AddSite(ctx context.Context, jno int64, user entity.User) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_site/AddSite err: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_site/AddSite rollback err: %v", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_site/AddSite commit err: %v", commitErr)
			}
		}
	}()

	err = s.Store.AddSite(ctx, s.SafeDB, tx, jno, user)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_site/AddSite err: %w", err)
	}

	return nil
}

// func: 현장 사용안함 변경
// @param
// -
func (s *ServiceSite) ModifySiteIsNonUse(ctx context.Context, sno int64) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_site/ModifySiteIsNonUse err: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_site/ModifySiteIsNonUse rollback err: %v", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_site/ModifySiteIsNonUse commit err: %v", commitErr)
			}
		}
	}()

	// 현장
	if err = s.Store.ModifySiteIsNonUse(ctx, tx, sno); err != nil {
		return fmt.Errorf("service_site/ModifySiteIsNonUse err: %w", err)
	}

	// 프로젝트
	if err = s.ProjectStore.ModifyProjectIsNonUse(ctx, tx, sno); err != nil {
		return fmt.Errorf("service_site/ModifySiteIsNonUse err: %w", err)
	}

	// 위치
	if err = s.SitePosStore.ModifySitePosIsNonUse(ctx, tx, sno); err != nil {
		return fmt.Errorf("service_site/ModifySiteIsNonUse err: %w", err)
	}

	// 날짜
	if err = s.SiteDateStore.ModifySiteDateIsNonUse(ctx, tx, sno); err != nil {
		return fmt.Errorf("service_site/ModifySiteDateIsNonUse err: %w", err)
	}

	return
}
