package db

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Tx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
}

type SQLTx struct {
	Tx *sqlx.Tx
}

func (s *SQLTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.Tx.Exec(query, args...)
}

func (s *SQLTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.Tx.Query(query, args...)
}

func (s *SQLTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.Tx.QueryRow(query, args...)
}

func (s *SQLTx) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return s.Tx.NamedExec(query, arg)
}

func (s *SQLTx) Commit() error {
	return s.Tx.Commit()
}

func (s *SQLTx) Rollback() error {
	return s.Tx.Rollback()
}
