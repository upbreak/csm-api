package store

import (
	"context"
	"csm-api/clock"
	"csm-api/config"
	"csm-api/utils"
	"database/sql"
	"fmt"
	"github.com/godror/godror"
	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
	"time"
)

// New creates a new database connection
func New(ctx context.Context, cfg *config.DBConfig) (*sqlx.DB, func(), error) {
	// godror.ConnectionParams 설정
	var P godror.ConnectionParams
	P.Username = cfg.UserName
	P.Password = godror.NewPassword(cfg.Password)
	P.ConnectString = fmt.Sprintf("%s:%s/%s", cfg.Host, cfg.Port, cfg.OracleSid)
	//P.ConnectString = fmt.Sprintf("%s:%s/%s?_enableTxReadWrite=0", cfg.Host, cfg.Port, cfg.OracleSid)
	P.StandaloneConnection = sql.NullBool{Bool: true, Valid: true}
	//P.Timezone = time.FixedZone("Asia/Seoul", 9*60*60) // 애플리케이션 타임존 설정
	//P.SetSessionParamOnInit("TIME_ZONE", "Asia/Seoul") // 세션 타임존 설정

	// OCI 세션 풀링 활성화 (godror SessionPool)
	P.PoolParams.MinSessions = 2                   // 최소 세션 2개
	P.PoolParams.MaxSessions = 10                  // 최대 세션 10개
	P.PoolParams.SessionIncrement = 2              // 풀 확장할 때마다 +2
	P.PoolParams.SessionTimeout = 60 * time.Second // 유휴 풀 세션 TTL
	P.PoolParams.WaitTimeout = 10 * time.Second    // 풀 대기 최대 시간

	// 디버깅 용도로 DSN 출력
	fmt.Printf("DSN: %s\n", P.StringWithPassword())

	// Connector 생성
	connector := godror.NewConnector(P)

	// sql.DB 생성
	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(0)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)

	// 연결 확인
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, utils.CustomErrorf(err)
	}

	// sqlx.DB 생성
	xdb := sqlx.NewDb(db, "godror")

	cleanup := func() {
		fmt.Printf("close db: %s\n", cfg.DBName)
		_ = db.Close()
	}

	return xdb, cleanup, nil
}

type Repository struct {
	Clocker clock.Clocker
}

type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Conn(ctx context.Context) (*sql.Conn, error)
}

type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type NamedExecer interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
	Rebind(query string) string
}

var (
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)
