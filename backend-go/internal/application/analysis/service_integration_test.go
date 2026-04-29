package analysis_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	apanalysis "github.com/fintech/cbpi/backend-go/internal/application/analysis"
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/application/uow"
	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	pgstore "github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence/postgres"
	infrauow "github.com/fintech/cbpi/backend-go/internal/infrastructure/uow"

	_ "github.com/lib/pq"
)

func TestRunAnalysis_PostgresIntegration_CommitsReportAndSnapshot(t *testing.T) {
	db := openIntegrationDB(t)
	resetIntegrationDB(t, db)

	portfolioRepo := pgstore.NewPortfolioRepo(db)
	analysisRepo := pgstore.NewAnalysisRepo(db)
	snapshotRepo := pgstore.NewSnapshotRepo(db)
	uow := infrauow.NewPostgresUoW(db, nil)
	valuation := apportfolio.NewValuationService(market.NewStubMarketDataProvider(), fx.NewStubFXRateProvider())

	create := apportfolio.NewCreatePortfolio(portfolioRepo)
	add := apportfolio.NewAddPosition(portfolioRepo)
	run := apanalysis.NewRunAnalysis(uow, valuation)
	latest := apanalysis.NewGetLatestAnalysis(analysisRepo)
	history := apsnapshotReader{repo: snapshotRepo}

	p, err := create.Execute("USD")
	if err != nil {
		t.Fatalf("create portfolio: %v", err)
	}
	if _, err := add.Execute(apportfolio.AddPositionInput{
		PortfolioID: p.ID,
		Symbol:      "AAPL",
		Quantity:    100,
		Price:       195,
		Currency:    "USD",
	}); err != nil {
		t.Fatalf("add position: %v", err)
	}

	report, err := run.Execute(context.Background(), p.ID)
	if err != nil {
		t.Fatalf("run analysis: %v", err)
	}

	gotReport, err := latest.Execute(p.ID)
	if err != nil {
		t.Fatalf("latest analysis: %v", err)
	}
	if gotReport.ID != report.ID {
		t.Fatalf("expected latest report %q, got %q", report.ID, gotReport.ID)
	}

	snapshots, err := history.FindByPortfolioID(p.ID)
	if err != nil {
		t.Fatalf("fetch snapshots: %v", err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(snapshots))
	}
	if snapshots[0].PortfolioID != p.ID {
		t.Fatalf("expected snapshot portfolio %q, got %q", p.ID, snapshots[0].PortfolioID)
	}
}

func TestPostgresUoW_RollsBackAnalysisAndSnapshotOnError(t *testing.T) {
	db := openIntegrationDB(t)
	resetIntegrationDB(t, db)

	portfolioRepo := pgstore.NewPortfolioRepo(db)
	analysisRepo := pgstore.NewAnalysisRepo(db)
	snapshotRepo := pgstore.NewSnapshotRepo(db)
	unitOfWork := infrauow.NewPostgresUoW(db, nil)

	create := apportfolio.NewCreatePortfolio(portfolioRepo)
	p, err := create.Execute("USD")
	if err != nil {
		t.Fatalf("create portfolio: %v", err)
	}

	wantErr := errors.New("force rollback")
	err = unitOfWork.Do(context.Background(), func(repos uow.Repositories) error {
		report := analysis.AnalysisReport{
			ID:                           "report-rollback",
			PortfolioID:                  p.ID,
			CreatedAt:                    time.Now().UTC(),
			TotalValueBRL:                10,
			TotalValueUSD:                2,
			BRLExposurePercent:           50,
			USDExposurePercent:           50,
			TopAssetConcentrationPercent: 100,
			Insights:                     []string{"test rollback"},
		}
		if err := repos.Analysis.Save(report); err != nil {
			t.Fatalf("save analysis inside uow: %v", err)
		}

		snap := snapshot.PortfolioSnapshot{
			ID:            "snapshot-rollback",
			PortfolioID:   p.ID,
			Timestamp:     time.Now().UTC(),
			TotalValueBRL: 10,
			TotalValueUSD: 2,
		}
		if err := repos.Snapshot.Save(snap); err != nil {
			t.Fatalf("save snapshot inside uow: %v", err)
		}

		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected rollback error %v, got %v", wantErr, err)
	}

	if _, err := analysisRepo.GetLatestByPortfolioID(p.ID); !errors.Is(err, analysis.ErrNotFound) {
		t.Fatalf("expected no persisted analysis after rollback, got %v", err)
	}

	snapshots, err := snapshotRepo.FindByPortfolioID(p.ID)
	if err != nil {
		t.Fatalf("fetch snapshots after rollback: %v", err)
	}
	if len(snapshots) != 0 {
		t.Fatalf("expected 0 snapshots after rollback, got %d", len(snapshots))
	}
}

type apsnapshotReader struct {
	repo snapshot.Repository
}

func (r apsnapshotReader) FindByPortfolioID(portfolioID string) ([]snapshot.PortfolioSnapshot, error) {
	return r.repo.FindByPortfolioID(portfolioID)
}

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run Postgres integration tests")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close test database: %v", err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping test database: %v", err)
	}

	applyMigration(t, db)
	return db
}

func applyMigration(t *testing.T, db *sql.DB) {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve current file path")
	}
	migrationPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "migrations", "001_init.sql")
	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	statements := strings.Split(string(sqlBytes), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("apply migration statement %q: %v", stmt, err)
		}
	}
}

func resetIntegrationDB(t *testing.T, db *sql.DB) {
	t.Helper()

	if _, err := db.Exec(`
		TRUNCATE TABLE
			portfolio_snapshots,
			analysis_reports,
			positions,
			portfolios
		RESTART IDENTITY CASCADE
	`); err != nil {
		t.Fatalf("reset test database: %v", err)
	}
}
