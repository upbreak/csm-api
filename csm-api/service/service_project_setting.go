package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"github.com/guregu/null"
	"time"
)

type ServiceProjectSetting struct {
	SafeDB        store.Queryer
	SafeTDB       store.Beginner
	Store         store.ProjectSettingStore
	WorkHourStore store.WorkHourStore
}

// func: 프로젝트에 설정된 공수 조회
// @param
// - jno: 프로젝트pk
func (s *ServiceProjectSetting) GetManHourList(ctx context.Context, jno int64) (*entity.ManHours, error) {

	manhours, err := s.Store.GetManHourList(ctx, s.SafeDB, jno)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.ManHours{}, fmt.Errorf("service_manHour/GetManHourList err: %w", err)
	}

	return manhours, nil

}

// func: 공수 수정 및 추가 (수정 시 기존 공수 삭제 후 새로 넣는 방식)
// @param
// - manHours: 공수 정보 배열
func (s *ServiceProjectSetting) MergeManHours(ctx context.Context, manHours *entity.ManHours) error {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// TODO: 에러 아카이브
				err = fmt.Errorf("service_project_setting/MergeProjectSetting Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				// TODO: 에러 아카이브
				err = fmt.Errorf("service_project_setting/MergeProjectSetting Commit error: %w", commitErr)
			}
		}
	}()

	jno := (*manHours)[0].Jno.Int64
	user := (*manHours)[0].Base

	// jno에 해당하는 공수 찾기
	deleteManhours, err := s.Store.GetManHourList(ctx, s.SafeDB, jno)
	if err != nil {
		// TODO: 에러 아카이브
		return fmt.Errorf("MergeManHours/GetManHourList err: %w", err)
	}

	// jno에 해당하는 공수 모두 삭제
	for _, deleteManhour := range *deleteManhours {
		deleteManhour.Message.Valid = true
		deleteManhour.Message.String = fmt.Sprintf(`[DELETE] mhno:[before:%d, after: N/A]|work_hour:[before: %d, after: N/A]|man_hour:[before:%.2f, after: N/A]|jno:[before:%d, after: N/A]|etc:[before:%s, after: N/A]`, deleteManhour.Mhno.Int64, deleteManhour.WorkHour.Int64, deleteManhour.ManHour.Float64, deleteManhour.Jno.Int64, deleteManhour.Etc.String)
		deleteManhour.Base = user

		// 삭제
		if err = s.DeleteManHour(ctx, deleteManhour.Mhno.Int64, *deleteManhour); err != nil {
			// TODO: 에러 아카이브
			return fmt.Errorf("MergeManHours err: %w", err)
		}
	}

	// 받아온 공수 추가
	for _, manHour := range *manHours {
		if !manHour.Message.Valid {
			continue
		}

		// 추가
		err = s.Store.AddManHour(ctx, tx, *manHour)
		if err != nil {
			// TODO: 에러 아카이브
			return fmt.Errorf("service_project_setting/AddManHour error: %w", err)
		}

		// 로그 남기기
		if err = s.Store.ManHourLog(ctx, tx, *manHour); err != nil {
			// TODO: 에러 아카이브
			return fmt.Errorf("service_project_setting/MergeManHour error: %w", err)
		}
	}

	// 공수에 맞춰 근로자 업데이트
	if err = s.WorkHourStore.ModifyWorkHourByJno(ctx, tx, jno, user, nil); err != nil {
		// TODO: 에러 아카이브
		return fmt.Errorf("service_project_setting/WorkHourStore/: %w", err)
	}

	return nil
}

// func: 프로젝트 설정 정보 추가 및 수정
// @param
// - ProjectSetting
func (s *ServiceProjectSetting) MergeProjectSetting(ctx context.Context, project entity.ProjectSetting) (err error) {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		// TODO: 에러 아카이브
		return fmt.Errorf("service_project_setting/ModifyProjectSetting BeginTx error: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// TODO: 에러 아카이브
				err = fmt.Errorf("service_project_setting/ModifyProjectSetting Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				// TODO: 에러 아카이브
				err = fmt.Errorf("service_project_setting/ModifyProjectSetting Commit error: %w", commitErr)
			}
		}
	}()

	if !project.Message.Valid {
		return
	}

	count, err := s.Store.MergeProjectSetting(ctx, tx, project)
	if err != nil {
		// TODO: 에러 아카이브
		return fmt.Errorf("service_project_setting/ModifyProjectSetting error: %w", err)
	}

	if count <= 0 {
		return
	} else {
		// 프로젝트 재설정 시 근로자 업데이트
		jno := project.Jno.Int64
		user := project.Base
		if err = s.WorkHourStore.ModifyWorkHourByJno(ctx, tx, jno, user, nil); err != nil {
			// TODO: 에러 아카이브
			return fmt.Errorf("service_project_setting/WorkHourStore error: %w", err)
		}
	}

	if err = s.Store.ProjectSettingLog(ctx, tx, project); err != nil {
		// TODO: 에러 아카이브
		return fmt.Errorf("service_project_setting/MergeManHour error: %w", err)
	}

	return
}

