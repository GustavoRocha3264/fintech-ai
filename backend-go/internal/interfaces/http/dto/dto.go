package dto

import (
	"time"

	appdashboard "github.com/fintech/cbpi/backend-go/internal/application/dashboard"
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	domainanalysis "github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	domainportfolio "github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	domainsnapshot "github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

type CreatePortfolioRequest struct {
	BaseCurrency string `json:"baseCurrency" binding:"required,oneof=BRL USD"`
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

type PortfolioResponse struct {
	ID           string             `json:"id"`
	BaseCurrency string             `json:"baseCurrency"`
	CreatedAt    time.Time          `json:"createdAt"`
	Positions    []PositionResponse `json:"positions"`
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
	Portfolio PortfolioResponse `json:"portfolio"`
	Valuation ValuationResponse `json:"valuation"`
}

type AnalysisResponse struct {
	ID                           string    `json:"id"`
	PortfolioID                  string    `json:"portfolioId"`
	CreatedAt                    time.Time `json:"createdAt"`
	TotalValueBRL                float64   `json:"totalValueBRL"`
	TotalValueUSD                float64   `json:"totalValueUSD"`
	BRLExposurePercent           float64   `json:"brlExposurePercent"`
	USDExposurePercent           float64   `json:"usdExposurePercent"`
	TopAssetConcentrationPercent float64   `json:"topAssetConcentrationPercent"`
	Insights                     []string  `json:"insights"`
}

type SnapshotResponse struct {
	ID            string    `json:"id"`
	PortfolioID   string    `json:"portfolioId"`
	Timestamp     time.Time `json:"timestamp"`
	TotalValueBRL float64   `json:"totalValueBRL"`
	TotalValueUSD float64   `json:"totalValueUSD"`
}

type FXResponse struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

type DashboardFXResponse struct {
	USDToBRL float64 `json:"usdToBRL"`
	BRLToUSD float64 `json:"brlToUSD"`
}

type DashboardResponse struct {
	Portfolio    PortfolioResponse   `json:"portfolio"`
	Valuation    ValuationResponse   `json:"valuation"`
	LatestReport *AnalysisResponse   `json:"latestReport"`
	Snapshots    []SnapshotResponse  `json:"snapshots"`
	FX           DashboardFXResponse `json:"fx"`
}

func NewPositionResponse(p domainportfolio.Position) PositionResponse {
	return PositionResponse{
		ID:       p.ID,
		Symbol:   p.Symbol,
		Quantity: p.Quantity,
		Price:    p.Price,
		Currency: p.Currency,
	}
}

func NewPortfolioResponse(p *domainportfolio.Portfolio) PortfolioResponse {
	positions := make([]PositionResponse, 0, len(p.Positions))
	for _, pos := range p.Positions {
		positions = append(positions, NewPositionResponse(pos))
	}

	return PortfolioResponse{
		ID:           p.ID,
		BaseCurrency: p.BaseCurrency,
		CreatedAt:    p.CreatedAt,
		Positions:    positions,
	}
}

func NewMoneyResponse(m domainportfolio.Money) MoneyResponse {
	return MoneyResponse{Amount: m.Amount, Currency: m.Currency}
}

func NewValuationResponse(v domainportfolio.Valuation) ValuationResponse {
	return ValuationResponse{
		TotalBRL:     NewMoneyResponse(v.TotalBRL),
		TotalUSD:     NewMoneyResponse(v.TotalUSD),
		PercentInBRL: v.PercentInBRL,
		PercentInUSD: v.PercentInUSD,
	}
}

func NewPortfolioWithValuationResponse(view *apportfolio.PortfolioView) PortfolioWithValuationResponse {
	return PortfolioWithValuationResponse{
		Portfolio: NewPortfolioResponse(view.Portfolio),
		Valuation: NewValuationResponse(view.Valuation),
	}
}

func NewAnalysisResponse(r *domainanalysis.AnalysisReport) AnalysisResponse {
	insights := make([]string, len(r.Insights))
	copy(insights, r.Insights)

	return AnalysisResponse{
		ID:                           r.ID,
		PortfolioID:                  r.PortfolioID,
		CreatedAt:                    r.CreatedAt,
		TotalValueBRL:                r.TotalValueBRL,
		TotalValueUSD:                r.TotalValueUSD,
		BRLExposurePercent:           r.BRLExposurePercent,
		USDExposurePercent:           r.USDExposurePercent,
		TopAssetConcentrationPercent: r.TopAssetConcentrationPercent,
		Insights:                     insights,
	}
}

func NewSnapshotResponse(s domainsnapshot.PortfolioSnapshot) SnapshotResponse {
	return SnapshotResponse{
		ID:            s.ID,
		PortfolioID:   s.PortfolioID,
		Timestamp:     s.Timestamp,
		TotalValueBRL: s.TotalValueBRL,
		TotalValueUSD: s.TotalValueUSD,
	}
}

func NewSnapshotResponses(items []domainsnapshot.PortfolioSnapshot) []SnapshotResponse {
	out := make([]SnapshotResponse, 0, len(items))
	for _, item := range items {
		out = append(out, NewSnapshotResponse(item))
	}
	return out
}

func NewFXResponse(from, to string, rate float64) FXResponse {
	return FXResponse{From: from, To: to, Rate: rate}
}

type MarketSymbolResponse struct {
	Ticker   string `json:"ticker"`
	Currency string `json:"currency"`
	Name     string `json:"name"`
}

type MarketQuoteResponse struct {
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

func NewDashboardResponse(d *appdashboard.Dashboard) DashboardResponse {
	var report *AnalysisResponse
	if d.LatestReport != nil {
		r := NewAnalysisResponse(d.LatestReport)
		report = &r
	}

	return DashboardResponse{
		Portfolio:    NewPortfolioResponse(d.Portfolio),
		Valuation:    NewValuationResponse(d.Valuation),
		LatestReport: report,
		Snapshots:    NewSnapshotResponses(d.Snapshots),
		FX: DashboardFXResponse{
			USDToBRL: d.USDToBRL,
			BRLToUSD: d.BRLToUSD,
		},
	}
}
