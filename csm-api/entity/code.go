package entity

import (
	"fmt"
	"github.com/guregu/null"
)

type Code struct {
	Level     null.Int    `json:"level" db:"LEVEL"`
	IDX       null.Int    `json:"idx" db:"IDX"`
	Code      null.String `json:"code" db:"CODE"`
	PCode     null.String `json:"p_code" db:"P_CODE"`
	CodeNm    null.String `json:"code_nm" db:"CODE_NM"`
	CodeColor null.String `json:"code_color" db:"CODE_COLOR"`
	UdfVal03  null.String `json:"udf_val_03" db:"UDF_VAL_03"`
	UdfVal04  null.String `json:"udf_val_04" db:"UDF_VAL_04"`
	UdfVal05  null.String `json:"udf_val_05" db:"UDF_VAL_05"`
	UdfVal06  null.String `json:"udf_val_06" db:"UDF_VAL_06"`
	UdfVal07  null.String `json:"udf_val_07" db:"UDF_VAL_07"`
	SortNo    null.Int    `json:"sort_no" db:"SORT_NO"`
	IsUse     null.String `json:"is_use" db:"IS_USE"`
	Etc       null.String `json:"etc" db:"ETC"`
	Base
}

type Codes []*Code

type CodeSort struct {
	IDX    null.Int `json:"idx" db:"IDX"`
	SortNo null.Int `json:"sort_no" db:"SORT_NO"`
}

type CodeSorts []*CodeSort

type CodeTree struct {
	IDX      null.Int    `json:"idx" db:"IDX"` // level이 1이 아니면 쌓아야함.
	Code     null.String `json:"code" db:"CODE"`
	Level    null.Int    `json:"level" db:"LEVEL"`
	PCode    null.String `json:"p_code" db:"P_CODE"`
	Expand   null.Bool   `json:"expand" db:"EXPAND"`
	Children *CodeTrees  `json:"code_trees" db:"CODE_TREES"`
	CodeSet  *Code       `json:"code_set" db:"CODE_SET"`
}

type CodeTrees []*CodeTree

// 코드를 코드트리 구조로 변환
func ConvertCodesToCodeTree(codes Codes, pCode string) (codeTrees CodeTrees, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("entity/ConvertCodesToCodeTree: %v", r)
		}
	}()
	if len(codes) == 0 {
		return CodeTrees{}, nil
	}

	for i := 0; i < len(codes); i++ {

		// 부모가 다른 경우 넘기기
		if codes[i].PCode.Valid && codes[i].PCode.String != pCode {
			continue
		}

		// code의 값 tree 형태로 변환
		codeTree := &CodeTree{}

		codeTree.IDX = codes[i].IDX
		codeTree.Level = codes[i].Level
		codeTree.Code = codes[i].Code

		// root는 root값 입력하기
		if codes[i].PCode.String == "" {
			codeTree.PCode.String = "root"
			codeTree.PCode.Valid = true
		} else {
			codeTree.PCode = codes[i].PCode
		}

		// 하위 code 넣기
		codeTree.Children = &CodeTrees{}
		child, err := ConvertCodesToCodeTree(codes[i+1:], codes[i].Code.String) // 현재 코드가 다음 레벨의 부모코드
		if err != nil {
			return codeTrees, err
		}
		if len(child) > 0 {
			codeTree.Children = &child
			codeTree.Expand.Bool = true
			codeTree.Expand.Valid = true
		} else {
			codeTree.Expand.Bool = false
			codeTree.Expand.Valid = true
		}
		codeTree.CodeSet = codes[i]
		codeTrees = append(codeTrees, codeTree)
	}

	return
}
