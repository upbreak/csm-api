package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ServiceProject struct {
	SafeDB    store.Queryer
	SafeTDB   store.Beginner
	Store     store.ProjectStore
	UserStore store.UserStore
}

// 현장 고유번호로 현장에 해당하는 프로젝트 리스트 조회 비즈니스
//
// @param sno: 현장 고유번호, , targetDate time.Time: 현재시간
func (p *ServiceProject) GetProjectList(ctx context.Context, sno int64, targetDate time.Time) (*entity.ProjectInfos, error) {
	projectInfos, err := p.Store.GetProjectList(ctx, p.SafeDB, sno, targetDate)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.ProjectInfos{}, fmt.Errorf("service_project/getProjectList error: %w", err)
	}

	// 안전관리자 수 조회
	safeInfos, err := p.Store.GetProjectSafeWorkerCountList(ctx, p.SafeDB, targetDate)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_project/getProjectSafeWorkerCountList error: %w", err)
	}

	// 프로젝트 정보 객체에 pm, pe 정보 삽입
	for _, projectInfo := range *projectInfos {
		var unoList []int
		// pm uno 조회
		if &projectInfo.JobPm != nil && projectInfo.JobPm.String != "" {
			uno, err := strconv.Atoi(projectInfo.JobPm.String)
			if err != nil {
				//TODO: 에러 아카이브
				return &entity.ProjectInfos{}, fmt.Errorf("service_project/strconv.Atoi(projectInfo.JobPm) parse err")
			}
			unoList = append(unoList, uno)
		}

		// pe uno 조회
		if &projectInfo.JobPe != nil && projectInfo.JobPe.String != "" {
			jobPeList := strings.Split(projectInfo.JobPe.String, ",")
			for _, jonPe := range jobPeList {
				uno, err := strconv.Atoi(jonPe)
				if err != nil {
					//TODO: 에러 아카이브
					return &entity.ProjectInfos{}, fmt.Errorf("service_project/strconv.Atoi(jonPe) parse err")
				}
				unoList = append(unoList, uno)
			}
		}

		// pm, pe 정보 일괄 조회
		userPmPeList, err := p.UserStore.GetUserInfoPmPeList(ctx, p.SafeDB, unoList)
		if err != nil {
			//TODO: 에러 아카이브
			return &entity.ProjectInfos{}, fmt.Errorf("service_project/GetUserInfoPmPeList error: %w", err)
		}
		projectInfo.ProjectPeList = userPmPeList

		// 안전, 공사 근로자 수
		for _, safe := range *safeInfos {
			if projectInfo.Sno == safe.Sno && projectInfo.Jno == safe.Jno {
				projectInfo.WorkerCountSafe = safe.SafeCount
				projectInfo.WorkerCountWork.Int64 = projectInfo.WorkerCountHtenc.Int64 - safe.SafeCount.Int64
				projectInfo.WorkerCountWork.Valid = true
				break
			}
		}
	}

	return projectInfos, nil
}

// func: 프로젝트 근로자 수 조회
// @param
// - sno int64 현장 번호, targetDate time.Time: 현재시간
func (p *ServiceProject) GetProjectWorkerCountList(ctx context.Context, targetDate time.Time) (*entity.ProjectInfos, error) {
	// 근로자 수 조회
	projectInfos, err := p.Store.GetProjectWorkerCountList(ctx, p.SafeDB, targetDate)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.ProjectInfos{}, fmt.Errorf("service_project/getProjectWorkerCountList error: %w", err)
	}

	// 안전관리자 수 조회
	safeInfos, err := p.Store.GetProjectSafeWorkerCountList(ctx, p.SafeDB, targetDate)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_project/getProjectSafeWorkerCountList error: %w", err)
	}

	// 안전, 공사 근로자 수
	for _, project := range *projectInfos {
		for _, safe := range *safeInfos {
			if project.Sno == safe.Sno && project.Jno == safe.Jno {
				project.WorkerCountSafe = safe.SafeCount
				project.WorkerCountWork.Int64 = project.WorkerCountHtenc.Int64 - safe.SafeCount.Int64
				project.WorkerCountWork.Valid = true
				break
			}
		}
	}

	return projectInfos, nil
}

