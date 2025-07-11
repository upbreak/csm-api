package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

func (r *Repository) GetCodeList(ctx context.Context, db Queryer, pCode string) (*entity.Codes, error) {
	list := entity.Codes{}

	query := `
				SELECT
					   t1.CODE,
				 	   t1.P_CODE,
					   t1.CODE_NM,
					   t1.CODE_COLOR
			      FROM IRIS_CODE_SET t1
			     WHERE P_CODE = :1
			       AND t1.IS_USE = 'Y'
			  ORDER BY t1."ORDER"`

	if err := db.SelectContext(ctx, &list, query, pCode); err != nil {
		return nil, fmt.Errorf("GetCodeList err: %w", err)
	}

	return &list, nil
}

// 코드트리 조회
func (r *Repository) GetCodeTree(ctx context.Context, db Queryer, pCode string) (*entity.Codes, error) {
	codes := entity.Codes{}

	query := fmt.Sprintf(`
			SELECT 
			    LEVEL, 
			    C.IDX, 
			    C.CODE,
			    C.P_CODE,
			    C.CODE_NM,
			    C.CODE_COLOR,
			    C.UDF_VAL_03,
			    C.UDF_VAL_04,
			    C.UDF_VAL_05,
			    C.UDF_VAL_06,
			    C.UDF_VAL_07,
			    C."ORDER" AS SORT_NO,
			    C.IS_USE,
			    C.ETC			    
			FROM IRIS_CODE_SET C
			WHERE DEL_YN = 'N'
			START WITH P_CODE = '%s'
			CONNECT BY PRIOR CODE = P_CODE
			ORDER SIBLINGS BY "ORDER" ASC
		`, pCode)

	if err := db.SelectContext(ctx, &codes, query); err != nil {
		return nil, fmt.Errorf("GetCodeTrees err: %w", err)
	}

	return &codes, nil
}

// func: 코드트리 수정 및 저장
func (r *Repository) MergeCode(ctx context.Context, tx Execer, code entity.Code) error {

	query := `
			MERGE INTO IRIS_CODE_SET C1
			USING (
				SELECT 
					:1 AS IDX, 
					:2 AS CODE,
					:3 AS P_CODE,
					:4 AS CODE_NM,
					:5 AS CODE_COLOR,
					:6 AS UDF_VAL_03,
					:7 AS UDF_VAL_04,
					:8 AS UDF_VAL_05,
					:9 AS UDF_VAL_06,
					:10 AS UDF_VAL_07,
					:11 AS SORT_NO,
					:12 AS IS_USE,
					:13 AS ETC,	
					:14 AS UNO,	
					:15 AS USER_NAME
				FROM DUAL
			) C2
			ON (
				C1.DEL_YN = 'N' 
				AND C2.IDX IS NOT NULL
				AND C1.IDX = C2.IDX
			) WHEN MATCHED THEN
				UPDATE SET
					C1.CODE = C2.CODE,
					C1.CODE_NM = C2.CODE_NM,
					C1.CODE_COLOR = C2.CODE_COLOR,
					C1.UDF_VAL_03 = C2.UDF_VAL_03,
					C1.UDF_VAL_04 = C2.UDF_VAL_04,
					C1.UDF_VAL_05 = C2.UDF_VAL_05,
					C1.UDF_VAL_06 = C2.UDF_VAL_06,
					C1.UDF_VAL_07 = C2.UDF_VAL_07,
					C1."ORDER" = C2.SORT_NO,
					C1.IS_USE = C2.IS_USE,
					C1.ETC = C2.ETC,
					C1.MOD_UNO = C2.UNO,
					C1.MOD_USER = C2.USER_NAME,
					C1.MOD_DATE = SYSDATE
			WHEN NOT MATCHED THEN
				INSERT ( IDX, CODE, P_CODE, CODE_NM, CODE_COLOR, UDF_VAL_03, UDF_VAL_04, UDF_VAL_05, UDF_VAL_06, UDF_VAL_07, "ORDER", IS_USE, DEL_YN, ETC, REG_UNO, REG_USER, REG_DATE )
 				VALUES (
					SEQ_IRIS_CODE_SET.NEXTVAL,
					C2.CODE,
					C2.P_CODE,
					C2.CODE_NM,
					C2.CODE_COLOR,
					C2.UDF_VAL_03,
					C2.UDF_VAL_04,
					C2.UDF_VAL_05,
					C2.UDF_VAL_06,
					C2.UDF_VAL_07,
					C2.SORT_NO,
					C2.IS_USE,
					'N',
					C2.ETC,
					C2.UNO,
					C2.USER_NAME,
					SYSDATE)
				`

	if _, err := tx.ExecContext(ctx, query,
		code.IDX, code.Code, code.PCode, code.CodeNm, code.CodeColor,
		code.UdfVal03, code.UdfVal04, code.UdfVal05, code.UdfVal06, code.UdfVal07,
		code.SortNo, code.IsUse, code.Etc, code.RegUno, code.RegUser); err != nil {

		return fmt.Errorf("MergeCode err: %w", err)
	}

	return nil
}

// func: 코드 삭제
func (r *Repository) RemoveCode(ctx context.Context, tx Execer, idx int64) error {
	query := `
		UPDATE IRIS_CODE_SET
		SET 
			DEL_YN = 'Y'
		WHERE 
			IDX = :1			
	`

	if _, err := tx.ExecContext(ctx, query, idx); err != nil {
		return fmt.Errorf("store_code/RemoveCode err: %w", err)
	}

	return nil
}

// func: 코드순서 변경
func (r *Repository) ModifySortNo(ctx context.Context, tx Execer, codeSort entity.CodeSort) error {
	query := `
		UPDATE 
		    IRIS_CODE_SET
		SET
			"ORDER" = :1
		WHERE
			IDX = :2                 
		`
	if _, err := tx.ExecContext(ctx, query, codeSort.SortNo, codeSort.IDX); err != nil {
		return fmt.Errorf("store_code/ModifyCodeSort err: %w", err)
	}
	return nil
}

// func: 코드 중복 검사
func (r *Repository) DuplicateCheckCode(ctx context.Context, db Queryer, code string) (int, error) {
	var count int

	query := `
		SELECT 
		    COUNT(*) 
		FROM
		    IRIS_CODE_SET
		WHERE
		    DEL_YN = 'N'
			AND CODE = :1
		`

	if err := db.GetContext(ctx, &count, query, code); err != nil {
		return -1, fmt.Errorf("store_code/DuplicateCheckCode err: %w", err)
	}

	return count, nil
}
