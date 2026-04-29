package portfolio_test

import (
	"math"
	"testing"

	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func TestValuate_MixedCurrencies(t *testing.T) {
	positions := []portfolio.Position{
		{Symbol: "PETR4", Quantity: 100, Currency: portfolio.CurrencyBRL},
		{Symbol: "AAPL", Quantity: 10, Currency: portfolio.CurrencyUSD},
	}
	prices := map[string]portfolio.Money{
		"PETR4": portfolio.NewMoney(40.0, portfolio.CurrencyBRL), // 4000 BRL
		"AAPL":  portfolio.NewMoney(200.0, portfolio.CurrencyUSD), // 2000 USD = 10000 BRL
	}
	v := portfolio.Valuate(positions, prices, 5.0)

	if !approx(v.TotalBRL.Amount, 14000.0) {
		t.Fatalf("total BRL = %v", v.TotalBRL.Amount)
	}
	if !approx(v.TotalUSD.Amount, 4000.0/5.0+2000.0) {
		t.Fatalf("total USD = %v", v.TotalUSD.Amount)
	}
	if !approx(v.PercentInBRL+v.PercentInUSD, 100.0) {
		t.Fatalf("percentages don't sum to 100: %v + %v", v.PercentInBRL, v.PercentInUSD)
	}
	// 4000 BRL native vs 10000 BRL-equivalent → ~28.57% / ~71.43%
	if !approx(v.PercentInBRL, 4000.0/14000.0*100) {
		t.Fatalf("percent BRL = %v", v.PercentInBRL)
	}
}

func TestValuate_EmptyPortfolio(t *testing.T) {
	v := portfolio.Valuate(nil, nil, 5.0)
	if v.TotalBRL.Amount != 0 || v.TotalUSD.Amount != 0 {
		t.Fatalf("expected zero totals")
	}
}
