package dto

import "time"

type CreatePortfolioRequest struct {
	BaseCurrency string `json:"baseCurrency" binding:"required,oneof=BRL USD"`
}

type PortfolioResponse struct {
	ID           string             `json:"id"`
	BaseCurrency string             `json:"baseCurrency"`
	CreatedAt    time.Time          `json:"createdAt"`
	Positions    []PositionResponse `json:"positions"`
}

type AddPositionRequest struct {
	Symbol   string  `json:"symbol" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,oneof=BRL USD"`
}

type PositionResponse struct {
	ID       string  `json:"id"`
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

type MoneyResponse struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type ValuationResponse struct {
	TotalBRL     MoneyResponse `json:"totalBRL"`
	TotalUSD     MoneyResponse `json:"totalUSD"`
	PercentInBRL float64       `json:"percentInBRL"`
	PercentInUSD float64       `json:"percentInUSD"`
}

type PortfolioWithValuationResponse struct {
	PortfolioResponse
	Valuation ValuationResponse `json:"valuation"`
}
