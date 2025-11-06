package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	sqlMap      = make(map[string]*sqlx.DB)
	sqlMutex    = sync.RWMutex{}
	DefaultName = "default"
)

func InitFor(dbName string, dsn string) error {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	sqlMutex.Lock()
	sqlMap[dbName] = db
	sqlMutex.Unlock()
	return nil
}

func Init(dsn string) error {
	return InitFor(DefaultName, dsn)
}

func NamedQueryFor(dbName, query string, arg interface{}) (*sqlx.Rows, error) {
	sqlMutex.RLock()
	sqlxdb, exists := sqlMap[dbName]
	sqlMutex.RUnlock()
	if !exists {
		return nil, fmt.Errorf("database not initialized")
	}
	return sqlxdb.NamedQuery(query, arg)
}

func NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return NamedQueryFor(DefaultName, query, arg)
}

func NamedExecFor(dbName, query string, arg interface{}) (sql.Result, error) {
	sqlMutex.RLock()
	sqlxdb, exists := sqlMap[dbName]
	sqlMutex.RUnlock()
	if !exists {
		return nil, fmt.Errorf("database not initialized")
	}
	return sqlxdb.NamedExec(query, arg)
}
func NamedExec(query string, arg interface{}) (sql.Result, error) {
	return NamedExecFor(DefaultName, query, arg)
}

func QueryRowFor(dbName, query string, args ...any) *sql.Row {

	sqlMutex.RLock()
	sqlxdb, exists := sqlMap[dbName]
	sqlMutex.RUnlock()
	if !exists {
		return nil
	}
	return sqlxdb.QueryRow(query, args...)
}
func QueryRow(query string, args ...any) *sql.Row {
	return QueryRowFor(DefaultName, query, args...)
}

func SelectFor(dbName string, dest interface{}, query string, args ...interface{}) error {

	sqlMutex.RLock()
	sqlxdb, exists := sqlMap[dbName]
	sqlMutex.RUnlock()
	if !exists {
		return fmt.Errorf("database not initialized")
	}
	return sqlxdb.Select(dest, query, args...)
}
func Select(dest interface{}, query string, args ...interface{}) error {
	return SelectFor(DefaultName, dest, query, args...)
}

func ExecContextFor(ctx context.Context, dbName string, query string, args ...any) (sql.Result, error) {
	sqlMutex.RLock()
	sqlxdb, exists := sqlMap[dbName]
	sqlMutex.RUnlock()
	if !exists {
		return nil, fmt.Errorf("database not initialized")
	}
	return sqlxdb.ExecContext(ctx, query, args...)
}
func ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return ExecContextFor(ctx, DefaultName, query, args...)
}

func NewTranscationFor(dbName string) (*SqlxTransaction, error) {
	sqlMutex.RLock()
	sqlxdb, exists := sqlMap[dbName]
	sqlMutex.RUnlock()
	if !exists {
		return nil, fmt.Errorf("database not initialized")
	}
	tx, err := sqlxdb.Beginx()
	if err != nil {
		return nil, fmt.Errorf("beginx: %w", err)
	}
	return NewSqlxTransaction(tx), nil
}
func NewTranscation() (*SqlxTransaction, error) {
	return NewTranscationFor(DefaultName)
}

type SqlxTransaction struct {
	tx *sqlx.Tx
}

func NewSqlxTransaction(tx *sqlx.Tx) *SqlxTransaction {
	return &SqlxTransaction{tx: tx}
}

func (t *SqlxTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *SqlxTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *SqlxTransaction) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return t.tx.NamedQuery(query, arg)
}
