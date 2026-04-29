package snapshot

import (
	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

type GetHistory struct {
	repo snapshot.Repository
}

func NewGetHistory(r snapshot.Repository) *GetHistory {
	return &GetHistory{repo: r}
}

func (uc *GetHistory) Execute(portfolioID string) ([]snapshot.PortfolioSnapshot, error) {
	return uc.repo.FindByPortfolioID(portfolioID)
}
