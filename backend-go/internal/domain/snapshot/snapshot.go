package snapshot

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioSnapshot struct {
	ID            string
	PortfolioID   string
	Timestamp     time.Time
	TotalValueBRL float64
	TotalValueUSD float64
}

func New(portfolioID string, totalBRL, totalUSD float64) *PortfolioSnapshot {
	return &PortfolioSnapshot{
		ID:            uuid.NewString(),
		PortfolioID:   portfolioID,
		Timestamp:     time.Now().UTC(),
		TotalValueBRL: totalBRL,
		TotalValueUSD: totalUSD,
	}
}
