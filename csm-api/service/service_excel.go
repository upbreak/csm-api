package service

import (
	"context"
	"csm-api/ctxutil"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/txutil"
	"csm-api/utils"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

type ServiceExcel struct {
	SafeDB      store.Queryer
	SafeTDB     store.Beginner
	Store       store.ExcelStore
	WorkerStore store.WorkerStore
}

func mustGet(f *excelize.File, sheet, cell string) string {
	val, _ := f.GetCellValue(sheet, cell)
	return strings.TrimSpace(val)
}

// TBM excel import
func (s *ServiceExcel) ImportTbm(ctx context.Context, path string, tbm entity.Tbm) (err error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	order, err := s.Store.GetTbmOrder(ctx, s.SafeDB, tbm)
	if err != nil {
		return utils.CustomErrorf(err)
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
						Sno:        tbm.Sno,
						Jno:        tbm.Jno,
						Department: tbm.Department,
						DiscName:   tbm.DiscName,
						TbmDate:    tbm.TbmDate,
						TbmOrder:   utils.ParseNullInt(order),
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

	tx, ok := ctxutil.GetTx(ctx)
	if !ok || tx == nil {
		conn, err := s.SafeTDB.Conn(ctx)
		if err != nil {
			return utils.CustomErrorf(err)
		}
		defer func() {
			if closeErr := conn.Close(); closeErr != nil {
				if err != nil {
					err = utils.CustomMessageErrorf(fmt.Sprintf("%v; conn.Close", err), closeErr)
				} else {
					err = utils.CustomMessageErrorf("conn.Close", closeErr)
				}
			}
		}()
		defer txutil.DeferTxx(tx, &err)
	}

	// db 저장
	if err = s.Store.AddTbmExcel(ctx, tx, tbmList); err != nil {
		return utils.CustomErrorf(err)
	}

	return
}

// 퇴직공제 excel import
func (s *ServiceExcel) ImportDeduction(ctx context.Context, path string, deduction entity.Deduction) (err error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	order, err := s.Store.GetDeductionOrder(ctx, s.SafeDB, deduction)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	siteNm, err := s.Store.GetDeductionSiteNameBySno(ctx, s.SafeDB, deduction.Sno.Int64)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	sheetName := f.GetSheetName(0)

	var deductionList []entity.Deduction

	// 3번째 행부터 (1-based index, 엑셀은 B3부터 시작)
	for rowIdx := 3; ; rowIdx++ {
		// B열(근무날짜)
		cellAddr := fmt.Sprintf("B%d", rowIdx)
		dateStr, err := f.GetCellValue(sheetName, cellAddr)
		if err != nil || strings.TrimSpace(dateStr) == "" {
			break
		}
		if dateStr == "" || utils.ConvertMMDDYYToYYMMDD(dateStr) != deduction.RecordDate.Time.Format("06-01-02") {
			continue
		}

		// C열(현장명)
		siteName, _ := f.GetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx))
		if utils.NormalizeForEqual(siteName) != utils.NormalizeForEqual(siteNm) {
			continue
		}

		// G열(이름)
		userNm := mustGet(f, sheetName, fmt.Sprintf("G%d", rowIdx))
		// F열(회사명)
		department := mustGet(f, sheetName, fmt.Sprintf("F%d", rowIdx))
		// I열(성별)
		gender := mustGet(f, sheetName, fmt.Sprintf("N%d", rowIdx))
		// H열(생년월일)
		regNo := mustGet(f, sheetName, fmt.Sprintf("H%d", rowIdx))
		// I열(전화번호)
		phone := mustGet(f, sheetName, fmt.Sprintf("I%d", rowIdx))
		normalizedPhone := strings.ReplaceAll(strings.ReplaceAll(phone, "-", ""), " ", "")
		if len(normalizedPhone) == 10 && strings.HasPrefix(normalizedPhone, "1") {
			normalizedPhone = "0" + normalizedPhone
		}
		// O열(출근시간)
		inTime := mustGet(f, sheetName, fmt.Sprintf("O%d", rowIdx))
		// P열(퇴근시간)
		outTime := mustGet(f, sheetName, fmt.Sprintf("P%d", rowIdx))

		newDeduction := entity.Deduction{
			Sno:          deduction.Sno,
			UserNm:       utils.ParseNullString(userNm),
			Department:   utils.ParseNullString(department),
			Gender:       utils.ParseNullString(gender),
			RegNo:        utils.ParseNullString(utils.ConvertMMDDYYToYYMMDD(regNo)),
			Phone:        utils.ParseNullString(normalizedPhone),
			InRecogTime:  utils.ParseNullDateTime(deduction.RecordDate.Time.Format("2006-01-02"), inTime),
			OutRecogTime: utils.ParseNullDateTime(deduction.RecordDate.Time.Format("2006-01-02"), outTime),
			RecordDate:   deduction.RecordDate,
			DeductOrder:  utils.ParseNullString(order),
			Base: entity.Base{
				RegUser: deduction.RegUser,
				RegUno:  deduction.RegUno,
			},
		}

		if newDeduction.UserNm.Valid && newDeduction.Department.Valid && newDeduction.Gender.Valid && newDeduction.Phone.Valid {
			deductionList = append(deductionList, newDeduction)
		}
	}

	tx, ok := ctxutil.GetTx(ctx)
	if !ok || tx == nil {
		tx, cleanup, err := txutil.BeginTxWithCleanMode(ctx, s.SafeTDB, false)
		if err != nil {
			return utils.CustomErrorf(err)
		}

		defer func() {
			txutil.DeferTx(tx, &err)
			cleanup()
		}()
	}

	if err = s.Store.AddDeductionExcel(ctx, tx, deductionList); err != nil {
		return utils.CustomErrorf(err)
	}

	return
}

