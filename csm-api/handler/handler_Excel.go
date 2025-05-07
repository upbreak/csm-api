package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-05-07
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @description: 엑셀 import, export
 */
type HandlerExcel struct {
	Service service.ExcelService
}

// TODO: 임시작성
// func: 일간 퇴직공제 export
// @param
// -
func (h *HandlerExcel) ExportDailyDeduction(w http.ResponseWriter, r *http.Request) {
	var rows []entity.DailyDeduction

	if err := json.NewDecoder(r.Body).Decode(&rows); err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	f, err := h.Service.ExportDailyDeduction(rows)
	if err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	// 파일 스트림 전송 (성공한 경우)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=retirement_deduction_%s.xlsx", time.Now().Format("20060102")))
	w.Header().Set("File-Name", fmt.Sprintf("retirement_deduction_%s.xlsx", time.Now().Format("20060102")))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, File-Name")

	if err = f.Write(w); err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("엑셀 파일 전송 실패: %v", err))
		return
	}
}

// TODO: 임시작성
// func: 퇴직공제 엑셀 import
// @param
// -
func (h *HandlerExcel) ImportDeduction(w http.ResponseWriter, r *http.Request) {
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

	// 임시 파일로 저장
	tempFilePath := "./excel_deduction/temp_" + header.Filename
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

	_, err = io.Copy(outFile, file)
	if err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("failed to save the uploaded file: %v", err))
		return
	}

	err = h.Service.ImportDeduction(tempFilePath)
	if err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("failed to parse Excel file: %v", err))
		return
	}

	SuccessResponse(r.Context(), w)
}
