// Package postgres implements the domain repositories on top of a SQL
// transaction. Each repo takes a *sql.Tx so the same connection / transaction
// flows through every write inside a UnitOfWork.Do.
//
// This package uses only stdlib database/sql. Wire your driver of choice at
// the composition root, e.g.:
//
//	import _ "github.com/lib/pq"
//	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
package postgres

import (
	"context"
	"database/sql"
)

// dbtx is the subset of *sql.Tx the repositories use. Pulled into a tiny
// interface so unit tests can stub it without a live database.
type dbtx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
