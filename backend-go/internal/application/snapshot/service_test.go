package snapshot_test

import (
	"context"
	"testing"

	apanalysis "github.com/fintech/cbpi/backend-go/internal/application/analysis"
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	apsnapshot "github.com/fintech/cbpi/backend-go/internal/application/snapshot"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
	infrauow "github.com/fintech/cbpi/backend-go/internal/infrastructure/uow"
)

func TestSnapshotsCapturedAfterAnalysis(t *testing.T) {
	portRepo := persistence.NewInMemoryPortfolioRepository()
	analysisRepo := persistence.NewInMemoryAnalysisRepository()
	snapshotRepo := persistence.NewInMemorySnapshotRepository()
	uow := infrauow.NewInMemoryUoW(portRepo, analysisRepo, snapshotRepo)
	val := apportfolio.NewValuationService(market.NewStubMarketDataProvider(), fx.NewStubFXRateProvider())

	create := apportfolio.NewCreatePortfolio(portRepo)
	add := apportfolio.NewAddPosition(portRepo)
	run := apanalysis.NewRunAnalysis(uow, val)
	history := apsnapshot.NewGetHistory(snapshotRepo)

	p, err := create.Execute("USD")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := add.Execute(apportfolio.AddPositionInput{
		PortfolioID: p.ID, Symbol: "AAPL", Quantity: 5, Price: 195, Currency: "USD",
	}); err != nil {
		t.Fatalf("add: %v", err)
	}

	for i := 0; i < 3; i++ {
		if _, err := run.Execute(context.Background(), p.ID); err != nil {
			t.Fatalf("run[%d]: %v", i, err)
		}
	}

	hist, err := history.Execute(p.ID)
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if len(hist) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(hist))
	}
	for i := 1; i < len(hist); i++ {
		if hist[i].Timestamp.Before(hist[i-1].Timestamp) {
			t.Fatalf("snapshots not ordered oldest-first")
		}
	}
}

func TestHistory_EmptyForUnknownPortfolio(t *testing.T) {
	history := apsnapshot.NewGetHistory(persistence.NewInMemorySnapshotRepository())
	got, err := history.Execute("nope")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}
