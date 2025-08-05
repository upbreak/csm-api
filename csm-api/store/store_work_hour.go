package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"fmt"
	"strings"
)

// 마감처리가 안된 특정프로젝트의 근로자의 공수 계산: jno는 필수, ids는 없으면 jno의 모든 근로자 계산 있으면 해당 id의 근로자만 계산
func (r *Repository) ModifyWorkHourByJno(ctx context.Context, tx Execer, jno int64, user entity.Base, uuids []string) error {
	var (
		query strings.Builder
		args  []interface{}
	)

	args = append(args, jno)

	query.WriteString(`
			MERGE INTO IRIS_WORKER_DAILY_SET T1
			USING (
				SELECT 
					T1.ROWID AS W_ROWID,
					T1.JNO,
					T1.RECORD_DATE,
					T1.IN_RECOG_TIME,
					T1.OUT_RECOG_TIME,
					T1.IS_OVERTIME,
					T2.IN_TIME,
					T2.OUT_TIME,
					T2.RESPITE_TIME,
					TO_CHAR(T2.IN_TIME + NUMTODSINTERVAL(T2.RESPITE_TIME, 'MINUTE'), 'HH24:MI') AS IN_LIMIT,
					TO_CHAR(T2.OUT_TIME - NUMTODSINTERVAL(T2.RESPITE_TIME, 'MINUTE'), 'HH24:MI') AS OUT_LIMIT,
					TO_CHAR(T1.IN_RECOG_TIME, 'HH24:MI') AS ACTUAL_IN,
					TO_CHAR(T1.OUT_RECOG_TIME, 'HH24:MI') AS ACTUAL_OUT,
					FLOOR((
						(
							CASE 
								WHEN T1.IS_OVERTIME = 'Y' 
								THEN TO_DATE(TO_CHAR(T1.RECORD_DATE + 1, 'YYYY-MM-DD') || ' ' || TO_CHAR(T1.OUT_RECOG_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI')
								ELSE TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' ' || TO_CHAR(T1.OUT_RECOG_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI')
							END
							-
							GREATEST(
								TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' ' || TO_CHAR(T1.IN_RECOG_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI'),
								TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' ' || TO_CHAR(T2.IN_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI')
							)
							-
							GREATEST(
								LEAST(
									CASE 
										WHEN T1.IS_OVERTIME = 'Y' 
										THEN TO_DATE(TO_CHAR(T1.RECORD_DATE + 1, 'YYYY-MM-DD') || ' ' || TO_CHAR(T1.OUT_RECOG_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI')
										ELSE TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' ' || TO_CHAR(T1.OUT_RECOG_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI')
									END,
									TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' 13:00', 'YYYY-MM-DD HH24:MI')
								)
								-
								GREATEST(
									GREATEST(
										TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' ' || TO_CHAR(T1.IN_RECOG_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI'),
										TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' ' || TO_CHAR(T2.IN_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI')
									),
									TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' 12:00', 'YYYY-MM-DD HH24:MI')
								),
								0
							)
						) * 24
					)) AS WORKED_HOUR
				FROM IRIS_WORKER_DAILY_SET T1
				JOIN IRIS_JOB_SET T2 ON T1.JNO = T2.JNO
				WHERE 
					TRUNC(T1.RECORD_DATE) < TRUNC(SYSDATE)
					AND T1.IN_RECOG_TIME IS NOT NULL
					AND T1.OUT_RECOG_TIME IS NOT NULL
					AND T1.JNO = :1
					AND T1.IS_DEADLINE = 'N'`)

	if len(uuids) > 0 {
		query.WriteString("\nAND T1.USER_KEY IN (")
		for i, id := range uuids {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString(fmt.Sprintf(":%d", i+2))
			args = append(args, id)
		}
		query.WriteString(")")
	}

	modUserIndex := len(args) + 1
	modUnoIndex := modUserIndex + 1

	args = append(args, user.ModUser, user.ModUno)

	query.WriteString(fmt.Sprintf(`
			) T3
			ON (T1.ROWID = T3.W_ROWID)
			WHEN MATCHED THEN
			UPDATE SET 
				T1.WORK_HOUR = (
					CASE
						WHEN T3.ACTUAL_IN <= T3.IN_LIMIT 
						 AND T3.ACTUAL_OUT >= T3.OUT_LIMIT THEN 1.0
						ELSE (
							SELECT NVL(MAX(T4.MAN_HOUR), 0)
							FROM IRIS_MAN_HOUR T4
							WHERE 
								T4.JNO = T3.JNO
								AND T4.WORK_HOUR <= T3.WORKED_HOUR
						)
					END
				),
				T1.MOD_DATE = SYSDATE,
				T1.MOD_USER = :%d,
				T1.MOD_UNO  = :%d`, modUserIndex, modUnoIndex))

	if _, err := tx.ExecContext(ctx, query.String(), args...); err != nil {
		return utils.CustomErrorf(err)
	}
	return nil
}

