package market

import "github.com/fintech/cbpi/backend-go/internal/domain/portfolio"

// Symbol describes a tradeable asset the platform knows about by default.
// The list is intentionally small and curated — Brapi handles arbitrary symbols
// at query time; this catalog only exists to populate the frontend's suggestions.
type Symbol struct {
	Ticker   string
	Currency string
	Name     string
}

// KnownSymbols is the curated list surfaced via GET /api/v1/market/symbols.
// It intentionally mirrors the stub's price map so both modes have identical
// defaults.
var KnownSymbols = []Symbol{
	{Ticker: "PETR4", Currency: portfolio.CurrencyBRL, Name: "Petrobras PN"},
	{Ticker: "VALE3", Currency: portfolio.CurrencyBRL, Name: "Vale ON"},
	{Ticker: "ITUB4", Currency: portfolio.CurrencyBRL, Name: "Itaú Unibanco PN"},
	{Ticker: "AAPL", Currency: portfolio.CurrencyUSD, Name: "Apple Inc."},
	{Ticker: "MSFT", Currency: portfolio.CurrencyUSD, Name: "Microsoft Corp."},
	{Ticker: "VOO", Currency: portfolio.CurrencyUSD, Name: "Vanguard S&P 500 ETF"},
}
