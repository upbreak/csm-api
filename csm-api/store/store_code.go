package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

func (r *Repository) GetCodeList(ctx context.Context, db Queryer, pCode string) (*entity.CodeSqls, error) {
	sqls := entity.CodeSqls{}

	query := `SELECT
				t1.CODE,
				t1.P_CODE,
				t1.CODE_NM 
			FROM IRIS_CODE_SET t1
			WHERE P_CODE = :1`

	if err := db.SelectContext(ctx, &sqls, query, pCode); err != nil {
		return nil, fmt.Errorf("GetCodeList err: %w", err)
	}

	return &sqls, nil
}
