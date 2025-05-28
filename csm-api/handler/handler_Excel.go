package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-05-07
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @description: 엑셀 import, export
 */
type HandlerExcel struct {
	Service     service.ExcelService
	FileService service.UploadFileService
}

// excel 자료 import
// fileType: WORK_LETTER (작업허가서), TBM (TBM 문서), DEDUCTION (퇴직공제), REPORT (작업일보)
func (h *HandlerExcel) ImportExcel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 최대 10MB
	if err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("failed to parse multipart form: %v", err))
		return
	}

	// 파일 받기
	file, header, err := r.FormFile("file")
	if err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("failed to receive the file: %v", err))
		return
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			FailResponse(r.Context(), w, fmt.Errorf("failed to file Close: %v", err))
			return
		}
	}(file)

	// 엑셀 파일 확장자 검사
	if !(len(header.Filename) > 5 && (header.Filename[len(header.Filename)-5:] == ".xlsx" || header.Filename[len(header.Filename)-4:] == ".xls")) {
		FailResponse(r.Context(), w, fmt.Errorf("only Excel files (.xlsx, .xls) are allowed"))
		return
	}

	// 날짜 (2000-01-01)
	workDate := r.FormValue("work_date")
	// 프로젝트
	jnoString := r.FormValue("jno")
	// 회사
	department := r.FormValue("department")
	// 종류
	fileType := r.FormValue("file_type")
	// 추가 파일 경로
	addDir := r.FormValue("add_dir")
	if workDate == "" || fileType == "" || jnoString == "" {
		FailResponse(r.Context(), w, fmt.Errorf("missing 'file_date' or 'jno' or 'file_type' field"))
		return
	}
	regUser := r.FormValue("reg_user")
	regUno := r.FormValue("reg_uno")

	dates := strings.Split(workDate, "-")
	if len(dates) != 3 {
		FailResponse(r.Context(), w, fmt.Errorf("invalid 'file_date' format (expected: YYYY-MM-DD)"))
		return
	}

	// 저장 경로 설정 및 생성
	var dir string
	if addDir == "" {
		dir = filepath.Join("uploads", strings.ToLower(fileType), dates[0], dates[1], dates[2], jnoString)
	} else {
		dir = filepath.Join("uploads", strings.ToLower(fileType), dates[0], dates[1], dates[2], jnoString, addDir)
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("failed to create upload directory: %v", err))
		return
	}

	tempFilePath := filepath.Join(dir, header.Filename)
	outFile, err := os.Create(tempFilePath)
	if err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("failed to create a temporary file: %v", err))
		return
	}
	defer func(outFile *os.File) {
		err = outFile.Close()
		if err != nil {
			FailResponse(r.Context(), w, fmt.Errorf("failed to outFile Close: %v", err))
			return
		}
	}(outFile)

	// 파일 복사(저장)
	_, err = io.Copy(outFile, file)
	if err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("failed to save the uploaded file: %v", err))
		return
	}

	// 파일 정보 저장
	uploadFile := entity.UploadFile{
		FileType: utils.ParseNullString(fileType),
		FilePath: utils.ParseNullString(dir),
		FileName: utils.ParseNullString(header.Filename),
		WorkDate: utils.ParseNullTime(workDate),
		Jno:      utils.ParseNullInt(jnoString),
		Base: entity.Base{
			RegUser: utils.ParseNullString(regUser),
			RegUno:  utils.ParseNullInt(regUno),
		},
	}
	if err = h.FileService.AddUploadFile(r.Context(), uploadFile); err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	// 엑셀 파싱 및 db 저장
	if fileType == "TBM" {
		tbm := entity.Tbm{
			Jno:        utils.ParseNullInt(jnoString),
			Department: utils.ParseNullString(department),
			TbmDate:    utils.ParseNullTime(workDate),
			Base: entity.Base{
				RegUser: utils.ParseNullString(regUser),
				RegUno:  utils.ParseNullInt(regUno),
			},
		}
		if err = h.Service.ImportTbm(r.Context(), tempFilePath, tbm); err != nil {
			FailResponse(r.Context(), w, err)
			return
		}
	}

	SuccessResponse(r.Context(), w)
}

// excel 자료 export
func (h *HandlerExcel) ExportExcel(w http.ResponseWriter, r *http.Request) {
	jno := r.URL.Query().Get("jno")
	workDate := r.URL.Query().Get("work_date")
	fileType := r.URL.Query().Get("file_type")
	if workDate == "" || fileType == "" || jno == "" {
		FailResponse(r.Context(), w, fmt.Errorf("missing 'work_date' or 'file_type' or 'jno' field"))
		return
	}

	file := entity.UploadFile{
		FileType: utils.ParseNullString(fileType),
		Jno:      utils.ParseNullInt(jno),
		WorkDate: utils.ParseNullTime(workDate),
	}

	// 파일 경로, 명칭 조회
	data, err := h.FileService.GetUploadFile(r.Context(), file)
	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	filePath := filepath.Join(data.FilePath.String, data.FileName.String)

	// 파일 존재 확인
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		FailResponse(r.Context(), w, fmt.Errorf("file does not exist: %v", filePath))
		return
	}

	// 파일 열기
	f, err := os.Open(filePath)
	if err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("failed to file open: %v", err))
		return
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			FailResponse(r.Context(), w, fmt.Errorf("failed to close file: %v", err))
			return
		}
	}(f)

	// 다운로드용 헤더 설정
	fileName := data.FileName.String
	encodedName := url.PathEscape(fileName)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedName))
	w.Header().Set("File-Name", data.FileName.String)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, File-Name")

	// 파일 스트림 복사
	if _, err = io.Copy(w, f); err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("파일 전송 실패: %v", err))
		return
	}
}
