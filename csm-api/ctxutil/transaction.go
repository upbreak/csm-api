package ctxutil

import (
	"context"
	"csm-api/utils"
	"github.com/jmoiron/sqlx"
)

type ctxTxKey struct{}

var txKey = ctxTxKey{}

func WithTx(ctx context.Context, db *sqlx.DB) (context.Context, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, txKey, tx)
	return ctx, nil
}

func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sqlx.Tx)
	return tx, ok
}

func DeferTx(ctx context.Context, handler string, errRef *error) func() {
	tx, ok := GetTx(ctx)
	if !ok || tx == nil {
		return func() {} // no-op
	}

	return func() {
		if *errRef != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				*errRef = utils.CustomMessageErrorfDepth(2, "failed to rollback transaction", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				*errRef = utils.CustomMessageErrorfDepth(2, "failed to commit transaction", commitErr)
			}
		}
	}
}
