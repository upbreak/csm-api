package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"database/sql"
	"fmt"
)

type ServiceOrganization struct {
	TimeSheetDB store.Queryer
	Store       store.OrganizationStore
}

// func: 조직도 조회: 고객사
// @param
// - JNO
func (s *ServiceOrganization) GetOrganizationClientList(ctx context.Context, jno int64) (*entity.OrganizationPartitions, error) {
	var jnoSql sql.NullInt64

	if jno != 0 {
		jnoSql = sql.NullInt64{Valid: true, Int64: jno}
	} else {
		jnoSql = sql.NullInt64{Valid: false}
	}

	clientSql := &entity.OrganizationSqls{}
	clientSql, err := s.Store.GetOrganizationClientList(ctx, s.TimeSheetDB, jnoSql)
	if err != nil {
		return &entity.OrganizationPartitions{}, fmt.Errorf("ServiceOrganization/GetOrganizationClientList: %w", err)
	}

	clients := &entity.Organizations{} //  []organization
	if err := entity.ConvertSliceToRegular(*clientSql, clients); err != nil {
		return &entity.OrganizationPartitions{}, fmt.Errorf("ServiceOrganization/CovertSliceToRegular: %w", err)
	}

	organizations := entity.OrganizationPartitions{}
	if len(*clients) != 0 {
		// 공종 별로 구분하여 데이터 반환
		funcName := (*clients)[0].FuncName
		funcClients := &entity.Organizations{}
		for _, client := range *clients {
			if funcName != client.FuncName {
				organization := &entity.OrganizationPartition{}
				organization.FuncName = funcName
				organization.OrganizationList = funcClients
				organizations = append(organizations, organization)
				funcName = client.FuncName
				funcClients = &entity.Organizations{}
			}

			*funcClients = append(*funcClients, client)
		}

		organization := &entity.OrganizationPartition{}
		organization.FuncName = funcName
		organization.OrganizationList = funcClients
		organizations = append(organizations, organization)
	}

	return &organizations, nil
}

// func: 조직도 조회: 계약자(외부직원, 내부직원, 협력사)
// @param
// - JNO
func (s *ServiceOrganization) GetOrganizationHtencList(ctx context.Context, jno int64) (*entity.OrganizationPartitions, error) {
	var jnoSql sql.NullInt64

	if jno != 0 {
		jnoSql = sql.NullInt64{Valid: true, Int64: jno}
	} else {
		jnoSql = sql.NullInt64{Valid: false}
	}

	funcNameSqls, err := s.Store.GetFuncNameList(ctx, s.TimeSheetDB)
	if err != nil {
		return &entity.OrganizationPartitions{}, fmt.Errorf("ServiceOrganization/GetFuncNameList: %w", err)
	}

	funcNames := &entity.FuncNames{}
	if err := entity.ConvertSliceToRegular(*funcNameSqls, funcNames); err != nil {
		return &entity.OrganizationPartitions{}, fmt.Errorf("ServiceOrganization/CovertSliceToRegular: %w", err)
	}

	organizations := entity.OrganizationPartitions{}
	for _, funcName := range *funcNames {
		var funcNoSql sql.NullInt64
		if funcName.FuncNo != 0 {
			funcNoSql = sql.NullInt64{Valid: true, Int64: funcName.FuncNo}
		} else {
			funcNoSql = sql.NullInt64{Valid: false}
		}

		organization := entity.OrganizationPartition{}
		hitechSql := &entity.OrganizationSqls{}
		hitechSql, err := s.Store.GetOrganizationHtencList(ctx, s.TimeSheetDB, jnoSql, funcNoSql)
		if err != nil {
			return &entity.OrganizationPartitions{}, fmt.Errorf("ServiceOrganization/GetOrganizationHtencList: %w", err)
		}
		if len(*hitechSql) == 0 {
			continue
		}

		hitech := &entity.Organizations{}
		if err := entity.ConvertSliceToRegular(*hitechSql, hitech); err != nil {
			return &entity.OrganizationPartitions{}, fmt.Errorf("ServiceOrganization/ConvertSliceToRegular: %w", err)
		}

		organization.FuncName = funcName.FuncName
		organization.OrganizationList = hitech

		organizations = append(organizations, &organization)
	}

	return &organizations, nil
}
