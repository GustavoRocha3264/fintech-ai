package portfolio

import (
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

type CreatePortfolio struct {
	repo portfolio.Repository
}

func NewCreatePortfolio(r portfolio.Repository) *CreatePortfolio {
	return &CreatePortfolio{repo: r}
}

func (uc *CreatePortfolio) Execute(baseCurrency string) (*portfolio.Portfolio, error) {
	p, err := portfolio.New(baseCurrency)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(*p); err != nil {
		return nil, err
	}
	return p, nil
}

type GetPortfolio struct {
	repo portfolio.Repository
}

func NewGetPortfolio(r portfolio.Repository) *GetPortfolio {
	return &GetPortfolio{repo: r}
}

func (uc *GetPortfolio) Execute(id string) (*portfolio.Portfolio, error) {
	return uc.repo.FindByID(id)
}
