package service

import (
	"context"
	"csm-api/ctxutil"
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
		return list, utils.CustomErrorf(err)
	}
	return list, nil
}

// 업로드 파일 정보
func (s *ServiceUploadFile) GetUploadFile(ctx context.Context, file entity.UploadFile) (entity.UploadFile, error) {
	data, err := s.Store.GetUploadFile(ctx, s.DB, file)
	if err != nil {
		return entity.UploadFile{}, utils.CustomErrorf(err)
	}
	return data, nil
}

// 업로드 파일 정보 저장
func (s *ServiceUploadFile) AddUploadFile(ctx context.Context, file entity.UploadFile) (err error) {
	tx, ok := ctxutil.GetTx(ctx)
	if !ok || tx == nil {
		tx, err = s.TDB.BeginTxx(ctx, nil)
		if err != nil {
			return utils.CustomErrorf(err)
		}

		defer func() {
			if r := recover(); r != nil {
				_ = tx.Rollback()
				err = utils.CustomMessageErrorf("panic", fmt.Errorf("%v", r))
				return
			}
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					err = utils.CustomMessageErrorf("rollback", rollbackErr)
				}
			} else {
				if commitErr := tx.Commit(); commitErr != nil {
					err = utils.CustomMessageErrorf("commit", commitErr)
				}
			}
		}()
	}

	// 차수
	uploadRound, err := s.Store.GetUploadRound(ctx, s.DB, file)
	if err != nil {
		return utils.CustomErrorf(err)
	}
	file.UploadRound = utils.ParseNullInt(strconv.Itoa(uploadRound))

	// 저장
	if err = s.Store.AddUploadFile(ctx, tx, file); err != nil {
		return utils.CustomErrorf(err)
	}

	return
}
