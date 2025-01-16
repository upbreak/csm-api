package store

import (
	"context"
	"csm-api/clock"
	"csm-api/config"
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
)

func New(ctx context.Context, cfg *config.DBConfig) (*sqlx.DB, func(), error) {
	db, err := sql.Open("godror", fmt.Sprintf("%s/%s@%s:%s/%s", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.OracleSid))
	if err != nil {
		return nil, func() {}, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}

	xdb := sqlx.NewDb(db, "godror")
	cleanup := func() {
		fmt.Println("close db: %s", cfg.DBName)
		_ = db.Close()
	}
	return xdb, cleanup, err
}

type Repository struct {
	Clocker clock.Clocker
}

type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

var (
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)
