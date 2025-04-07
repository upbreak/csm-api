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
