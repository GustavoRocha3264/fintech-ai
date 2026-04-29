package fx

import (
	"log"

	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

// Fallback wraps a primary provider and falls back to a secondary one on
// error. Used to keep the app resilient when the live FX API is unreachable.
type Fallback struct {
	primary   portfolio.FXRateProvider
	secondary portfolio.FXRateProvider
}

func NewFallback(primary, secondary portfolio.FXRateProvider) *Fallback {
	return &Fallback{primary: primary, secondary: secondary}
}

func (f *Fallback) GetRate(from, to string) (float64, error) {
	r, err := f.primary.GetRate(from, to)
	if err == nil {
		return r, nil
	}
	log.Printf("fx primary failed (%v), using fallback", err)
	return f.secondary.GetRate(from, to)
}
