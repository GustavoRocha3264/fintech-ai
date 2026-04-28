package dto

import "time"

type CreatePortfolioRequest struct {
	BaseCurrency string `json:"baseCurrency" binding:"required,oneof=BRL USD"`
}

type PortfolioResponse struct {
	ID           string    `json:"id"`
	BaseCurrency string    `json:"baseCurrency"`
	CreatedAt    time.Time `json:"createdAt"`
}
