package portfolio

import "errors"

var ErrNotFound = errors.New("portfolio not found")

type Repository interface {
	Save(p Portfolio) error
	FindByID(id string) (*Portfolio, error)
}

type MarketDataProvider interface {
	GetPrice(symbol string) (price float64, currency string, err error)
}

type FXRateProvider interface {
	// GetRate returns the rate to convert 1 unit of `from` into `to`.
	GetRate(from, to string) (float64, error)
}
