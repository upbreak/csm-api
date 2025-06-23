package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

// func: 프로젝트에 설정된 공수 조회
// @param
// - jno: 프로젝트pk
func (r *Repository) GetManHourList(ctx context.Context, db Queryer, jno int64) (*entity.ManHours, error) {
	manHours := entity.ManHours{}

	query := `
			SELECT
			    MHNO,
			    WORK_HOUR,
			    MAN_HOUR,
			    JNO
			FROM 
			    IRIS_MAN_HOUR MH
			WHERE
			    MH.DEL_YN = 'N' AND
				MH.JNO = :1
			ORDER BY
			    WORK_HOUR ASC
		`

	if err := db.SelectContext(ctx, &manHours, query, jno); err != nil {
		return nil, fmt.Errorf("GetManHourList err:%v", err)
	}

	return &manHours, nil
}

// func: 공수 수정 및 추가
// @param
// - manHour: 공수 정보
func (r *Repository) MergeManHour(ctx context.Context, tx Execer, manHour entity.ManHour) (err error) {
	query := `
		MERGE INTO SAFE.IRIS_MAN_HOUR J1
		USING (
			SELECT 
				:1 AS MHNO,
				:2 AS WORK_HOUR,
				:3 AS MAN_HOUR,
				:4 AS JNO, 
				:5 AS UNO,	
				:6 AS USER_NAME
			FROM DUAL
		) J2
		ON (
			J1.DEL_YN = 'N' AND
			J1.MHNO = J2.MHNO
		) WHEN MATCHED THEN
			UPDATE SET
				J1.WORK_HOUR = J2.WORK_HOUR,
				J1.MAN_HOUR = J2.MAN_HOUR,
				J1.JNO = J2.JNO,
				J1.MOD_UNO = J2.UNO,	
				J1.MOD_USER = J2.USER_NAME,
				J1.MOD_DATE = SYSDATE
		WHEN NOT MATCHED THEN
			INSERT ( MHNO, WORK_HOUR, MAN_HOUR, JNO, REG_UNO, REG_USER, REG_DATE )
			VALUES (
				SEQ_IRIS_MAN_HOUR.NEXTVAL,
				J2.WORK_HOUR,
				J2.MAN_HOUR,
				J2.JNO,
				J2.UNO,	
				J2.USER_NAME,
				SYSDATE
			)
		`
	if _, err = tx.ExecContext(ctx, query, manHour.Mhno, manHour.WorkHour, manHour.ManHour, manHour.Jno, manHour.RegUno, manHour.RegUser); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("MargeManHour err: %w", err)

	}

	return
}

// func: 프로젝트 설정 정보 수정
// @param: ProjectSetting
// -
func (r *Repository) MergeProjectSetting(ctx context.Context, tx Execer, project entity.ProjectSetting) error {
	//agent := utils.GetAgent()

	query := `
				MERGE INTO IRIS_JOB_SET J1
				USING (
					SELECT 
						:1 AS JNO,
						:2 AS IN_TIME,
						:3 AS OUT_TIME,
						:4 AS RESPITE_TIME,
						:5 AS CANCEL_CODE,
						:6 AS UNO,	
						:7 AS USER_NAME
					FROM DUAL
				) J2
				ON (
					J1.JNO = J2.JNO
				) WHEN MATCHED THEN
					UPDATE SET
						J1.IN_TIME = J2.IN_TIME,
						J1.OUT_TIME = J2.OUT_TIME,
						J1.RESPITE_TIME = J2.RESPITE_TIME,
						J1.CANCEL_CODE = J2.CANCEL_CODE,
						J1.MOD_UNO = J2.UNO,	
						J1.MOD_USER = J2.USER_NAME,
						J1.MOD_DATE = SYSDATE
				WHEN NOT MATCHED THEN
					INSERT ( JNO, IN_TIME, OUT_TIME, RESPITE_TIME, CANCEL_CODE, REG_UNO, REG_USER, REG_DATE )
					VALUES (
						J2.JNO,
						J2.IN_TIME,
						J2.OUT_TIME,
						J2.RESPITE_TIME,
						J2.CANCEL_CODE,
						J2.UNO,	
						J2.USER_NAME,
						SYSDATE		
			)`
	if _, err := tx.ExecContext(ctx, query, project.Jno, project.InTime, project.OutTime, project.RespiteTime, project.CancelCode, project.RegUno, project.RegUser); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("MergeProject. Failed to modify project setting: %w", err)
	}

	return nil
}

// func: 프로젝트 미설정 정보 조회(스케줄러)
// @param
// -
func (r *Repository) GetCheckProjectSetting(ctx context.Context, db Queryer) (projects *entity.ProjectSettings, err error) {
	projects = &entity.ProjectSettings{}

	query := `
				SELECT 
				    DISTINCT(JNO) 
				FROM 
				    IRIS_SITE_JOB 
				WHERE 
				    JNO NOT IN (SELECT JNO FROM IRIS_JOB_SET)`

	if err = db.SelectContext(ctx, projects, query); err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("GetCheckProjectSetting err: %w", err)
	}

	return
}

// func: 프로젝트 기본 설정 정보 조회
// @param
// - jno
func (r *Repository) GetProjectSetting(ctx context.Context, db Queryer, jno int64) (*entity.ProjectSettings, error) {
	setting := entity.ProjectSettings{}
	query := fmt.Sprintf(`
			SELECT 
				J.JNO,
				J.IN_TIME,
				J.OUT_TIME,
				J.RESPITE_TIME,
				J.CANCEL_CODE,
				J.REG_DATE,
				J.REG_UNO,
				J.REG_USER,
				J.MOD_DATE,
				J.MOD_UNO,
				J.MOD_USER
			FROM IRIS_JOB_SET J
			WHERE
				J.JNO = :1
			`)

	if err := db.SelectContext(ctx, &setting, query, jno); err != nil {
		//TODO: 에러 아카이브
		return &setting, fmt.Errorf("GetProjectSetting fail: %v", err)
	}

	return &setting, nil
}

// func: 공수 삭제
// @param
// - mhno: 공수pk
func (r *Repository) DeleteManHour(ctx context.Context, tx Execer, mhno int64) error {
	query := fmt.Sprintf(`
			UPDATE IRIS_MAN_HOUR
			SET
			    DEL_YN = 'Y'
			WHERE
			    MHNO = :1			    
			`)
	if _, err := tx.ExecContext(ctx, query, mhno); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("DeleteManHour fail: %v", err)
	}

	return nil
}
