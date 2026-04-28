package portfolio

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	CurrencyBRL = "BRL"
	CurrencyUSD = "USD"
)

var ErrInvalidBaseCurrency = errors.New("base currency must be BRL or USD")

type Portfolio struct {
	ID           string
	BaseCurrency string
	CreatedAt    time.Time
}

func New(baseCurrency string) (*Portfolio, error) {
	if baseCurrency != CurrencyBRL && baseCurrency != CurrencyUSD {
		return nil, ErrInvalidBaseCurrency
	}
	return &Portfolio{
		ID:           uuid.NewString(),
		BaseCurrency: baseCurrency,
		CreatedAt:    time.Now().UTC(),
	}, nil
}