// func: 프로젝트 조회(이름)
// @param
// -
func (p *ServiceProject) GetProjectNmList(ctx context.Context) (*entity.ProjectInfos, error) {
	list, err := p.Store.GetProjectNmList(ctx, p.SafeDB)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.ProjectInfos{}, fmt.Errorf("service_project/getProjectNmList error: %w", err)
	}

	return list, nil
}

// func: 공사관리시스템 등록 프로젝트 전체 조회
// @param
// -
func (p *ServiceProject) GetUsedProjectList(ctx context.Context, page entity.Page, search entity.JobInfo, retry string) (*entity.JobInfos, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.JobInfos{}, fmt.Errorf("service_project/OfPageSql error: %w", err)
	}

	jobInfos, err := p.Store.GetUsedProjectList(ctx, p.SafeDB, pageSql, search, retry)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.JobInfos{}, fmt.Errorf("service_project/GetUsedProjectList error: %w", err)
	}

	return jobInfos, nil
}

// func: 프로젝트 전체 조회 개수
// @param
// -
func (p *ServiceProject) GetUsedProjectCount(ctx context.Context, search entity.JobInfo, retry string) (int, error) {
	count, err := p.Store.GetUsedProjectCount(ctx, p.SafeDB, search, retry)
	if err != nil {
		//TODO: 에러 아카이브
		return 0, fmt.Errorf("service_project/GetUsedProjectCount error: %w", err)
	}

	return count, nil
}

// func: 프로젝트 전체 조회
// @param
// -
func (p *ServiceProject) GetAllProjectList(ctx context.Context, page entity.Page, search entity.JobInfo, isAll int) (*entity.JobInfos, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.JobInfos{}, fmt.Errorf("service_project/OfPageSql error: %w", err)
	}

	jobInfos, err := p.Store.GetAllProjectList(ctx, p.SafeDB, pageSql, search, isAll)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.JobInfos{}, fmt.Errorf("service_project/GetUsedProjectList error: %w", err)
	}

	return jobInfos, nil

}

// func: 프로젝트 개수 조회
// @param
// -
func (p *ServiceProject) GetAllProjectCount(ctx context.Context, search entity.JobInfo, isAll int) (int, error) {
	count, err := p.Store.GetAllProjectCount(ctx, p.SafeDB, search)
	if isAll == 1 {
		count += 1
	}
	if err != nil {
		//TODO: 에러 아카이브
		return 0, fmt.Errorf("service_project/GetAllProjectCount error: %w", err)
	}

	return count, nil
}

// func: 본인이 속한 프로젝트 조회
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
		//TODO: 에러 아카이브
		return &entity.JobInfos{}, fmt.Errorf("service_project/OfPageSql error: %w", err)
	}

	jobInfos, err := p.Store.GetStaffProjectList(ctx, p.SafeDB, pageSql, search, unoSql)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.JobInfos{}, fmt.Errorf("service_project/GetStaffProjectList: %w", err)
	}
	return jobInfos, nil

}

// func: 본인이 속한 프로젝트 개수
// @param
// - UNO
func (p *ServiceProject) GetStaffProjectCount(ctx context.Context, search entity.JobInfo, uno int64) (int, error) {
	var unoSql sql.NullInt64

	if uno != 0 {
		unoSql = sql.NullInt64{Valid: true, Int64: uno}
	} else {
		unoSql = sql.NullInt64{Valid: false}
	}

	count, err := p.Store.GetStaffProjectCount(ctx, p.SafeDB, search, unoSql)
	if err != nil {
		//TODO: 에러 아카이브
		return 0, fmt.Errorf("service_project/GetStaffProjectCount error: %w", err)
	}

	return count, nil

}