func (s *ServiceExcel) ImportAddDailyWorker(ctx context.Context, path string, worker entity.WorkerDaily) (err error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	sheet := f.GetSheetName(0)
	var excels []entity.WorkerDailyExcel

	row := 2
	for {
		// B열 (이름) 기준으로 값이 없으면 종료
		userNm, err := f.GetCellValue(sheet, fmt.Sprintf("B%d", row))
		if err != nil || strings.TrimSpace(userNm) == "" {
			break
		}

		department, _ := f.GetCellValue(sheet, fmt.Sprintf("C%d", row)) // 부서/조직명

		rawPhone, _ := f.GetCellValue(sheet, fmt.Sprintf("D%d", row)) // 핸드폰번호
		normalizedPhone := strings.ReplaceAll(strings.ReplaceAll(rawPhone, "-", ""), " ", "")
		if strings.HasPrefix(normalizedPhone, "1") {
			normalizedPhone = "0" + normalizedPhone
		}

		workDate, _ := f.GetCellValue(sheet, fmt.Sprintf("E%d", row)) // 날짜
		if !utils.IsYYYYMMDD(workDate) {
			workDate = utils.NormalizeYYMMDD(utils.ConvertMMDDYYToYYMMDD(workDate))
		}

		// F, G열 (출근/퇴근시간) → 시간 서식으로 저장됨
		inTimeRaw, err := f.GetCellValue(sheet, fmt.Sprintf("F%d", row))
		if err != nil {
			return utils.CustomErrorf(err)
		}
		inTime := inTimeRaw
		if timeVal, err := f.GetCellValue(sheet, fmt.Sprintf("F%d", row), excelize.Options{RawCellValue: false}); err == nil {
			inTime = timeVal
		}

		outTimeRaw, err := f.GetCellValue(sheet, fmt.Sprintf("G%d", row))
		if err != nil {
			return utils.CustomErrorf(err)
		}
		outTime := outTimeRaw
		if timeVal, err := f.GetCellValue(sheet, fmt.Sprintf("G%d", row), excelize.Options{RawCellValue: false}); err == nil {
			outTime = timeVal
		}

		workHour, _ := f.GetCellValue(sheet, fmt.Sprintf("H%d", row)) // 공수

		excels = append(excels, entity.WorkerDailyExcel{
			Department: department,
			UserNm:     userNm,
			Phone:      normalizedPhone,
			WorkDate:   workDate,
			InTime:     inTime,
			OutTime:    outTime,
			WorkHour:   workHour,
		})

		row++
	}

	var workers []entity.WorkerDaily
	for _, excel := range excels {
		temp := entity.WorkerDaily{
			Sno:          worker.Sno,
			Jno:          worker.Jno,
			Department:   utils.ParseNullString(excel.Department),
			UserNm:       utils.ParseNullString(excel.UserNm),
			UserId:       utils.ParseNullString(excel.Phone),
			RecordDate:   utils.ParseNullDate(excel.WorkDate),
			InRecogTime:  utils.ParseNullDateTime(excel.WorkDate, utils.NormalizeHHMM(excel.InTime)),
			OutRecogTime: utils.ParseNullDateTime(excel.WorkDate, utils.NormalizeHHMM(excel.OutTime)),
			WorkHour:     utils.ParseNullFloat(excel.WorkHour),
			CompareState: utils.ParseNullString("X"),
			WorkState:    utils.ParseNullString("02"),
			Base: entity.Base{
				RegUser: worker.RegUser,
				RegUno:  worker.RegUno,
			},
		}
		workers = append(workers, temp)
	}

	tx, cleanup, err := txutil.BeginTxWithCleanMode(ctx, s.SafeTDB, false)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer func() {
		txutil.DeferTx(tx, &err)
		cleanup()
	}()

	var list entity.WorkerDailys
	if list, err = s.WorkerStore.AddDailyWorkers(ctx, s.SafeDB, tx, workers); err != nil {
		return utils.CustomErrorf(err)
	}

	if err = s.WorkerStore.MergeSiteBaseWorkerLog(ctx, tx, list); err != nil {
		return utils.CustomErrorf(err)
	}

	return
}
