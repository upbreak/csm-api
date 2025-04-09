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
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("GetCodeList err: %w", err)
	}

	return &list, nil
}

// 코드트리 조회
func (r *Repository) GetCodeTree(ctx context.Context, db Queryer) (*entity.Codes, error) {
	codes := entity.Codes{}

	query := `
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
			START WITH P_CODE IS NULL
			CONNECT BY PRIOR CODE = P_CODE
			ORDER SIBLINGS BY "ORDER" ASC
		`

	if err := db.SelectContext(ctx, &codes, query); err != nil {
		//TODO: 에러 아카이브
		return nil, fmt.Errorf("GetCodeTrees err: %w", err)
	}
	return &codes, nil
}