// 마감처이가 안되고 출퇴근이 둘다 있는 모든 근로자의 공수 계산
func (r *Repository) ModifyWorkHour(ctx context.Context, tx Execer, user entity.Base) error {

	query := `
			MERGE INTO IRIS_WORKER_DAILY_SET T1
			USING (
				SELECT 
					T1.ROWID AS W_ROWID,
					T1.JNO,
					--T1.USER_ID,
					T1.RECORD_DATE,
					T1.IN_RECOG_TIME,
					T1.OUT_RECOG_TIME,
					T1.IS_OVERTIME,
					T2.IN_TIME,
					T2.OUT_TIME,
					T2.RESPITE_TIME,
					TO_CHAR(T2.IN_TIME + NUMTODSINTERVAL(T2.RESPITE_TIME, 'MINUTE'), 'HH24:MI') AS IN_LIMIT,
					TO_CHAR(T2.OUT_TIME - NUMTODSINTERVAL(T2.RESPITE_TIME, 'MINUTE'), 'HH24:MI') AS OUT_LIMIT,
					TO_CHAR(T1.IN_RECOG_TIME, 'HH24:MI') AS ACTUAL_IN,
					TO_CHAR(T1.OUT_RECOG_TIME, 'HH24:MI') AS ACTUAL_OUT,
					GREATEST(
						TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' ' || TO_CHAR(T1.IN_RECOG_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI'),
						TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' ' || TO_CHAR(T2.IN_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI')
					) AS IN_DATETIME,
					CASE 
						WHEN T1.IS_OVERTIME = 'Y' 
						THEN TO_DATE(TO_CHAR(T1.RECORD_DATE + 1, 'YYYY-MM-DD') || ' ' || TO_CHAR(T1.OUT_RECOG_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI')
						ELSE TO_DATE(TO_CHAR(T1.RECORD_DATE, 'YYYY-MM-DD') || ' ' || TO_CHAR(T1.OUT_RECOG_TIME, 'HH24:MI'), 'YYYY-MM-DD HH24:MI')
					END AS OUT_DATETIME
			
				FROM IRIS_WORKER_DAILY_SET T1
				JOIN IRIS_JOB_SET T2 ON T1.JNO = T2.JNO
				WHERE 
					TRUNC(T1.RECORD_DATE) < TRUNC(SYSDATE)
					AND T1.IN_RECOG_TIME IS NOT NULL
					AND T1.OUT_RECOG_TIME IS NOT NULL
					AND T1.IS_DEADLINE = 'N'
					AND T1.COMPARE_STATE = 'S'
					AND T1.WORK_HOUR IS NULL 
			) T3
			ON (T1.ROWID = T3.W_ROWID)
			WHEN MATCHED THEN
			UPDATE SET 
				T1.WORK_HOUR = (
					CASE
						WHEN T3.ACTUAL_IN <= T3.IN_LIMIT 
						 AND T3.ACTUAL_OUT >= T3.OUT_LIMIT THEN 1.0
						ELSE (
							SELECT NVL(MAX(T4.MAN_HOUR), 0)
							FROM IRIS_MAN_HOUR T4
							WHERE 
								T4.JNO = T3.JNO
								AND T4.WORK_HOUR <= FLOOR((
									(T3.OUT_DATETIME - T3.IN_DATETIME)
									- GREATEST(
										LEAST(T3.OUT_DATETIME, TO_DATE(TO_CHAR(T3.RECORD_DATE, 'YYYY-MM-DD') || ' 13:00', 'YYYY-MM-DD HH24:MI'))
										- GREATEST(T3.IN_DATETIME, TO_DATE(TO_CHAR(T3.RECORD_DATE, 'YYYY-MM-DD') || ' 12:00', 'YYYY-MM-DD HH24:MI')),
										0
									)
								) * 24)
						)
					END
				),
				T1.MOD_DATE = SYSDATE,
				T1.MOD_USER = :1,
				T1.MOD_UNO  = :2`

	if _, err := tx.ExecContext(ctx, query, user.ModUser, user.ModUno); err != nil {
		return utils.CustomErrorf(err)
	}
	return nil
}
