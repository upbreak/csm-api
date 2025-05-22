package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
	"fmt"
	"strconv"
)

type ServiceUploadFile struct {
	DB    store.Queryer
	TDB   store.Beginner
	Store store.UploadFileStore
}

// 업로드 파일 리스트
func (s *ServiceUploadFile) GetUploadFileList(ctx context.Context, file entity.UploadFile) ([]entity.UploadFile, error) {
	list, err := s.Store.GetUploadFileList(ctx, s.DB, file)
	if err != nil {
		return list, fmt.Errorf("serviceUploadFile.GetUploadFileList: %w", err)
	}
	return list, nil
}

// 업로드 파일 정보
func (s *ServiceUploadFile) GetUploadFile(ctx context.Context, file entity.UploadFile) (entity.UploadFile, error) {
	data, err := s.Store.GetUploadFile(ctx, s.DB, file)
	if err != nil {
		return entity.UploadFile{}, fmt.Errorf("serviceUploadFile.UploadFile: %w", err)
	}
	return data, nil
}

// 업로드 파일 정보 저장
func (s *ServiceUploadFile) AddUploadFile(ctx context.Context, file entity.UploadFile) (err error) {
	tx, err := s.TDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("serviceUploadFile.AddUploadFile: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("serviceUploadFile.AddUploadFile rollback error: %w", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("serviceUploadFile.AddUploadFile commit error: %w", commitErr)
			}
		}
	}()

	// 차수
	uploadRound, err := s.Store.GetUploadRound(ctx, s.DB, file)
	if err != nil {
		return fmt.Errorf("serviceUploadFile.AddUploadFile: %w", err)
	}
	file.UploadRound = utils.ParseNullInt(strconv.Itoa(uploadRound))

	// 저장
	if err = s.Store.AddUploadFile(ctx, tx, file); err != nil {
		return fmt.Errorf("serviceUploadFile.AddUploadFile: %w", err)
	}

	return nil
}
