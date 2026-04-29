package dto

import (
	"encoding/json"
	"testing"
	"time"

	appdashboard "github.com/fintech/cbpi/backend-go/internal/application/dashboard"
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	domainanalysis "github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	domainportfolio "github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	domainsnapshot "github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

func TestNewPortfolioResponse_MapsDomainWithoutLeakingFields(t *testing.T) {
	createdAt := time.Date(2026, 4, 29, 10, 0, 0, 0, time.UTC)
	p := &domainportfolio.Portfolio{
		ID:           "portfolio-1",
		BaseCurrency: domainportfolio.CurrencyUSD,
		CreatedAt:    createdAt,
		Positions: []domainportfolio.Position{
			{
				ID:          "position-1",
				PortfolioID: "portfolio-1",
				Symbol:      "AAPL",
				Quantity:    10,
				Price:       195,
				Currency:    domainportfolio.CurrencyUSD,
			},
		},
	}

	got := NewPortfolioResponse(p)
	if got.ID != p.ID {
		t.Fatalf("expected id %q, got %q", p.ID, got.ID)
	}
	if got.BaseCurrency != p.BaseCurrency {
		t.Fatalf("expected base currency %q, got %q", p.BaseCurrency, got.BaseCurrency)
	}
	if got.CreatedAt != createdAt {
		t.Fatalf("expected createdAt %v, got %v", createdAt, got.CreatedAt)
	}
	if len(got.Positions) != 1 {
		t.Fatalf("expected 1 position, got %d", len(got.Positions))
	}

	raw, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal portfolio response: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal portfolio response: %v", err)
	}

	if _, exists := payload["portfolioId"]; exists {
		t.Fatalf("portfolio response leaked internal position relationship field")
	}
	if _, exists := payload["name"]; exists {
		t.Fatalf("portfolio response leaked unexpected field")
	}
}

func TestNewDashboardResponse_JSONContract(t *testing.T) {
	now := time.Date(2026, 4, 29, 12, 0, 0, 0, time.UTC)
	dashboard := &appdashboard.Dashboard{
		Portfolio: &domainportfolio.Portfolio{
			ID:           "portfolio-1",
			BaseCurrency: domainportfolio.CurrencyBRL,
			CreatedAt:    now,
			Positions: []domainportfolio.Position{
				{ID: "position-1", Symbol: "PETR4", Quantity: 20, Price: 32, Currency: domainportfolio.CurrencyBRL},
			},
		},
		Valuation: domainportfolio.Valuation{
			TotalBRL:     domainportfolio.NewMoney(640, domainportfolio.CurrencyBRL),
			TotalUSD:     domainportfolio.NewMoney(128, domainportfolio.CurrencyUSD),
			PercentInBRL: 80,
			PercentInUSD: 20,
		},
		LatestReport: &domainanalysis.AnalysisReport{
			ID:                           "analysis-1",
			PortfolioID:                  "portfolio-1",
			CreatedAt:                    now,
			TotalValueBRL:                640,
			TotalValueUSD:                128,
			BRLExposurePercent:           80,
			USDExposurePercent:           20,
			TopAssetConcentrationPercent: 60,
			Insights:                     []string{"Concentrated portfolio"},
		},
		Snapshots: []domainsnapshot.PortfolioSnapshot{
			{ID: "snap-1", PortfolioID: "portfolio-1", Timestamp: now, TotalValueBRL: 640, TotalValueUSD: 128},
		},
		USDToBRL: 5,
		BRLToUSD: 0.2,
	}

	response := NewDashboardResponse(dashboard)
	if response.Portfolio.ID != dashboard.Portfolio.ID {
		t.Fatalf("expected portfolio id %q, got %q", dashboard.Portfolio.ID, response.Portfolio.ID)
	}
	if response.LatestReport == nil {
		t.Fatal("expected latest report to be present")
	}

	raw, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal dashboard response: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal dashboard response: %v", err)
	}

	if _, ok := payload["portfolio"]; !ok {
		t.Fatal("expected portfolio field in dashboard response")
	}
	if _, ok := payload["valuation"]; !ok {
		t.Fatal("expected valuation field in dashboard response")
	}
	if _, ok := payload["latestReport"]; !ok {
		t.Fatal("expected latestReport field in dashboard response")
	}
	if _, ok := payload["snapshots"]; !ok {
		t.Fatal("expected snapshots field in dashboard response")
	}
	if _, ok := payload["fx"]; !ok {
		t.Fatal("expected fx field in dashboard response")
	}
}

func TestNewPortfolioWithValuationResponse_UsesNestedTransportShape(t *testing.T) {
	view := &apportfolio.PortfolioView{
		Portfolio: &domainportfolio.Portfolio{
			ID:           "portfolio-2",
			BaseCurrency: domainportfolio.CurrencyUSD,
			CreatedAt:    time.Date(2026, 4, 29, 14, 0, 0, 0, time.UTC),
		},
		Valuation: domainportfolio.Valuation{
			TotalBRL:     domainportfolio.NewMoney(500, domainportfolio.CurrencyBRL),
			TotalUSD:     domainportfolio.NewMoney(100, domainportfolio.CurrencyUSD),
			PercentInBRL: 50,
			PercentInUSD: 50,
		},
	}

	got := NewPortfolioWithValuationResponse(view)
	if got.Portfolio.ID != view.Portfolio.ID {
		t.Fatalf("expected portfolio id %q, got %q", view.Portfolio.ID, got.Portfolio.ID)
	}
	if got.Valuation.TotalUSD.Amount != view.Valuation.TotalUSD.Amount {
		t.Fatalf("expected total USD %v, got %v", view.Valuation.TotalUSD.Amount, got.Valuation.TotalUSD.Amount)
	}
}
