package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"fmt"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-17
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

type ServiceWorker struct {
	SafeDB  store.Queryer
	SafeTDB store.Beginner
	Store   store.WorkerStore
}

// func: 전체 근로자 조회
// @param
// - page entity.PageSql: 정렬, 리스트 수
// - search entity.WorkerSql: 검색 단어
// - retry string: 통합검색 텍스트
func (s *ServiceWorker) GetWorkerTotalList(ctx context.Context, page entity.Page, search entity.Worker, retry string) (*entity.Workers, error) {
	// regular type ->  sql type 변환
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		return nil, fmt.Errorf("service_worker;total/OfPageSql err: %v", err)
	}

	// 조회
	list, err := s.Store.GetWorkerTotalList(ctx, s.SafeDB, pageSql, search, retry)
	if err != nil {
		return nil, fmt.Errorf("service_worker/GetWorkerTotalList err: %v", err)
	}

	return list, nil
}

// func: 전체 근로자 개수 조회
// @param
// - searchTime string: 조회 날짜
// - retry string: 통합검색 텍스트
func (s *ServiceWorker) GetWorkerTotalCount(ctx context.Context, search entity.Worker, retry string) (int, error) {
	count, err := s.Store.GetWorkerTotalCount(ctx, s.SafeDB, search, retry)
	if err != nil {
		return 0, fmt.Errorf("service_worker/GetWorkerTotalCount err: %v", err)
	}
	return count, nil
}

// func: 근로자 검색(현장근로자 추가시 사용)
// @param
// - userId string
func (s *ServiceWorker) GetAbsentWorkerList(ctx context.Context, page entity.Page, search entity.WorkerDaily, retry string) (*entity.Workers, error) {
	// regular type ->  sql type 변환
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		return nil, fmt.Errorf("service_worker;ByUserId/OfPageSql err: %v", err)
	}

	// 조회
	list, err := s.Store.GetAbsentWorkerList(ctx, s.SafeDB, pageSql, search, retry)
	if err != nil {
		return nil, fmt.Errorf("service_worker/GetAbsentWorkerList err: %v", err)
	}

	return list, nil
}

// func: 근로자 개수 검색(현장근로자 추가시 사용)
// @param
// - userId string
func (s *ServiceWorker) GetAbsentWorkerCount(ctx context.Context, search entity.WorkerDaily, retry string) (int, error) {
	count, err := s.Store.GetAbsentWorkerCount(ctx, s.SafeDB, search, retry)
	if err != nil {
		return 0, fmt.Errorf("service_worker;ByUserId/GetAbsentWorkerCount err: %v", err)
	}
	return count, nil
}

// 프로젝트에 참여한 회사명 리스트
func (s *ServiceWorker) GetWorkerDepartList(ctx context.Context, jno int64) ([]string, error) {
	list, err := s.Store.GetWorkerDepartList(ctx, s.SafeDB, jno)
	if err != nil {
		return nil, fmt.Errorf("service_worker;ByUserId/GetWorkerDepartList err: %v", err)
	}
	return list, nil
}

// func: 근로자 추가
// @param
// -
func (s *ServiceWorker) AddWorker(ctx context.Context, worker entity.Worker) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_worker;AddWorker err: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_worker;AddWorker panic err: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker;AddWorker err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker;AddWorker err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	err = s.Store.AddWorker(ctx, tx, worker)
	if err != nil {
		return fmt.Errorf("service_worker/AddWorker err: %v", err)
	}
	return
}

// func: 근로자 수정
// @param
// -
func (s *ServiceWorker) ModifyWorker(ctx context.Context, worker entity.Worker) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_worker;ModifyWorker err: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_worker;ModifyWorker panic err: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	err = s.Store.ModifyWorker(ctx, tx, worker)
	if err != nil {
		return fmt.Errorf("service_worker/ModifyWorker err: %v", err)
	}
	return
}

