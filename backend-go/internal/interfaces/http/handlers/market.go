package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	domainportfolio "github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	inframarket "github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/dto"
)

type MarketHandler struct {
	provider domainportfolio.MarketDataProvider
}

func NewMarketHandler(provider domainportfolio.MarketDataProvider) *MarketHandler {
	return &MarketHandler{provider: provider}
}

// GetSymbols returns the curated catalog of known symbols with their native
// currency. Used by the frontend to populate the symbol combobox.
func (h *MarketHandler) GetSymbols(c *gin.Context) {
	resp := make([]dto.MarketSymbolResponse, len(inframarket.KnownSymbols))
	for i, s := range inframarket.KnownSymbols {
		resp[i] = dto.MarketSymbolResponse{Ticker: s.Ticker, Currency: s.Currency, Name: s.Name}
	}
	c.JSON(http.StatusOK, resp)
}

// GetQuote resolves the current price and native currency for a single symbol
// via the configured MarketDataProvider (live or stub fallback).
func (h *MarketHandler) GetQuote(c *gin.Context) {
	symbol := strings.ToUpper(strings.TrimSpace(c.Param("symbol")))
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}
	price, currency, err := h.provider.GetPrice(symbol)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.MarketQuoteResponse{Symbol: symbol, Price: price, Currency: currency})
}
