package market

import (
	"log"

	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

// Fallback wraps a primary MarketDataProvider and falls back to a secondary
// one on error. Mirrors the same pattern used for FX rate providers so the
// app stays resilient when the live market API is unreachable.
type Fallback struct {
	primary   portfolio.MarketDataProvider
	secondary portfolio.MarketDataProvider
}

func NewFallback(primary, secondary portfolio.MarketDataProvider) *Fallback {
	return &Fallback{primary: primary, secondary: secondary}
}

func (f *Fallback) GetPrice(symbol string) (float64, string, error) {
	price, currency, err := f.primary.GetPrice(symbol)
	if err == nil {
		return price, currency, nil
	}
	log.Printf("market primary failed for %s (%v), using fallback", symbol, err)
	return f.secondary.GetPrice(symbol)
}