// 본인이 속한 프로젝트 이름 목록
func (p *ServiceProject) GetProjectNmUnoList(ctx context.Context, uno int64, role string) (*entity.ProjectInfos, error) {

	var unoSql sql.NullInt64
	if uno != 0 {
		unoSql = sql.NullInt64{Valid: true, Int64: uno}
	} else {
		unoSql = sql.NullInt64{Valid: false}
	}

	var roleInt int
	if role == "ADMIN" {
		roleInt = 1
	} else {
		roleInt = 0
	}
	projectInfos, err := p.Store.GetProjectNmUnoList(ctx, p.SafeDB, unoSql, roleInt)

	if err != nil {
		//TODO: 에러 아카이브
		return &entity.ProjectInfos{}, fmt.Errorf("service_project/getProjectNmList error: %w", err)
	}

	return projectInfos, nil
}

// func: 현장근태 사용되지 않은 프로젝트
// @param
// -
func (s *ServiceProject) GetNonUsedProjectList(ctx context.Context, page entity.Page, search entity.NonUsedProject, retry string) (*entity.NonUsedProjects, error) {
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_project/GetNonUsedProjectList ofPageSql error: %w", err)
	}

	list, err := s.Store.GetNonUsedProjectList(ctx, s.SafeDB, pageSql, search, retry)

	if err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("service_project/GetNonUsedProjectList error: %w", err)
	}

	return list, nil
}

// func: 현장근태 사용되지 않은 프로젝트 수
// @param
// -
func (s *ServiceProject) GetNonUsedProjectCount(ctx context.Context, search entity.NonUsedProject, retry string) (int, error) {
	count, err := s.Store.GetNonUsedProjectCount(ctx, s.SafeDB, search, retry)

	if err != nil {
		//TODO: 에러 아카이브
		return 0, fmt.Errorf("service_project/GetNonUsedProjectCount error: %w", err)
	}

	return count, nil
}

// func: 현장 프로젝트 추가
// @param
// -
func (s *ServiceProject) AddProject(ctx context.Context, project entity.ReqProject) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_project/AddProject BeginTx error: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_project/AddProject Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_project/AddProject Commit error: %w", commitErr)
			}
		}
	}()

	err = s.Store.AddProject(ctx, tx, project)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_project/AddProject error: %w", err)
	}
	return
}

// func: 현장 기본 프로젝트 변경
// @param
// -
func (s *ServiceProject) ModifyDefaultProject(ctx context.Context, project entity.ReqProject) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_project/ModifyDefaultProject BeginTx error: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_project/ModifyDefaultProject Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_project/ModifyDefaultProject Commit error: %w", commitErr)
			}
		}
	}()

	err = s.Store.ModifyDefaultProject(ctx, tx, project)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_project/ModifyDefaultProject error: %w", err)
	}
	return
}

// func: 현장 프로젝트 사용여부 변경
// @param
// -
func (s *ServiceProject) ModifyUseProject(ctx context.Context, project entity.ReqProject) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_project/ModifyUseProject BeginTx error: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_project/ModifyUseProject Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_project/ModifyUseProject Commit error: %w", commitErr)
			}
		}
	}()

	err = s.Store.ModifyUseProject(ctx, tx, project)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_project/ModifyUseProject error: %w", err)
	}
	return
}

// func: 현장 프로젝트 삭제
// @param
// -
func (s *ServiceProject) RemoveProject(ctx context.Context, sno int64, jno int64) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_project/RemoveProject BeginTx error: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_project/RemoveProject Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_project/RemoveProject Commit error: %w", commitErr)
			}
		}
	}()

	err = s.Store.RemoveProject(ctx, tx, sno, jno)
	if err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("service_project/RemoveProject error: %w", err)
	}
	return
}
