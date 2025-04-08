package service

import (
	"context"
	"csm-api/api"
	"csm-api/entity"
	"csm-api/store"
	"database/sql"
	"encoding/json"
	"fmt"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-18
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

type ServiceCompany struct {
	SafeDB      store.Queryer
	TimeSheetDB store.Queryer
	Store       store.CompanyStore
}

// func: job 정보 조회
// @param
// - jno sql.NullInt64: 프로젝트 고유번호
func (s *ServiceCompany) GetJobInfo(ctx context.Context, jno int64) (*entity.JobInfo, error) {
	jnoSql := entity.ToSQLNulls(jno).(sql.NullInt64)
	data, err := s.Store.GetJobInfo(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_conpany/GetJobInfo err: %w", err)
	}
	return data, nil
}

// func: 현장소장 조회
// @param
// - jno int64: 프로젝트 고유번호
func (s *ServiceCompany) GetSiteManagerList(ctx context.Context, jno int64) (*entity.Managers, error) {
	jnoSql := entity.ToSQLNulls(jno).(sql.NullInt64)
	list, err := s.Store.GetSiteManagerList(ctx, s.TimeSheetDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetSiteManagerList err: %w", err)
	}

	return list, nil
}

// func: 안전관리자 조회
// @param
// - jno int64: 프로젝트 고유번호
func (s *ServiceCompany) GetSafeManagerList(ctx context.Context, jno int64) (*entity.Managers, error) {
	jnoSql := entity.ToSQLNulls(jno).(sql.NullInt64)
	list, err := s.Store.GetSafeManagerList(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetSafeManagerList err: %w", err)
	}

	return list, nil
}

// func: 관리감독자 조회
// @param
// - jno int64: 프로젝트 고유번호
func (s *ServiceCompany) GetSupervisorList(ctx context.Context, jno int64) (*entity.Supervisors, error) {
	jnoSql := entity.ToSQLNulls(jno).(sql.NullInt64)
	list, err := s.Store.GetSupervisorList(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetSupervisorList err: %w", err)
	}

	return list, nil
}

// func: 공종 정보 조회
// @param
func (s *ServiceCompany) GetWorkInfoList(ctx context.Context) (*entity.WorkInfos, error) {
	list, err := s.Store.GetWorkInfoList(ctx, s.SafeDB)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetWorkInfoList err: %w", err)
	}

	return list, nil
}

// func: 협력업체 정보 조회
// @param
// - jno int64: 프로젝트 고유번호
func (s *ServiceCompany) GetCompanyInfoList(ctx context.Context, jno int64) (*entity.CompanyInfoResList, error) {
	list := &entity.CompanyInfoResList{}
	jnoSql := entity.ToSQLNulls(jno).(sql.NullInt64)

	// 협력업체 정보 조회
	companyList, err := s.Store.GetCompanyInfoList(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetCompanyInfoList err: %w", err)
	}
	for _, item := range *companyList {
		temp := &entity.CompanyInfoRes{}
		temp.Jno = item.Jno.Int64
		temp.Cno = item.Cno.Int64
		temp.Id = item.Id.String
		temp.Cellphone = item.Cellphone.String
		temp.Email = item.Email.String
		*list = append(*list, temp)
	}

	// JOB별 협력업체 리스트 API
	url := fmt.Sprintf("http://wcfservice.hi-techeng.co.kr/apipcs/getcontractinfo?jno=%d&contracttype=C", jno)
	response, err := api.CallGetAPI(url)
	if err != nil {
		return nil, fmt.Errorf("service_conpany;companyInfo/call Get Api err: %w", err)
	}
	companyApiReq := &entity.CompanyApiReq{}
	if err = json.Unmarshal([]byte(response), companyApiReq); err != nil {
		return nil, fmt.Errorf("service_conpany;companyInfo/json.Unmarshal err: %w", err)
	}
	if companyApiReq.ResultType != "Success" {
		return nil, fmt.Errorf("service_conpany;companyInfo/Api ResultType not Success")
	}

	// 공종 조회
	workInfoList, err := s.Store.GetCompanyWorkInfoList(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetCompanyWorkInfoList err: %w", err)
	}

	companyApiValues := entity.CompanyApiValues{}
	duplicate := make(map[int64]bool)
	for _, company := range companyApiReq.Value {
		if !duplicate[int64(company.CompCno)] {
			duplicate[int64(company.CompCno)] = true
			temp := &entity.CompanyApiValue{}
			temp.Jno = company.Jno
			temp.CompCno = company.CompCno
			temp.CompNameKr = company.CompNameKr
			temp.WorkerName = company.WorkerName
			companyApiValues = append(companyApiValues, temp)
		}
	}

	matched := make(map[int]bool)
	for _, item := range *list {
		for idx, company := range companyApiValues {
			if item.Jno == int64(company.Jno) && item.Cno == int64(company.CompCno) {
				item.CompNameKr = company.CompNameKr
				item.WorkerName = company.WorkerName
				matched[idx] = true
				break
			}
		}

		for _, work := range *workInfoList {
			if item.Jno == work.Jno.Int64 && item.Cno == work.Cno.Int64 {
				item.WorkInfo = append(item.WorkInfo, work.FuncNo.Int64)
			}
		}
	}

	for idx, company := range companyApiValues {
		if !matched[idx] {
			temp := &entity.CompanyInfoRes{}
			temp.Jno = int64(company.Jno)
			temp.Cno = int64(company.CompCno)
			temp.CompNameKr = company.CompNameKr
			temp.WorkerName = company.WorkerName
			temp.WorkInfo = make([]int64, 0)
			*list = append(*list, temp)
		}
	}

	return list, nil
}
