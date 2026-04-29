// Package uow contains concrete UnitOfWork implementations.
//
// InMemoryUoW provides real atomic semantics over the in-memory repositories
// by snapshotting their state before the callback runs and restoring it on
// failure. This keeps tests fast and offline while still exercising the same
// transaction boundary the production Postgres adapter uses.
package uow

import (
	"context"
	"sync"

	"github.com/fintech/cbpi/backend-go/internal/application/uow"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
)

type InMemoryUoW struct {
	portfolios *persistence.InMemoryPortfolioRepository
	analyses   *persistence.InMemoryAnalysisRepository
	snapshots  *persistence.InMemorySnapshotRepository

	// One transaction at a time. The in-memory store has no per-row locks, so
	// serializing is simpler and sufficient for the use cases we have.
	mu sync.Mutex
}

func NewInMemoryUoW(
	p *persistence.InMemoryPortfolioRepository,
	a *persistence.InMemoryAnalysisRepository,
	s *persistence.InMemorySnapshotRepository,
) *InMemoryUoW {
	return &InMemoryUoW{portfolios: p, analyses: a, snapshots: s}
}

func (u *InMemoryUoW) Do(ctx context.Context, fn func(uow.Repositories) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	pBackup := u.portfolios.Snapshot()
	aBackup := u.analyses.Snapshot()
	sBackup := u.snapshots.Snapshot()

	repos := uow.Repositories{
		Portfolio: u.portfolios,
		Analysis:  u.analyses,
		Snapshot:  u.snapshots,
	}

	err := fn(repos)
	if err != nil {
		u.portfolios.Restore(pBackup)
		u.analyses.Restore(aBackup)
		u.snapshots.Restore(sBackup)
		return err
	}
	return nil
}
