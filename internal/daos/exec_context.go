package daos

import (
	"database/sql"
)

// ExecContext allows both *sql.DB and *sql.Tx to be used
type ExecContext interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
