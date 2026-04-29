package persistence

import (
	"sync"

	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
)

type InMemoryAnalysisRepository struct {
	mu      sync.RWMutex
	latest  map[string]analysis.AnalysisReport
}

func NewInMemoryAnalysisRepository() *InMemoryAnalysisRepository {
	return &InMemoryAnalysisRepository{latest: map[string]analysis.AnalysisReport{}}
}

func (r *InMemoryAnalysisRepository) Save(report analysis.AnalysisReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	prev, ok := r.latest[report.PortfolioID]
	if !ok || report.CreatedAt.After(prev.CreatedAt) {
		r.latest[report.PortfolioID] = report
	}
	return nil
}

func (r *InMemoryAnalysisRepository) GetLatestByPortfolioID(portfolioID string) (*analysis.AnalysisReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rep, ok := r.latest[portfolioID]
	if !ok {
		return nil, analysis.ErrNotFound
	}
	return &rep, nil
}

func (r *InMemoryAnalysisRepository) Snapshot() map[string]analysis.AnalysisReport {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]analysis.AnalysisReport, len(r.latest))
	for k, v := range r.latest {
		insights := make([]string, len(v.Insights))
		copy(insights, v.Insights)
		v.Insights = insights
		out[k] = v
	}
	return out
}

func (r *InMemoryAnalysisRepository) Restore(snap map[string]analysis.AnalysisReport) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latest = snap
}
