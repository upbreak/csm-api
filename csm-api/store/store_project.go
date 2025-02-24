package store

import (
	"context"
	"csm-api/entity"
	"database/sql"
	"fmt"
)

func (r *Repository) GetProjectList(ctx context.Context, db Queryer, sno int64) (*entity.ProjectInfoSqls, error) {
	projectInfoSqls := entity.ProjectInfoSqls{}

	// sno 변환: 0이면 NULL 처리, 아니면 Valid 값으로 설정
	var snoParam sql.NullInt64
	if sno != 0 {
		snoParam = sql.NullInt64{Valid: true, Int64: sno}
	} else {
		snoParam = sql.NullInt64{Valid: false}
	}

	sql := `SELECT
				t1.SNO,
				t1.JNO,
				t1.IS_USE,
				t1.IS_DEFAULT,
				t1.REG_DATE,
				t1.REG_USER,
				t1.REG_UNO,
				t1.MOD_DATE,
				t1.MOD_USER,
				t1.MOD_UNO,
				t2.JOB_TYPE as PROJECT_TYPE,
				t2.JOB_NO as PROJECT_NO,
				t2.JOB_NAME as PROJECT_NM,
				t2.JOB_YEAR as PROJECT_YEAR,
				t2.JOB_LOC as PROJECT_LOC,
				t2.JOB_CODE as PROJECT_CODE,
				t4.KIND_NAME as PROJECT_CODE_NAME,
				t3.SITE_NM,
				t2.COMP_CODE,
				t2.COMP_NICK,
				t2.COMP_NAME,
				t2.COMP_ETC,
				t2.ORDER_COMP_CODE,
				t2.ORDER_COMP_NICK,
				t2.ORDER_COMP_NAME,
				t2.ORDER_COMP_JOB_NAME,
				t2.JOB_LOC_NAME as PROJECT_LOC_NAME,
				t2.JOB_PM,
				t2.JOB_PE,
				TO_DATE(t2.JOB_SD, 'YYYY-MM-DD') as PROJECT_STDT,
				TO_DATE(t2.JOB_ED, 'YYYY-MM-DD') as PROJECT_EDDT,
				TO_DATE(t2.REG_DATE, 'YYYY-MM-DD HH24:MI:SS') as PROJECT_REG_DATE,
				TO_DATE(t2.MOD_DATE, 'YYYY-MM-DD HH24:MI:SS') as PROJECT_MOD_DATE,
				t2.JOB_STATE as PROJECT_STATE,
				t2.MOC_NO,
				t2.WO_NO
			FROM
				IRIS_SITE_JOB t1
				INNER JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
				INNER JOIN IRIS_SITE_SET t3 ON t1.SNO = t3.SNO
				INNER JOIN TIMESHEET.JOB_KIND_CODE t4 ON t2.JOB_CODE = t4.KIND_CODE
			WHERE
				t1.SNO > 100
				AND (:1 IS NULL OR t1.SNO = :2)
			ORDER BY
				t1.IS_DEFAULT DESC`
	if err := db.SelectContext(ctx, &projectInfoSqls, sql, snoParam, snoParam); err != nil {
		return &projectInfoSqls, fmt.Errorf("getProjectList fail: %v", err)
	}

	return &projectInfoSqls, nil
}

func (r *Repository) GetProjectNmList(ctx context.Context, db Queryer) (*entity.ProjectInfoSqls, error) {
	projectInfoSqls := entity.ProjectInfoSqls{}

	sql := `SELECT
    			t1.SNO,
				t1.JNO,
				t2.JOB_NAME as PROJECT_NM
			FROM
				IRIS_SITE_JOB t1
				INNER JOIN S_JOB_INFO t2 ON t1.JNO = t2.JNO
				INNER JOIN IRIS_SITE_SET t3 ON t1.SNO = t3.SNO
				INNER JOIN TIMESHEET.JOB_KIND_CODE t4 ON t2.JOB_CODE = t4.KIND_CODE
			WHERE t1.sno > 100
			ORDER BY
				t1.IS_DEFAULT DESC`
	if err := db.SelectContext(ctx, &projectInfoSqls, sql); err != nil {
		return &projectInfoSqls, fmt.Errorf("GetProjectNmList fail: %v", err)
	}

	return &projectInfoSqls, nil
}
