package service

import (
	"context"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/utils"
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
