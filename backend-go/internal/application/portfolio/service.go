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

type AddPositionInput struct {
	PortfolioID string
	Symbol      string
	Quantity    float64
	Price       float64
	Currency    string
}

type AddPosition struct {
	repo portfolio.Repository
}

func NewAddPosition(r portfolio.Repository) *AddPosition {
	return &AddPosition{repo: r}
}

func (uc *AddPosition) Execute(in AddPositionInput) (*portfolio.Position, error) {
	p, err := uc.repo.FindByID(in.PortfolioID)
	if err != nil {
		return nil, err
	}
	pos, err := portfolio.NewPosition(p.ID, in.Symbol, in.Quantity, in.Price, in.Currency)
	if err != nil {
		return nil, err
	}
	p.AddPosition(*pos)
	if err := uc.repo.Save(*p); err != nil {
		return nil, err
	}
	return pos, nil
}

type PortfolioView struct {
	Portfolio *portfolio.Portfolio
	Valuation portfolio.Valuation
}

type GetPortfolioWithValuation struct {
	repo   portfolio.Repository
	market portfolio.MarketDataProvider
	fx     portfolio.FXRateProvider
}

func NewGetPortfolioWithValuation(r portfolio.Repository, m portfolio.MarketDataProvider, f portfolio.FXRateProvider) *GetPortfolioWithValuation {
	return &GetPortfolioWithValuation{repo: r, market: m, fx: f}
}

func (uc *GetPortfolioWithValuation) Execute(id string) (*PortfolioView, error) {
	p, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	prices := make(map[string]portfolio.Money, len(p.Positions))
	for _, pos := range p.Positions {
		price, currency, err := uc.market.GetPrice(pos.Symbol)
		if err != nil {
			return nil, err
		}
		prices[pos.Symbol] = portfolio.NewMoney(price, currency)
	}
	rate, err := uc.fx.GetRate(portfolio.CurrencyUSD, portfolio.CurrencyBRL)
	if err != nil {
		return nil, err
	}
	return &PortfolioView{Portfolio: p, Valuation: portfolio.Valuate(p.Positions, prices, rate)}, nil
}
