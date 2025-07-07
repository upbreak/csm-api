package handler

import (
	"csm-api/config"
	"csm-api/ctxutil"
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	Service     service.ExcelService
	FileService service.UploadFileService
	DB          *sqlx.DB
}

// excel 자료 import
// fileType: WORK_LETTER (작업허가서), TBM (TBM 문서), DEDUCTION (퇴직공제), REPORT (작업일보), ADD_DAILY_WORKER (현장 근로자 등록)
// POST ROW DATA
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
	// 현장
	snoString := r.FormValue("sno")
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
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	var dir string
	if addDir == "" {
		dir = filepath.Join(cfg.UploadPath, strings.ToLower(fileType), dates[0], dates[1], dates[2], jnoString)
	} else {
		dir = filepath.Join(cfg.UploadPath, strings.ToLower(fileType), dates[0], dates[1], dates[2], jnoString, addDir)
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

	// context 안에 트랜잭션 저장
	ctx, err := ctxutil.WithTx(r.Context(), h.DB)
	if err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("failed to begin transaction: %v", err))
		return
	}
	defer ctxutil.DeferTx(ctx, "ImportExcel", &err)()

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
	if err = h.FileService.AddUploadFile(ctx, uploadFile); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 엑셀 파싱 및 db 저장
	if fileType == "TBM" {
		tbm := entity.Tbm{
			Sno:        utils.ParseNullInt(snoString),
			Department: utils.ParseNullString(department),
			TbmDate:    utils.ParseNullTime(workDate),
			Base: entity.Base{
				RegUser: utils.ParseNullString(regUser),
				RegUno:  utils.ParseNullInt(regUno),
			},
		}
		if err = h.Service.ImportTbm(ctx, tempFilePath, tbm); err != nil {
			FailResponse(ctx, w, err)
			return
		}
	} else if fileType == "DEDUCTION" {
		deduction := entity.Deduction{
			Sno:        utils.ParseNullInt(snoString),
			RecordDate: utils.ParseNullTime(workDate),
			Base: entity.Base{
				RegUser: utils.ParseNullString(regUser),
				RegUno:  utils.ParseNullInt(regUno),
			},
		}
		if err = h.Service.ImportDeduction(ctx, tempFilePath, deduction); err != nil {
			FailResponse(ctx, w, err)
			return
		}
	} else if fileType == "ADD_DAILY_WORKER" {
		workDaily := entity.WorkerDaily{
			Sno:        utils.ParseNullInt(snoString),
			Jno:        utils.ParseNullInt(jnoString),
			RecordDate: utils.ParseNullTime(workDate),
			Base: entity.Base{
				RegUser: utils.ParseNullString(regUser),
				RegUno:  utils.ParseNullInt(regUno),
			},
		}
		if err = h.Service.ImportAddDailyWorker(ctx, tempFilePath, workDaily); err != nil {
			FailResponse(ctx, w, err)
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

// 현장 근로자 엑셀 양식 다운로드
func (h *HandlerExcel) DailyWorkerFormExport(w http.ResponseWriter, r *http.Request) {
	f := excelize.NewFile()
	sheet := "Sheet1"

	// 헤더
	headers := []string{"No.", "이름", "부서/조직명", "핸드폰번호", "근로날짜", "출근시간", "퇴근시간", "공수"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	// 스타일
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#C6EFCE"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	borderStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	dateFmt := "yyyy-mm-dd"
	timeFmt := "hh:mm"
	dateStyle, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &dateFmt,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	timeStyle, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &timeFmt,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	f.SetCellStyle(sheet, "A1", "H1", headerStyle)

	// 날짜 및 시간 데이터 파싱
	layoutDate := "2006-01-02"
	date1, _ := time.Parse(layoutDate, "2025-07-01")
	date2, _ := time.Parse(layoutDate, "2025-07-01")

	start1 := time.Date(1899, 12, 31, 7, 26, 0, 0, time.UTC)
	end1 := time.Date(1899, 12, 31, 15, 41, 0, 0, time.UTC)

	start2 := time.Date(1899, 12, 31, 11, 21, 0, 0, time.UTC)
	end2 := time.Date(1899, 12, 31, 14, 12, 0, 0, time.UTC)

	rows := [][]interface{}{
		{1, "홍길동1", "진웅종합건설", "010-1234-5678", date1, start1, end1, 1},
		{2, "홍길동2", "진웅종합건설", "010-1234-5678", date2, start2, end2, 0.5},
	}

	// 셀 값 입력
	for rowIdx, row := range rows {
		for colIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
		}
	}

	// 스타일 적용
	f.SetCellStyle(sheet, "A2", "H10", borderStyle)
	f.SetCellStyle(sheet, "E2", "E10", dateStyle)
	f.SetCellStyle(sheet, "F2", "F10", timeStyle)
	f.SetCellStyle(sheet, "G2", "G10", timeStyle)

	// 컬럼 너비
	f.SetColWidth(sheet, "A", "H", 15)

	// 파일 이름 생성
	fileName := "현장근로자 입력 양식.xlsx"
	encodedName := url.PathEscape(fileName)

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedName))
	w.Header().Set("File-Name", fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, File-Name")

	// 파일 출력
	if err := f.Write(w); err != nil {
		http.Error(w, "엑셀 파일 생성 실패", http.StatusInternalServerError)
		return
	}
}
