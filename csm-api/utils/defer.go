package utils

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func DeferTx(tx *sql.Tx, err *error) {
	if r := recover(); r != nil {
		_ = tx.Rollback()
		*err = CustomMessageErrorfDepth(2, "panic", fmt.Errorf("%v", r))
		return
	}

	if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			*err = CustomMessageErrorfDepth(2, "rollback", rollbackErr)
		}
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			*err = CustomMessageErrorfDepth(2, "commit", commitErr)
		}
	}
}

func DeferTxx(tx *sqlx.Tx, err *error) {
	if r := recover(); r != nil {
		_ = tx.Rollback()
		*err = CustomMessageErrorfDepth(2, "panic", fmt.Errorf("%v", r))
		return
	}

	if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			*err = CustomMessageErrorfDepth(2, "rollback", rollbackErr)
		}
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			*err = CustomMessageErrorfDepth(2, "commit", commitErr)
		}
	}
}
