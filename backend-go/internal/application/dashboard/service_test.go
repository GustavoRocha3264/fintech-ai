package dashboard_test

import (
	"errors"
	"testing"

	apanalysis "github.com/fintech/cbpi/backend-go/internal/application/analysis"
	"github.com/fintech/cbpi/backend-go/internal/application/dashboard"
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
)

type fixture struct {
	portRepo     *persistence.InMemoryPortfolioRepository
	analysisRepo *persistence.InMemoryAnalysisRepository
	snapshotRepo *persistence.InMemorySnapshotRepository
	fx           *fx.StubProvider
	val          apportfolio.ValuationService
}

func newFixture() fixture {
	portRepo := persistence.NewInMemoryPortfolioRepository()
	analysisRepo := persistence.NewInMemoryAnalysisRepository()
	snapshotRepo := persistence.NewInMemorySnapshotRepository()
	fxProv := fx.NewStubFXRateProvider()
	val := apportfolio.NewValuationService(market.NewStubMarketDataProvider(), fxProv)
	return fixture{portRepo, analysisRepo, snapshotRepo, fxProv, val}
}

func TestDashboard_FullyPopulatedAfterAnalysis(t *testing.T) {
	f := newFixture()
	create := apportfolio.NewCreatePortfolio(f.portRepo)
	add := apportfolio.NewAddPosition(f.portRepo)
	run := apanalysis.NewRunAnalysis(f.portRepo, f.analysisRepo, f.snapshotRepo, f.val)
	get := dashboard.NewGetDashboard(f.portRepo, f.val, f.analysisRepo, f.snapshotRepo, f.fx)

	p, err := create.Execute("USD")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := add.Execute(apportfolio.AddPositionInput{
		PortfolioID: p.ID, Symbol: "AAPL", Quantity: 5, Price: 195, Currency: "USD",
	}); err != nil {
		t.Fatalf("add: %v", err)
	}
	if _, err := run.Execute(p.ID); err != nil {
		t.Fatalf("run: %v", err)
	}

	d, err := get.Execute(p.ID)
	if err != nil {
		t.Fatalf("dashboard: %v", err)
	}
	if d.Portfolio == nil || d.Portfolio.ID != p.ID {
		t.Fatalf("portfolio missing: %+v", d.Portfolio)
	}
	if d.Valuation.TotalUSD.Amount <= 0 {
		t.Fatalf("expected positive USD total, got %+v", d.Valuation)
	}
	if d.LatestReport == nil {
		t.Fatalf("expected latest report after running analysis")
	}
	if len(d.Snapshots) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(d.Snapshots))
	}
	if d.USDToBRL <= 0 || d.BRLToUSD <= 0 {
		t.Fatalf("FX block not populated: %+v", d)
	}
}

func TestDashboard_NoAnalysisYet_ReportIsNil(t *testing.T) {
	f := newFixture()
	create := apportfolio.NewCreatePortfolio(f.portRepo)
	get := dashboard.NewGetDashboard(f.portRepo, f.val, f.analysisRepo, f.snapshotRepo, f.fx)

	p, err := create.Execute("USD")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	d, err := get.Execute(p.ID)
	if err != nil {
		t.Fatalf("dashboard: %v", err)
	}
	if d.LatestReport != nil {
		t.Fatalf("expected nil report, got %+v", d.LatestReport)
	}
	if len(d.Snapshots) != 0 {
		t.Fatalf("expected no snapshots, got %d", len(d.Snapshots))
	}
}

func TestDashboard_PortfolioNotFound(t *testing.T) {
	f := newFixture()
	get := dashboard.NewGetDashboard(f.portRepo, f.val, f.analysisRepo, f.snapshotRepo, f.fx)

	_, err := get.Execute("missing")
	if !errors.Is(err, portfolio.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
