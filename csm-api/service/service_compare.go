package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"github.com/guregu/null"
	"strings"
)

type ServiceCompare struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.CompareStore
}

// 일일 근로자 비교 리스트
func (s *ServiceCompare) GetCompareList(ctx context.Context, jno int64, startDate null.Time, retry string, order string) ([]entity.Compare, error) {
	workerlist, err := s.Store.GetDailyWorkerList(ctx, s.SafeDB, jno, startDate, retry, order)
	if err != nil {
		return nil, fmt.Errorf("ServiceCompare.GetCompareList GetDailyWorkerList :%w", err)
	}

	tbmList, err := s.Store.GetTbmList(ctx, s.SafeDB, jno, startDate, retry, order)
	if err != nil {
		return nil, fmt.Errorf("ServiceCompare.GetCompareList GetTbmList :%w", err)
	}

	deductionList, err := s.Store.GetDeductionList(ctx, s.SafeDB, jno, startDate, retry, order)
	if err != nil {
		return nil, fmt.Errorf("ServiceCompare.GetCompareList GetDeductionList :%w", err)
	}

	// 주민번호 " ", "-" 제거
	replacer := strings.NewReplacer(" ", "", "-", "")
	cleanRegNo := func(reg string) string {
		return replacer.Replace(reg)
	}

	// TBM map: 동명이인 고려한 slice map
	tbmMap := make(map[entity.TbmKey][]entity.Tbm)
	for _, tbm := range tbmList {
		key := entity.TbmKey{tbm.Jno.Int64, tbm.UserNm.String, tbm.Department.String, tbm.TbmDate.Time}
		tbmMap[key] = append(tbmMap[key], tbm)
	}

	// 공제 map: 주민번호 기준과 동명이인 기준
	deductionRegMap := make(map[entity.DeductionRegKey]entity.Deduction)
	deductionMap := make(map[entity.DeductionKey][]entity.Deduction)
	for _, d := range deductionList {
		regKey := entity.DeductionRegKey{
			Jno:        d.Jno.Int64,
			RegNo:      cleanRegNo(d.RegNo.String),
			RecordDate: d.RecordDate.Time,
		}
		deductionRegMap[regKey] = d

		key := entity.DeductionKey{
			Jno:        d.Jno.Int64,
			UserNm:     d.UserNm.String,
			Department: d.Department.String,
			RecordDate: d.RecordDate.Time,
		}
		deductionMap[key] = append(deductionMap[key], d)
	}

	var compareList []entity.Compare

	// 근태 기준 비교
	for _, worker := range workerlist {
		compare := entity.Compare{
			Jno:           worker.Jno,
			UserId:        worker.UserId,
			UserNm:        worker.UserNm,
			Department:    worker.Department,
			DiscName:      worker.DiscName,
			IsTbm:         utils.ParseNullString("N"),
			RecordDate:    worker.RecordDate,
			WorkerInTime:  worker.InRecogTime,
			WorkerOutTime: worker.OutRecogTime,
			CompareState:  worker.CompareState,
			IsDeadline:    worker.IsDeadline,
		}

		// TBM 비교
		tbmKey := entity.TbmKey{worker.Jno.Int64, worker.UserNm.String, worker.Department.String, worker.RecordDate.Time}
		if tbms, ok := tbmMap[tbmKey]; ok && len(tbms) > 0 {
			compare.IsTbm = utils.ParseNullString("Y")
			compare.DiscName = tbms[0].DiscName
			tbmMap[tbmKey] = tbms[1:] // 하나만 소비
			if len(tbmMap[tbmKey]) == 0 {
				delete(tbmMap, tbmKey)
			}
		}

		// 공제 비교 RegNo 기준
		deductKey := entity.DeductionRegKey{Jno: worker.Jno.Int64, RegNo: cleanRegNo(worker.RegNo.String), RecordDate: worker.RecordDate.Time}
		if deduction, ok := deductionRegMap[deductKey]; ok {
			compare.DeductionInTime = deduction.InRecogTime
			compare.DeductionOutTime = deduction.OutRecogTime
			delete(deductionRegMap, deductKey)

			dk := entity.DeductionKey{Jno: deduction.Jno.Int64, UserNm: deduction.UserNm.String, Department: deduction.Department.String, RecordDate: deduction.RecordDate.Time}
			if list := deductionMap[dk]; len(list) > 0 {
				deductionMap[dk] = list[1:]
				if len(deductionMap[dk]) == 0 {
					delete(deductionMap, dk)
				}
			}
		}

		compareList = append(compareList, compare)
	}

	// 근태x 남은 TBM-공제 비교
	for _, tbmSlice := range tbmMap {
		for _, tbm := range tbmSlice {
			compare := entity.Compare{
				Jno:          tbm.Jno,
				UserNm:       tbm.UserNm,
				Department:   tbm.Department,
				DiscName:     tbm.DiscName,
				IsTbm:        utils.ParseNullString("Y"),
				RecordDate:   tbm.TbmDate,
				CompareState: utils.ParseNullString("C"),
			}

			deductionKey := entity.DeductionKey{Jno: tbm.Jno.Int64, UserNm: tbm.UserNm.String, Department: tbm.Department.String, RecordDate: tbm.TbmDate.Time}
			if dList, ok := deductionMap[deductionKey]; ok && len(dList) > 0 {
				d := dList[0]
				compare.DeductionInTime = d.InRecogTime
				compare.DeductionOutTime = d.OutRecogTime
				deductionMap[deductionKey] = dList[1:]
				if len(deductionMap[deductionKey]) == 0 {
					delete(deductionMap, deductionKey)
				}

				// 삭제 동기화 (regMap도)
				regKey := entity.DeductionRegKey{Jno: d.Jno.Int64, RegNo: cleanRegNo(d.RegNo.String), RecordDate: d.RecordDate.Time}
				delete(deductionRegMap, regKey)
			}

			compareList = append(compareList, compare)
		}
	}

	// 근태x TBMx 공제 처리
	for _, dList := range deductionMap {
		for _, d := range dList {
			compare := entity.Compare{
				Jno:              d.Jno,
				UserNm:           d.UserNm,
				Department:       d.Department,
				IsTbm:            utils.ParseNullString("N"),
				CompareState:     utils.ParseNullString("C"),
				RecordDate:       d.RecordDate,
				DeductionInTime:  d.InRecogTime,
				DeductionOutTime: d.OutRecogTime,
			}
			compareList = append(compareList, compare)
		}
	}

	return compareList, nil
}

// 근로자 비교 반영/취소
func (s *ServiceCompare) ModifyWorkerCompareState(ctx context.Context, workers entity.WorkerDailys) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service.ModifyWorkerCompareState begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service.ModifyWorkerCompareState tx rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service.ModifyWorkerCompareState tx commit error: %w", commitErr)
			}
		}
	}()

	// 비교 상태 수정
	if err = s.Store.ModifyWorkerCompareState(ctx, tx, workers); err != nil {
		return fmt.Errorf("service.ModifyWorkerCompareState store error: %w", err)
	}

	// 비교 상태 수정 로그 등록
	if err = s.Store.AddCompareLog(ctx, tx, workers); err != nil {
		return fmt.Errorf("service.ModifyWorkerCompareState store error: %w", err)
	}

	return
}
