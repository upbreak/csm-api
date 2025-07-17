package txutil

import (
	"context"
	"csm-api/store"
	"csm-api/utils"
	"database/sql"
	"errors"
	"fmt"
)

func DeferTx(tx *sql.Tx, err *error) {
	if tx == nil {
		return
	}

	if r := recover(); r != nil {
		_ = tx.Rollback()
		*err = utils.CustomMessageErrorfDepth(2, "panic", fmt.Errorf("%v", r))
		return
	}

	if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			*err = utils.CustomMessageErrorfDepth(2, "rollback", rollbackErr)
		}
	} else {
		if commitErr := tx.Commit(); commitErr != nil && !errors.Is(commitErr, sql.ErrTxDone) {
			*err = utils.CustomMessageErrorfDepth(2, "commit", commitErr)
		}
	}
}

func BeginTxWithMode(ctx context.Context, db store.Beginner, readOnly bool) (*sql.Tx, error) {
	opts := &sql.TxOptions{
		ReadOnly: false,
	}

	conn, err := db.Conn(ctx)
	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		return nil, utils.CustomMessageErrorfDepth(2, "begin tx", err)
	}

	if readOnly {
		if _, err = tx.Exec("SET TRANSACTION READ ONLY"); err != nil {
			_ = tx.Rollback()
			return nil, utils.CustomMessageErrorfDepth(2, "set transaction", err)
		}
	}

	return tx, nil
}
