package txutil

import (
	"context"
	"csm-api/store"
	"csm-api/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
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

func DeferTxx(tx *sqlx.Tx, err *error) {
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

func BeginTxWithMode(ctx context.Context, db store.Beginner, readOnly bool) (tx *sql.Tx, err error) {
	tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: readOnly})
	if err != nil {
		orig := err
		return nil, utils.CustomMessageErrorfDepth(2, "begin tx", orig)
	}

	if readOnly {
		if _, e := tx.Exec("SET TRANSACTION READ ONLY"); e != nil {
			_ = tx.Rollback()
			orig := e
			return nil, utils.CustomMessageErrorfDepth(2, "set transaction", orig)
		}
	}

	return
}
