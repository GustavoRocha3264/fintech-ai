package persistence

import (
	"sort"
	"sync"

	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

type InMemorySnapshotRepository struct {
	mu       sync.RWMutex
	byPort   map[string][]snapshot.PortfolioSnapshot
}

func NewInMemorySnapshotRepository() *InMemorySnapshotRepository {
	return &InMemorySnapshotRepository{byPort: map[string][]snapshot.PortfolioSnapshot{}}
}

func (r *InMemorySnapshotRepository) Save(s snapshot.PortfolioSnapshot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byPort[s.PortfolioID] = append(r.byPort[s.PortfolioID], s)
	return nil
}

// FindByPortfolioID returns snapshots ordered oldest-first, suitable for
// charting. A copy is returned so callers can't mutate the underlying slice.
func (r *InMemorySnapshotRepository) FindByPortfolioID(portfolioID string) ([]snapshot.PortfolioSnapshot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	src := r.byPort[portfolioID]
	out := make([]snapshot.PortfolioSnapshot, len(src))
	copy(out, src)
	sort.Slice(out, func(i, j int) bool { return out[i].Timestamp.Before(out[j].Timestamp) })
	return out, nil
}
