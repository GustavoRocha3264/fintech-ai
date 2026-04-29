package market

import "github.com/fintech/cbpi/backend-go/internal/domain/portfolio"

// StubProvider returns deterministic mock prices. Symbols ending in well-known
// BR tickers are priced in BRL; everything else defaults to USD.
type StubProvider struct {
	prices map[string]quote
}

type quote struct {
	price    float64
	currency string
}

func NewStubMarketDataProvider() *StubProvider {
	return &StubProvider{
		prices: map[string]quote{
			"PETR4": {price: 38.50, currency: portfolio.CurrencyBRL},
			"VALE3": {price: 65.20, currency: portfolio.CurrencyBRL},
			"ITUB4": {price: 32.10, currency: portfolio.CurrencyBRL},
			"AAPL":  {price: 195.40, currency: portfolio.CurrencyUSD},
			"MSFT":  {price: 425.10, currency: portfolio.CurrencyUSD},
			"VOO":   {price: 510.00, currency: portfolio.CurrencyUSD},
		},
	}
}

func (s *StubProvider) GetPrice(symbol string) (float64, string, error) {
	if q, ok := s.prices[symbol]; ok {
		return q.price, q.currency, nil
	}
	return 100.0, portfolio.CurrencyUSD, nil
}
