package fx

import (
	"fmt"

	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

// StubProvider hard-codes a USD/BRL rate for development.
type StubProvider struct {
	usdToBRL float64
}

func NewStubFXRateProvider() *StubProvider {
	return &StubProvider{usdToBRL: 5.10}
}

func (s *StubProvider) GetRate(from, to string) (float64, error) {
	if from == to {
		return 1.0, nil
	}
	switch {
	case from == portfolio.CurrencyUSD && to == portfolio.CurrencyBRL:
		return s.usdToBRL, nil
	case from == portfolio.CurrencyBRL && to == portfolio.CurrencyUSD:
		return 1.0 / s.usdToBRL, nil
	}
	return 0, fmt.Errorf("unsupported pair: %s/%s", from, to)
}