// func: 현장 근로자 조회
// @param
// - page entity.PageSql: 정렬, 리스트 수
// - search entity.WorkerSql: 검색 단어
func (s *ServiceWorker) GetWorkerSiteBaseList(ctx context.Context, page entity.Page, search entity.WorkerDaily, retry string) (*entity.WorkerDailys, error) {
	// regular type ->  sql type 변환
	pageSql := entity.PageSql{}
	pageSql, err := pageSql.OfPageSql(page)
	if err != nil {
		return nil, fmt.Errorf("service_worker;site_base/OfPageSql err: %v", err)
	}

	// 조회
	list, err := s.Store.GetWorkerSiteBaseList(ctx, s.SafeDB, pageSql, search, retry)
	if err != nil {
		return nil, fmt.Errorf("service_worker/GetWorkerSiteBaseList err: %v", err)
	}

	return list, nil
}

// func: 현장 근로자 개수 조회
// @param
// - searchTime string: 조회 날짜
func (s *ServiceWorker) GetWorkerSiteBaseCount(ctx context.Context, search entity.WorkerDaily, retry string) (int, error) {
	count, err := s.Store.GetWorkerSiteBaseCount(ctx, s.SafeDB, search, retry)
	if err != nil {
		return 0, fmt.Errorf("service_worker/GetWorkerSiteBaseCount err: %v", err)
	}
	return count, nil
}

// func: 현장 근로자 추가/수정
// @param
// -
func (s *ServiceWorker) MergeSiteBaseWorker(ctx context.Context, workers entity.WorkerDailys) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_worker;MergeSiteBaseWorker err: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_worker;MergeSiteBaseWorker panic err: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker;MergeWorker err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker;MergeWorker err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	// 추가/수정
	if err = s.Store.MergeSiteBaseWorker(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker/MergeSiteBaseWorker err: %v", err)
	}

	// 변경사항 로그 저장
	if err = s.Store.MergeSiteBaseWorkerLog(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker/MergeSiteBaseWorker err: %v", err)
	}

	return
}

// func: 현장 근로자 일괄마감
// @param
// -
func (s *ServiceWorker) ModifyWorkerDeadline(ctx context.Context, workers entity.WorkerDailys) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_worker;ModifyWorkerDeadline err: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_worker;ModifyWorkerDeadline panic err: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	// 마감처리
	if err = s.Store.ModifyWorkerDeadline(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker/ModifyWorkerDeadline err: %v", err)
	}

	// 마감 로그 저장
	if err = s.Store.MergeSiteBaseWorkerLog(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker/MergeSiteBaseWorker err: %v", err)
	}
	return
}

// func: 현장 근로자 프로젝트 변경
// @param
// -
func (s *ServiceWorker) ModifyWorkerProject(ctx context.Context, workers entity.WorkerDailys) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_worker;ModifyWorkerProject err: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			if r := recover(); r != nil {
				_ = tx.Rollback()
				err = fmt.Errorf("service_site/AddSite panic: %v", r)
				return
			}
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; commit err: %v", err, commitErr)
			}
		}
	}()
	// 전체 근로자 프로젝트 변경
	if err = s.Store.ModifyWorkerDefaultProject(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker/ModifyWorkerDefaultProject err: %v", err)
	}

	// 현장 근로자 프로젝트 변경
	if err = s.Store.ModifyWorkerProject(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker/ModifyWorkerProject err: %v", err)
	}

	// 프로젝트 변경 로그 저장
	if err = s.Store.MergeSiteBaseWorkerLog(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker/MergeSiteBaseWorker err: %v", err)
	}

	return
}

// func: 현장 근로자 일일 마감처리
// @param
// -
func (s *ServiceWorker) ModifyWorkerDeadlineInit(ctx context.Context) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_worker;ModifyWorkerDeadlineInit err: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_worker;ModifyWorkerDeadlineInit panic err: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	if err = s.Store.ModifyWorkerDeadlineInit(ctx, tx); err != nil {
		return fmt.Errorf("service_worker/ModifyWorkerDeadlineInit err: %v", err)
	}
	return
}

