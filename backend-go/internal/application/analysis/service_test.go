package analysis_test

import (
	"errors"
	"testing"

	apanalysis "github.com/fintech/cbpi/backend-go/internal/application/analysis"
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
)

func TestRunAnalysis_GeneratesReportAndStoresIt(t *testing.T) {
	portRepo := persistence.NewInMemoryPortfolioRepository()
	analysisRepo := persistence.NewInMemoryAnalysisRepository()
	snapshotRepo := persistence.NewInMemorySnapshotRepository()
	marketProv := market.NewStubMarketDataProvider()
	fxProv := fx.NewStubFXRateProvider()

	create := apportfolio.NewCreatePortfolio(portRepo)
	add := apportfolio.NewAddPosition(portRepo)
	run := apanalysis.NewRunAnalysis(portRepo, analysisRepo, snapshotRepo, marketProv, fxProv)
	latest := apanalysis.NewGetLatestAnalysis(analysisRepo)

	p, err := create.Execute("USD")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := add.Execute(apportfolio.AddPositionInput{
		PortfolioID: p.ID, Symbol: "AAPL", Quantity: 100, Price: 195, Currency: "USD",
	}); err != nil {
		t.Fatalf("add: %v", err)
	}
	if _, err := add.Execute(apportfolio.AddPositionInput{
		PortfolioID: p.ID, Symbol: "PETR4", Quantity: 10, Price: 38, Currency: "BRL",
	}); err != nil {
		t.Fatalf("add: %v", err)
	}

	report, err := run.Execute(p.ID)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.PortfolioID != p.ID {
		t.Fatalf("portfolio id mismatch")
	}
	if report.TotalValueBRL <= 0 || report.TotalValueUSD <= 0 {
		t.Fatalf("expected positive totals: %+v", report)
	}
	if report.TopAssetConcentrationPercent <= 50 {
		t.Fatalf("expected high concentration, got %v", report.TopAssetConcentrationPercent)
	}

	got, err := latest.Execute(p.ID)
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	if got.ID != report.ID {
		t.Fatalf("expected stored report id %s, got %s", report.ID, got.ID)
	}
}

func TestRunAnalysis_PortfolioNotFound(t *testing.T) {
	portRepo := persistence.NewInMemoryPortfolioRepository()
	analysisRepo := persistence.NewInMemoryAnalysisRepository()
	snapshotRepo := persistence.NewInMemorySnapshotRepository()
	run := apanalysis.NewRunAnalysis(portRepo, analysisRepo, snapshotRepo, market.NewStubMarketDataProvider(), fx.NewStubFXRateProvider())

	_, err := run.Execute("missing")
	if !errors.Is(err, portfolio.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetLatestAnalysis_NotFound(t *testing.T) {
	analysisRepo := persistence.NewInMemoryAnalysisRepository()
	latest := apanalysis.NewGetLatestAnalysis(analysisRepo)

	_, err := latest.Execute("nope")
	if !errors.Is(err, analysis.ErrNotFound) {
		t.Fatalf("expected analysis.ErrNotFound, got %v", err)
	}
}
