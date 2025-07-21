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

func BeginTxWithCleanMode(ctx context.Context, db store.Beginner, readOnly bool) (tx *sql.Tx, cleanup func(), err error) {
	//conn, err := db.Conn(ctx)
	//if err != nil {
	//	return nil, nil, utils.CustomMessageErrorfDepth(2, "db.Conn", err)
	//}
	//
	//tx, err = conn.BeginTx(ctx, &sql.TxOptions{ReadOnly: readOnly})
	tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: readOnly})
	if err != nil {
		orig := err
		//_ = conn.Close()
		return nil, nil, utils.CustomMessageErrorfDepth(2, "begin tx", orig)
	}

	if readOnly {
		if _, e := tx.Exec("SET TRANSACTION READ ONLY"); e != nil {
			_ = tx.Rollback()
			orig := e
			//_ = conn.Close()
			return nil, nil, utils.CustomMessageErrorfDepth(2, "set transaction", orig)
		}
	}

	cleanup = func() {
		//if closeErr := conn.Close(); closeErr != nil {
		//	if err != nil {
		//		err = fmt.Errorf("%v; cleanup conn.Close: %w", err, closeErr)
		//	} else {
		//		err = utils.CustomMessageErrorfDepth(2, "cleanup conn.Close", closeErr)
		//	}
		//}
	}

	return
}