// func: 현장 근로자 철야 처리
// @param
// -
func (s *ServiceWorker) ModifyWorkerOverTime(ctx context.Context) (count int, err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("service_worker;ModifyWorkerOverTime err: %v", err)
	}

	// 철야 근로자 존재 여부 확인
	workerOverTimes := &entity.WorkerOverTimes{}
	workerOverTimes, err = s.Store.GetWorkerOverTime(ctx, s.SafeDB)
	if err != nil {
		return 0, fmt.Errorf("service_worker;GetWorkerOverTime err: %v", err)
	}
	count = len(*workerOverTimes)

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_worker;ModifyWorkerOverTime panic err: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorker err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	for _, workerOverTime := range *workerOverTimes {

		// 철야 근로자 철야 표시 및 퇴근시간 합치기.
		if err = s.Store.ModifyWorkerOverTime(ctx, tx, *workerOverTime); err != nil {
			return 0, fmt.Errorf("service_worker/ModifyWorkerOverTime err: %v", err)
		}

		// 다음날 퇴근 표시 삭제
		if err = s.Store.DeleteWorkerOverTime(ctx, tx, (*workerOverTime).AfterCno); err != nil {
			return 0, fmt.Errorf("service_worker;DeleteWorkerOverTime err: %v", err)
		}
	}
	return
}

// 현장 근로자 삭제
func (s *ServiceWorker) RemoveSiteBaseWorkers(ctx context.Context, workers entity.WorkerDailys) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_worker.RemoveSiteBaseWorkers BeginTx err: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_worker.RemoveSiteBaseWorkers panic err: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker.RemoveSiteBaseWorkers err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker.RemoveSiteBaseWorkers err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	// 현장 근로자 삭제
	if err = s.Store.RemoveSiteBaseWorkers(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker.RemoveSiteBaseWorkers err: %v", err)
	}

	// 삭제 로그 저장
	if err = s.Store.MergeSiteBaseWorkerLog(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker.MergeSiteBaseWorkerLog err: %v", err)
	}

	return
}

// 마감 취소
func (s *ServiceWorker) ModifyDeadlineCancel(ctx context.Context, workers entity.WorkerDailys) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_worker.ModifyDeadlineCancel BeginTx err: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_worker.ModifyDeadlineCancel panic err: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker.ModifyDeadlineCancel err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker.ModifyDeadlineCancel err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	// 마감 취소
	if err = s.Store.ModifyDeadlineCancel(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker.ModifyDeadlineCancel err: %v", err)
	}

	// 마감 취소 로그 저장
	if err = s.Store.MergeSiteBaseWorkerLog(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker.MergeSiteBaseWorkerLog err: %v", err)
	}

	return
}

// 프로젝트, 기간내 모든 현장근로자 근태정보 조회
func (s *ServiceWorker) GetDailyWorkersByJnoAndDate(ctx context.Context, param entity.RecordDailyWorkerReq) ([]entity.RecordDailyWorkerRes, error) {
	list, err := s.Store.GetDailyWorkersByJnoAndDate(ctx, s.SafeDB, param)
	if err != nil {
		return []entity.RecordDailyWorkerRes{}, fmt.Errorf("service_worker;GetDailyWorkersByJnoAndDate err: %v", err)
	}
	return list, nil
}

// 현장근로자 일괄 공수 변경
func (s *ServiceWorker) ModifyWorkHours(ctx context.Context, workers entity.WorkerDailys) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("service_worker;ModifyWorkHours err: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("service_worker;ModifyWorkHours panic err: %v", r)
			return
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorkHours err: %v; rollback err: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("service_worker;ModifyWorkHours err: %v; commit err: %v", err, commitErr)
			}
		}
	}()

	// 공수 변경
	if err = s.Store.ModifyWorkHours(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker/ModifyWorkHours err: %v", err)
	}

	// 공수 변경 로그 저장
	if err = s.Store.MergeSiteBaseWorkerLog(ctx, tx, workers); err != nil {
		return fmt.Errorf("service_worker.MergeSiteBaseWorkerLog err: %v", err)
	}

	return
}
