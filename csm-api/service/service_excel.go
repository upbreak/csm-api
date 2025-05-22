package service

import (
	"csm-api/entity"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

type ServiceExcel struct{}

// TODO: 임시작성
// func: 일간 퇴직공제 export
// @param
// -
func (s *ServiceExcel) ExportDailyDeduction(rows []entity.DailyDeduction) (*excelize.File, error) {
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
			return nil, err
		}
		err = f.SetCellValue(sheetName, h.StartCell, h.Title)
		if err != nil {
			return nil, err
		}
		err = f.SetCellStyle(sheetName, h.StartCell, h.EndCell, headerStyle)
		if err != nil {
			return nil, err
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
			return nil, err
		}
		err = f.SetCellStyle(sheetName, cell, cell, headerStyle)
		if err != nil {
			return nil, err
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
				return nil, err
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
				err = f.SetCellStyle(sheetName, cell, cell, centerStyle)
				if err != nil {
					return nil, err
				}
			} else {
				err = f.SetCellStyle(sheetName, cell, cell, dataStyle)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// 열 너비 조정 (B~W)
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
			return nil, err
		}
	}

	return f, nil
}

// TODO: 임시작성
// func: 퇴직공제 엑셀 import
// @param
// -
func (s *ServiceExcel) ImportDeduction(path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	// 첫 번째 시트명 가져오기
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return fmt.Errorf("no sheet found in Excel file")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	var deductions []entity.DailyDeduction

	for rowNumber := 4; rowNumber < len(rows)+1; rowNumber++ {
		values := make([]string, 21)

		for colIdx := 0; colIdx < 21; colIdx++ {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+3, rowNumber)
			val, _ := f.GetCellValue(sheetName, cellName)
			values[colIdx] = val
		}

		Value18str := values[17]
		Value19str := values[18]
		Value21str := values[20]

		var Value18 int
		if Value18str != "" {
			Value18, err = strconv.Atoi(Value18str)
			if err != nil {
				Value18 = 0
			}
		}
		var Value19 int
		if Value19str != "" {
			Value19, err = strconv.Atoi(Value19str)
			if err != nil {
				Value19 = 0
			}
		}
		var Value21 int
		if Value21str != "" {
			Value21, err = strconv.Atoi(Value21str)
			if err != nil {
				Value21 = 0
			}
		}

		deduction := entity.DailyDeduction{
			Value1:  values[0],
			Value2:  values[1],
			Value3:  values[2],
			Value4:  values[3],
			Value5:  values[4],
			Value6:  values[5],
			Value7:  values[6],
			Value8:  values[7],
			Value9:  values[8],
			Value10: values[9],
			Value11: values[10],
			Value12: values[11],
			Value13: values[12],
			Value14: values[13],
			Value15: values[14],
			Value16: values[15],
			Value17: values[16],
			Value18: Value18,
			Value19: Value19,
			Value20: values[19],
			Value21: Value21,
		}
		deductions = append(deductions, deduction)
	}

	return nil
}

// 작업허가서 import
func (s *ServiceExcel) ImportWorkLetter(path string) (int64, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open Excel file: %w", err)
	}

	sheetName := f.GetSheetName(0)

	// '근로자' 열 찾기
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0, fmt.Errorf("failed to read sheet rows: %w", err)
	}

	if len(rows) < 4 {
		return 0, fmt.Errorf("row 4 does not exist in the sheet")
	}

	var targetCol string
	for colIdx, cell := range rows[3] {
		if cell == "근로자" {
			colLetter, _ := excelize.ColumnNumberToName(colIdx + 1)
			targetCol = colLetter
			break
		}
	}
	if targetCol == "" {
		return 0, fmt.Errorf("could not find column labeled '근로자' in row 4")
	}

	// 근로자열의 마지막 행 값 가져오기(총 근로자 수)
	rowNum := 5 // 마지막행을 찾기 위해 리스트 부분인 5행부터 시작
	var lastValue string
	for {
		cellRef := fmt.Sprintf("%s%d", targetCol, rowNum)
		val, err := f.GetCellValue(sheetName, cellRef)
		if err != nil {
			return 0, fmt.Errorf("failed to read cell %s: %w", cellRef, err)
		}
		if strings.TrimSpace(val) == "" {
			break
		}
		lastValue = val
		rowNum++
	}

	result, err := strconv.ParseInt(lastValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse value '%s' as int64: %w", lastValue, err)
	}
	return result, nil
}
