package handler

import (
	"csm-api/entity"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"net/http"
	"time"
)

type HandlerExcel struct {
}

func (h *HandlerExcel) DailyDeduction(w http.ResponseWriter, r *http.Request) {
	var rows []entity.DailyDeduction

	if err := json.NewDecoder(r.Body).Decode(&rows); err != nil {
		FailResponse(r.Context(), w, err)
		return
	}

	f := excelize.NewFile()
	sheetName := "Sheet1"

	// 헤더 스타일
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#000000",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#B7DEE8"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// 헤더 구조 정의
	mainHeaders := []struct {
		Title     string
		StartCell string
		EndCell   string
	}{
		{"No.", "B2", "B3"}, {"일자", "C2", "C3"}, {"현장", "D2", "D3"},
		{"공제가입번호", "E2", "E3"}, {"업체명", "F2", "F3"}, {"소속업체", "G2", "G3"},
		{"성명(한국명)", "H2", "H3"}, {"생년월일", "I2", "I3"}, {"휴대전화번호", "J2", "J3"},
		{"직종", "K2", "K3"}, {"퇴직공제", "L2", "L3"}, {"비대상사유", "M2", "M3"},
		{"카드발급", "N2", "N3"}, {"성별", "O2", "O3"},
		{"태그내역", "P2", "Q2"}, {"근무시간", "R2", "S2"}, {"자동집계 일수", "T2", "U2"},
		{"인증방식", "V2", "V3"}, {"단말기번호", "W2", "W3"},
	}

	for _, h := range mainHeaders {
		err := f.MergeCell(sheetName, h.StartCell, h.EndCell)
		if err != nil {
			FailResponse(r.Context(), w, err)
			return
		}
		err = f.SetCellValue(sheetName, h.StartCell, h.Title)
		if err != nil {
			FailResponse(r.Context(), w, err)
			return
		}
		err = f.SetCellStyle(sheetName, h.StartCell, h.EndCell, headerStyle)
		if err != nil {
			FailResponse(r.Context(), w, err)
			return
		}
	}

	subHeaders := []string{
		"출근시간", "퇴근시간", "출근시간", "퇴근시간", "최초(태그기준)", "최종(수정기준)",
	}
	subCols := []string{"P", "Q", "R", "S", "T", "U"}

	for i, text := range subHeaders {
		cell := fmt.Sprintf("%s3", subCols[i])
		err := f.SetCellValue(sheetName, cell, text)
		if err != nil {
			FailResponse(r.Context(), w, err)
			return
		}
		err = f.SetCellStyle(sheetName, cell, cell, headerStyle)
		if err != nil {
			FailResponse(r.Context(), w, err)
			return
		}
	}

	// 본문 데이터
	for i, row := range rows {
		values := []interface{}{
			i + 1, row.Value1, row.Value2, row.Value3, row.Value4, row.Value5,
			row.Value6, row.Value7, row.Value8, row.Value9, row.Value10,
			row.Value11, row.Value12, row.Value13, row.Value14, row.Value15,
			row.Value16, row.Value17, row.Value18, row.Value19, row.Value20,
			row.Value21,
		}

		for j, value := range values {
			colIndex := j + 2                                        // B=2부터 시작
			cell, _ := excelize.CoordinatesToCellName(colIndex, i+4) // 행은 4부터 시작
			err := f.SetCellValue(sheetName, cell, value)
			if err != nil {
				FailResponse(r.Context(), w, err)
				return
			}

			// 가운데 정렬 적용할 열 인덱스 (1부터 시작)
			centerAlignedCols := map[string]bool{
				"B": true, // No.
				"I": true, // 생년월일
				"J": true, // 휴대전화번호
				"L": true, // 퇴직공제
				"N": true, // 카드발급
				"O": true, // 성별
				"V": true, // 인증방식
				"W": true, // 단말기번호
				"P": true, // 출근시간
				"Q": true, // 퇴근시간
				"R": true, // 출근시간
				"S": true, // 퇴근시간
				"T": true, // 최초(태그기준)
				"U": true, // 최종(수정기준)
			}

			colName, _ := excelize.ColumnNumberToName(colIndex)
			if centerAlignedCols[colName] {
				err := f.SetCellStyle(sheetName, cell, cell, centerStyle)
				if err != nil {
					return
				}
			} else {
				err := f.SetCellStyle(sheetName, cell, cell, dataStyle)
				if err != nil {
					return
				}
			}
		}
	}

	// 열 너비 자동 조정 (B~W)
	for col := 'B'; col <= 'W'; col++ {
		colName := string(col)
		width := 15
		if colName == "B" { // No.
			width = 10
		} else if colName == "D" { // 현장
			width = 45
		}
		err := f.SetColWidth(sheetName, colName, colName, float64(width))
		if err != nil {
			FailResponse(r.Context(), w, err)
			return
		}
	}

	// 파일 스트림 전송 (성공한 경우)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=retirement_deduction_%s.xlsx", time.Now().Format("20060102")))
	w.Header().Set("File-Name", fmt.Sprintf("retirement_deduction_%s.xlsx", time.Now().Format("20060102")))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, File-Name")

	if err := f.Write(w); err != nil {
		FailResponse(r.Context(), w, fmt.Errorf("엑셀 파일 전송 실패: %v", err))
		return
	}
}
