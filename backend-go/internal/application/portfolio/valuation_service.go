package portfolio

import (
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

// ValuationResult bundles everything a downstream use case might need from a
// single valuation pass: the aggregated totals/percentages, the per-symbol
// prices that fed it, and the FX rate used. Returning all three in one place
// avoids re-fetching prices when the caller wants to compute additional
// metrics (e.g. analysis concentration).
type ValuationResult struct {
	Valuation portfolio.Valuation
	Prices    map[string]portfolio.Money
	USDToBRL  float64
}

// ValuationService is the single seam for "give me a portfolio's current
// value." Every use case that needs valuation depends on this interface so we
// don't duplicate the price-fetching / FX-fetching / Valuate orchestration.
type ValuationService interface {
	Calculate(p portfolio.Portfolio) (ValuationResult, error)
}

type valuationService struct {
	market portfolio.MarketDataProvider
	fx     portfolio.FXRateProvider
}

func NewValuationService(m portfolio.MarketDataProvider, f portfolio.FXRateProvider) *valuationService {
	return &valuationService{market: m, fx: f}
}

func (s *valuationService) Calculate(p portfolio.Portfolio) (ValuationResult, error) {
	prices := make(map[string]portfolio.Money, len(p.Positions))
	for _, pos := range p.Positions {
		price, currency, err := s.market.GetPrice(pos.Symbol)
		if err != nil {
			return ValuationResult{}, err
		}
		prices[pos.Symbol] = portfolio.NewMoney(price, currency)
	}
	rate, err := s.fx.GetRate(portfolio.CurrencyUSD, portfolio.CurrencyBRL)
	if err != nil {
		return ValuationResult{}, err
	}
	return ValuationResult{
		Valuation: portfolio.Valuate(p.Positions, prices, rate),
		Prices:    prices,
		USDToBRL:  rate,
	}, nil
}
