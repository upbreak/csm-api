package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"database/sql"
	"errors"
	"fmt"
)

// 업로드할 파일 차수
func (r *Repository) GetUploadRound(ctx context.Context, db Queryer, file entity.UploadFile) (int, error) {
	var maxRound sql.NullInt64

	query := `
		SELECT MAX(UPLOAD_ROUND)
		FROM IRIS_UPLOADED_FILES
		WHERE FILE_PATH = :1
		AND JNO = :2
		AND TRUNC(WORK_DATE) = TRUNC(:3)
		AND FILE_TYPE = :4`

	if err := db.GetContext(ctx, &maxRound, query, file.FilePath, file.Jno, file.WorkDate, file.FileType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("GetUploadRound: %w", err)
	}

	if !maxRound.Valid {
		return 0, nil
	}
	return int(maxRound.Int64 + 1), nil
}

// 업로드 파일 리스트
func (r *Repository) GetUploadFileList(ctx context.Context, db Queryer, file entity.UploadFile) ([]entity.UploadFile, error) {
	var uploadFileList []entity.UploadFile

	query := `
		SELECT 
			T1.FILE_TYPE,
			T1.FILE_PATH,
			T1.FILE_NAME,
			T1.UPLOAD_ROUND
		FROM IRIS_UPLOADED_FILES T1
		JOIN (
			SELECT 
				FILE_TYPE,
				MAX(UPLOAD_ROUND) AS UPLOAD_ROUND
			FROM IRIS_UPLOADED_FILES
			WHERE JNO = :2
			AND TRUNC(WORK_DATE) = TRUNC(:3)
			GROUP BY FILE_TYPE, FILE_PATH
		) T2 ON T1.FILE_TYPE = T2.FILE_TYPE AND T1.UPLOAD_ROUND = T2.UPLOAD_ROUND`

	if err := db.SelectContext(ctx, &uploadFileList, query, file.Jno, file.WorkDate); err != nil {
		return nil, fmt.Errorf("GetUploadFileList: %w", err)
	}
	return uploadFileList, nil
}

// 업로드 파일
func (r *Repository) GetUploadFile(ctx context.Context, db Queryer, file entity.UploadFile) (entity.UploadFile, error) {
	var uploadFile entity.UploadFile

	query := `
		SELECT 
		    MAX(FILE_PATH) as FILE_PATH,
			MAX(FILE_NAME) as FILE_NAME, 
			MAX(UPLOAD_ROUND) as UPLOAD_ROUND
		FROM IRIS_UPLOADED_FILES
		WHERE JNO = :1
		AND TRUNC(WORK_DATE) = TRUNC(:2)
		AND FILE_TYPE = :3`

	if err := db.GetContext(ctx, &uploadFile, query, file.Jno, file.WorkDate, file.FileType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uploadFile, fmt.Errorf("GetUploadFile: %w", err)
		}
		return uploadFile, fmt.Errorf("GetUploadFile: %w", err)
	}
	return uploadFile, nil
}

// 업로드 파일 정보 저장
func (r *Repository) AddUploadFile(ctx context.Context, tx Execer, file entity.UploadFile) error {
	agent := utils.GetAgent()

	query := `
		INSERT INTO IRIS_UPLOADED_FILES(FILE_TYPE, FILE_PATH, FILE_NAME, UPLOAD_ROUND, WORK_DATE, JNO, REG_DATE, REG_USER, REG_UNO, REG_AGENT)
		VALUES(:1, :2, :3, :4, :5, :6, SYSDATE, :7, :8, :9)`

	_, err := tx.ExecContext(ctx, query, file.FileType, file.FilePath, file.FileName, file.UploadRound, file.WorkDate, file.Jno, file.RegUser, file.RegUno, agent)
	if err != nil {
		return fmt.Errorf("AddUploadFile: %w", err)
	}

	return nil
}
