package portfolio

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidPositionCurrency = errors.New("position currency must be BRL or USD")
	ErrInvalidQuantity         = errors.New("position quantity must be positive")
	ErrInvalidPrice            = errors.New("position price must be positive")
	ErrEmptySymbol             = errors.New("position symbol must not be empty")
)

type Position struct {
	ID          string
	PortfolioID string
	Symbol      string
	Quantity    float64
	Price       float64
	Currency    string
}

func NewPosition(portfolioID, symbol string, quantity, price float64, currency string) (*Position, error) {
	if symbol == "" {
		return nil, ErrEmptySymbol
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if price <= 0 {
		return nil, ErrInvalidPrice
	}
	if currency != CurrencyBRL && currency != CurrencyUSD {
		return nil, ErrInvalidPositionCurrency
	}
	return &Position{
		ID:          uuid.NewString(),
		PortfolioID: portfolioID,
		Symbol:      symbol,
		Quantity:    quantity,
		Price:       price,
		Currency:    currency,
	}, nil
}
