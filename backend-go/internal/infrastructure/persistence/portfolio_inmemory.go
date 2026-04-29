package persistence

import (
	"sync"

	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

type InMemoryPortfolioRepository struct {
	mu    sync.RWMutex
	store map[string]portfolio.Portfolio
}

func NewInMemoryPortfolioRepository() *InMemoryPortfolioRepository {
	return &InMemoryPortfolioRepository{store: map[string]portfolio.Portfolio{}}
}

func (r *InMemoryPortfolioRepository) Save(p portfolio.Portfolio) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[p.ID] = p
	return nil
}

func (r *InMemoryPortfolioRepository) FindByID(id string) (*portfolio.Portfolio, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return nil, portfolio.ErrNotFound
	}
	return &p, nil
}

// Snapshot returns a deep copy of the current state. Used by the in-memory
// UoW to support rollback on transaction failure.
func (r *InMemoryPortfolioRepository) Snapshot() map[string]portfolio.Portfolio {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]portfolio.Portfolio, len(r.store))
	for k, v := range r.store {
		positions := make([]portfolio.Position, len(v.Positions))
		copy(positions, v.Positions)
		v.Positions = positions
		out[k] = v
	}
	return out
}

// Restore overwrites the current state with the supplied snapshot.
func (r *InMemoryPortfolioRepository) Restore(snap map[string]portfolio.Portfolio) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store = snap
}
