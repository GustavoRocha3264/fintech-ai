package uow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fintech/cbpi/backend-go/internal/application/uow"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence/postgres"
)

// PostgresUoW maps the UnitOfWork boundary onto a real SQL transaction.
//
//	tx, _ := db.BeginTx(ctx, nil)
//	repos := { tx-scoped postgres repos }
//	if fn(repos) errors → tx.Rollback()
//	else                → tx.Commit()
//
// Repositories are constructed per-call from the *sql.Tx so every Save inside
// the callback rides on the same transaction. Reads done before the callback
// returns also see the in-flight writes (READ COMMITTED at minimum; the
// caller can override the isolation via context-bound options if needed).
//
// Note: this adapter uses only stdlib database/sql — bring your own driver
// (pgx/lib/pq) at the composition root via blank import.
type PostgresUoW struct {
	db   *sql.DB
	opts *sql.TxOptions
}

func NewPostgresUoW(db *sql.DB, opts *sql.TxOptions) *PostgresUoW {
	return &PostgresUoW{db: db, opts: opts}
}

func (u *PostgresUoW) Do(ctx context.Context, fn func(uow.Repositories) error) error {
	tx, err := u.db.BeginTx(ctx, u.opts)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	repos := uow.Repositories{
		Portfolio: postgres.NewPortfolioRepo(tx),
		Analysis:  postgres.NewAnalysisRepo(tx),
		Snapshot:  postgres.NewSnapshotRepo(tx),
	}

	if err := fn(repos); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			return fmt.Errorf("%w (rollback also failed: %v)", err, rbErr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
