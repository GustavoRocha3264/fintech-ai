// Package uow defines the Unit of Work boundary used by application services
// that need multi-repository atomicity. It lives in the application layer so
// use cases depend on an interface, not on any specific transaction
// mechanism (in-memory snapshot, sql.Tx, etc.). Concrete implementations live
// in /internal/infrastructure/uow.
package uow

import (
	"context"

	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

// Repositories is the transaction-scoped bundle handed to each UoW callback.
// Inside Do(...) every save goes through these — never through the
// "outside-the-uow" repos — so a failure rolls everything back together.
type Repositories struct {
	Portfolio portfolio.Repository
	Analysis  analysis.Repository
	Snapshot  snapshot.Repository
}

// UnitOfWork wraps a unit of work in a transaction. The callback runs against
// transaction-scoped repositories. If the callback returns an error, every
// write inside it is reverted; if it returns nil, the writes are committed
// atomically.
type UnitOfWork interface {
	Do(ctx context.Context, fn func(repos Repositories) error) error
}