// func: 프로젝트 미설정 정보 업데이트(스케줄러)
// @param
// -
func (s *ServiceProjectSetting) CheckProjectSetting(ctx context.Context) (count int, err error) {

	projects := &entity.ProjectSettings{}
	if projects, err = s.Store.GetCheckProjectSetting(ctx, s.SafeDB); err != nil {
		// TODO: 에러 아카이브
		return 0, fmt.Errorf("service_project_setting/CheckProjectSetting error: %w", err)
	}

	for _, project := range *projects {

		// 기본 공수 추가하기
		manHourMore := &entity.ManHour{}
		manHourLess := &entity.ManHour{}

		manHourMore.WorkHour = utils.ParseNullInt("8")
		manHourMore.ManHour = utils.ParseNullFloat("0.5")
		manHourMore.Jno = project.Jno

		manHours := entity.ManHours{manHourMore, manHourLess}

		if err = s.MergeManHours(ctx, &manHours); err != nil {
			// TODO: 에러 아카이브
			return 0, fmt.Errorf("service_manhours/MergeManHours error: %w", err)
		}

		// 프로젝트 기본값으로 설정하기
		setting := &entity.ProjectSetting{}

		setting.Jno = project.Jno
		loc, _ := time.LoadLocation("Asia/Seoul")
		setting.InTime = null.NewTime(time.Date(2006, 01, 02, 8, 0, 0, 0, loc), true)
		setting.OutTime = null.NewTime(time.Date(2006, 01, 02, 17, 0, 0, 0, loc), true)
		setting.RespiteTime = utils.ParseNullInt("30")
		setting.CancelCode = utils.ParseNullString("NO_DAY")

		if err = s.MergeProjectSetting(ctx, *setting); err != nil {
			// TODO: 에러 아카이브
			return 0, fmt.Errorf("service_project_setting/CheckProjectSetting error: %w", err)
		}

	}

	count = len(*projects)
	return
}

// func: 프로젝트 설정 정보 가져오기
// @param
// - jno: 프로젝트PK
func (s *ServiceProjectSetting) GetProjectSetting(ctx context.Context, jno int64) (*entity.ProjectSettings, error) {

	setting, err := s.Store.GetProjectSetting(ctx, s.SafeDB, jno)
	if err != nil {
		//TODO: 에러 아카이브
		return &entity.ProjectSettings{}, fmt.Errorf("service_project_setting/GetProjectSetting: %w", err)
	}

	manHours, err := s.GetManHourList(ctx, jno)
	if err != nil {
		// TODO: 에러 아카이브
		return &entity.ProjectSettings{}, fmt.Errorf("service_project_setting/GetProjectSetting: %w", err)
	}

	if len(*setting) > 0 {
		(*setting)[0].ManHours = manHours
	}

	return setting, nil

}

// func: 공수 삭제
// @param
// - mhno: 공수pk
func (s *ServiceProjectSetting) DeleteManHour(ctx context.Context, mhno int64, manhour entity.ManHour) error {
	tx, err := s.SafeTDB.BeginTx(ctx, nil)
	if err != nil {
		// TODO: 에러 아카이브
		return fmt.Errorf("service_project_setting/ModifyProjectSetting BeginTx error: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// TODO: 에러 아카이브
				err = fmt.Errorf("service_project_setting/ModifyProjectSetting Rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				// TODO: 에러 아카이브
				err = fmt.Errorf("service_project_setting/ModifyProjectSetting Commit error: %w", commitErr)
			}
		}
	}()

	// 공수 삭제
	if err = s.Store.DeleteManHour(ctx, tx, mhno); err != nil {
		// TODO: 에러 아카이브
		return fmt.Errorf("service_project_setting/DeleteManHour error: %w", err)
	}

	// 공수 삭제 시 근로자 업데이트
	jno := manhour.Jno.Int64
	user := manhour.Base
	if err = s.WorkHourStore.ModifyWorkHourByJno(ctx, tx, jno, user, nil); err != nil {
		// TODO: 에러 아카이브
		return fmt.Errorf("service_project_setting/WorkHourStore error: %w", err)
	}

	// 로그 기록
	if err = s.Store.ManHourLog(ctx, tx, manhour); err != nil {
		// TODO: 에러 아카이브
		return fmt.Errorf("service_project_setting/MergeManHour error: %w", err)
	}
	return nil
}
