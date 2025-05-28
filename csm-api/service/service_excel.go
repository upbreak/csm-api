package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"github.com/xuri/excelize/v2"
)

type ServiceExcel struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.ExcelStore
}

// TBM excel import
func (s *ServiceExcel) ImportTbm(ctx context.Context, path string, tbm entity.Tbm) (err error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("ImportTbm.failed to open Excel file: %w", err)
	}

	var tbmList []entity.Tbm

	type cellPair struct {
		nameCol string
		signCol string
	}

	pairs := []cellPair{
		{"C", "D"},
		{"G", "H"},
		{"K", "L"},
	}

	// 시트 전체 순회
	for _, sheetName := range f.GetSheetList() {
		for row := 25; row <= 34; row++ {
			for _, pair := range pairs {
				nameCell := pair.nameCol + fmt.Sprint(row)
				signCell := pair.signCol + fmt.Sprint(row)

				name, _ := f.GetCellValue(sheetName, nameCell)
				sign, _ := f.GetCellValue(sheetName, signCell)

				if name != "" && sign != "" {
					newTbm := entity.Tbm{
						Jno:        tbm.Jno,
						Department: tbm.Department,
						DiscName:   tbm.DiscName,
						TbmDate:    tbm.TbmDate,
						UserNm:     utils.ParseNullString(name),
						Base: entity.Base{
							RegUser: tbm.RegUser,
							RegUno:  tbm.RegUno,
						},
					}
					tbmList = append(tbmList, newTbm)
				}
			}
		}
	}

	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("ImportTbm.failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("ImportTbm.failed to rollback transaction: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("ImportTbm.failed to commit transaction: %w", commitErr)
			}
		}
	}()

	// 기존 db 삭제
	if err = s.Store.ModifyTbmExcel(ctx, tx, tbm); err != nil {
		return fmt.Errorf("ImportTbm.failed to modify tbm: %w", err)
	}

	// db 저장
	if err = s.Store.AddTbmExcel(ctx, tx, tbmList); err != nil {
		return fmt.Errorf("ImportTbm.failed to add tbm sheet: %w", err)
	}

	return
}

// 퇴직공제 excel import
func (s *ServiceExcel) ImportDeduction(ctx context.Context, path string, deduction entity.Deduction) error {
	return nil
}
