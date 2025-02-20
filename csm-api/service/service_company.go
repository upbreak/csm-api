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
	sqlData, err := s.Store.GetJobInfo(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_conpany/GetJobInfo err: %w", err)
	}

	data := &entity.JobInfo{}
	if err = entity.ConvertToRegular(*sqlData, data); err != nil {
		return nil, fmt.Errorf("service_conpany;jonInfo/ConvertSliceToRegular err: %w", err)
	}
	return data, nil
}

// func: 현장소장 조회
// @param
// - jno int64: 프로젝트 고유번호
func (s *ServiceCompany) GetSiteManagerList(ctx context.Context, jno int64) (*entity.Managers, error) {
	jnoSql := entity.ToSQLNulls(jno).(sql.NullInt64)
	sqlList, err := s.Store.GetSiteManagerList(ctx, s.TimeSheetDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetSiteManagerList err: %w", err)
	}

	list := &entity.Managers{}
	if err = entity.ConvertSliceToRegular(*sqlList, list); err != nil {
		return nil, fmt.Errorf("service_company;siteManager/ConvertSliceToRegular err: %w", err)
	}

	return list, nil
}

// func: 안전관리자 조회
// @param
// - jno int64: 프로젝트 고유번호
func (s *ServiceCompany) GetSafeManagerList(ctx context.Context, jno int64) (*entity.Managers, error) {
	jnoSql := entity.ToSQLNulls(jno).(sql.NullInt64)
	sqlList, err := s.Store.GetSafeManagerList(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetSafeManagerList err: %w", err)
	}

	list := &entity.Managers{}
	if err = entity.ConvertSliceToRegular(*sqlList, list); err != nil {
		return nil, fmt.Errorf("service_company;safeManager/ConvertSliceToRegular err: %w", err)
	}

	return list, nil
}

// func: 관리감독자 조회
// @param
// - jno int64: 프로젝트 고유번호
func (s *ServiceCompany) GetSupervisorList(ctx context.Context, jno int64) (*entity.Supervisors, error) {
	jnoSql := entity.ToSQLNulls(jno).(sql.NullInt64)
	sqlList, err := s.Store.GetSupervisorList(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetSupervisorList err: %w", err)
	}

	list := &entity.Supervisors{}
	if err = entity.ConvertSliceToRegular(*sqlList, list); err != nil {
		return nil, fmt.Errorf("service_company;supervisor/ConvertSliceToRegular err: %w", err)
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
	sqlCompanyList, err := s.Store.GetCompanyInfoList(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetCompanyInfoList err: %w", err)
	}
	companyList := &entity.CompanyInfos{}
	if err = entity.ConvertSliceToRegular(*sqlCompanyList, companyList); err != nil {
		return nil, fmt.Errorf("service_company;companyInfo/ConvertSliceToRegular err: %w", err)
	}
	for _, item := range *companyList {
		temp := &entity.CompanyInfoRes{}
		temp.Jno = item.Jno
		temp.Cno = item.Cno
		temp.Id = item.Id
		temp.Cellphone = item.Cellphone
		temp.Email = item.Email
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
	sqlWorkInfoList, err := s.Store.GetCompanyWorkInfoList(ctx, s.SafeDB, jnoSql)
	if err != nil {
		return nil, fmt.Errorf("service_company/GetCompanyWorkInfoList err: %w", err)
	}
	workInfoList := &entity.WorkInfos{}
	if err = entity.ConvertSliceToRegular(*sqlWorkInfoList, workInfoList); err != nil {
		return nil, fmt.Errorf("service_company;companyWorkInfo/ConvertSliceToRegular err: %w", err)
	}

	for _, item := range *list {
		for _, company := range companyApiReq.Value {
			if item.Jno == int64(company.Jno) && item.Cno == int64(company.CompCno) {
				item.CompNameKr = company.CompNameKr
				item.CompCeoName = company.CompCeoName
			}
		}
		for _, work := range *workInfoList {
			if item.Jno == work.Jno && item.Cno == work.Cno {
				item.WorkInfo = append(item.WorkInfo, work.FuncNo)
			}
		}
	}

	return list, nil
}
