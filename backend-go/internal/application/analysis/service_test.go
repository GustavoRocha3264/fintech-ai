package analysis_test

import (
	"context"
	"errors"
	"testing"

	apanalysis "github.com/fintech/cbpi/backend-go/internal/application/analysis"
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/application/uow"
	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
	infrauow "github.com/fintech/cbpi/backend-go/internal/infrastructure/uow"
)

type fixture struct {
	port  *persistence.InMemoryPortfolioRepository
	an    *persistence.InMemoryAnalysisRepository
	snaps *persistence.InMemorySnapshotRepository
	uow   *infrauow.InMemoryUoW
	val   apportfolio.ValuationService
}

func newFixture() fixture {
	port := persistence.NewInMemoryPortfolioRepository()
	an := persistence.NewInMemoryAnalysisRepository()
	snaps := persistence.NewInMemorySnapshotRepository()
	return fixture{
		port:  port,
		an:    an,
		snaps: snaps,
		uow:   infrauow.NewInMemoryUoW(port, an, snaps),
		val:   apportfolio.NewValuationService(market.NewStubMarketDataProvider(), fx.NewStubFXRateProvider()),
	}
}

func TestRunAnalysis_GeneratesReportAndStoresIt(t *testing.T) {
	f := newFixture()
	create := apportfolio.NewCreatePortfolio(f.port)
	add := apportfolio.NewAddPosition(f.port)
	run := apanalysis.NewRunAnalysis(f.uow, f.val)
	latest := apanalysis.NewGetLatestAnalysis(f.an)

	p, err := create.Execute("USD")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := add.Execute(apportfolio.AddPositionInput{
		PortfolioID: p.ID, Symbol: "AAPL", Quantity: 100, Price: 195, Currency: "USD",
	}); err != nil {
		t.Fatalf("add: %v", err)
	}

	report, err := run.Execute(context.Background(), p.ID)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.PortfolioID != p.ID {
		t.Fatalf("portfolio id mismatch")
	}

	got, err := latest.Execute(p.ID)
	if err != nil || got.ID != report.ID {
		t.Fatalf("latest mismatch: %v / %v", got, err)
	}
	hist, _ := f.snaps.FindByPortfolioID(p.ID)
	if len(hist) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(hist))
	}
}

func TestRunAnalysis_PortfolioNotFound(t *testing.T) {
	f := newFixture()
	run := apanalysis.NewRunAnalysis(f.uow, f.val)

	_, err := run.Execute(context.Background(), "missing")
	if !errors.Is(err, portfolio.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetLatestAnalysis_NotFound(t *testing.T) {
	f := newFixture()
	latest := apanalysis.NewGetLatestAnalysis(f.an)

	_, err := latest.Execute("nope")
	if !errors.Is(err, analysis.ErrNotFound) {
		t.Fatalf("expected analysis.ErrNotFound, got %v", err)
	}
}

// TestRunAnalysis_AtomicRollback exercises the UoW boundary directly: a
// callback that saves an analysis report and then returns an error must leave
// the analysis store untouched.
func TestRunAnalysis_AtomicRollback(t *testing.T) {
	f := newFixture()
	create := apportfolio.NewCreatePortfolio(f.port)
	p, _ := create.Execute("USD")

	// Run a UoW that saves a report then deliberately fails.
	wantErr := errors.New("boom")
	err := f.uow.Do(context.Background(), func(repos uow.Repositories) error {
		// Save something inside the UoW, then fail. The save must be
		// reverted when the callback returns the error.
		_ = repos.Analysis.Save(analysis.AnalysisReport{ID: "tmp", PortfolioID: p.ID})
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected boom, got %v", err)
	}
	if _, err := f.an.GetLatestByPortfolioID(p.ID); !errors.Is(err, analysis.ErrNotFound) {
		t.Fatalf("expected nothing persisted, got %v", err)
	}
}
