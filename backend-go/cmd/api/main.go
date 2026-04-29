package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apanalysis "github.com/fintech/cbpi/backend-go/internal/application/analysis"
	"github.com/fintech/cbpi/backend-go/internal/application/dashboard"
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	apsnapshot "github.com/fintech/cbpi/backend-go/internal/application/snapshot"
	appuow "github.com/fintech/cbpi/backend-go/internal/application/uow"
	domainanalysis "github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	domainportfolio "github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	domainsnapshot "github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
	pgstore "github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence/postgres"
	infrauow "github.com/fintech/cbpi/backend-go/internal/infrastructure/uow"
	httpiface "github.com/fintech/cbpi/backend-go/internal/interfaces/http"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/handlers"

	_ "github.com/lib/pq"
)

func main() {
	addr := envOr("HTTP_ADDR", ":8080")
	databaseURL := os.Getenv("DATABASE_URL")
	fxAPIURL := envOr("FX_API_URL", "https://open.er-api.com/v6/latest")
	fxTTL := envDurationOr("FX_CACHE_TTL", 5*time.Minute)

	marketProvider := market.NewStubMarketDataProvider()
	fxProvider := fx.NewFallback(
		fx.NewHTTPProvider(fxAPIURL, fxTTL),
		fx.NewStubFXRateProvider(),
	)

	portfolioRepo, analysisRepo, snapshotRepo, unitOfWork, persistenceMode, cleanup := mustBuildPersistence(databaseURL)
	defer cleanup()

	valuationSvc := apportfolio.NewValuationService(marketProvider, fxProvider)
	createUC := apportfolio.NewCreatePortfolio(portfolioRepo)
	getUC := apportfolio.NewGetPortfolio(portfolioRepo)
	addPosUC := apportfolio.NewAddPosition(portfolioRepo)
	getValuedUC := apportfolio.NewGetPortfolioWithValuation(portfolioRepo, valuationSvc)
	runAnalysisUC := apanalysis.NewRunAnalysis(unitOfWork, valuationSvc)
	latestAnalysisUC := apanalysis.NewGetLatestAnalysis(analysisRepo)
	historyUC := apsnapshot.NewGetHistory(snapshotRepo)
	dashboardUC := dashboard.NewGetDashboard(portfolioRepo, valuationSvc, analysisRepo, snapshotRepo, fxProvider)

	portfolioHandler := handlers.NewPortfolioHandler(createUC, getUC, addPosUC, getValuedUC)
	analysisHandler := handlers.NewAnalysisHandler(runAnalysisUC, latestAnalysisUC)
	snapshotHandler := handlers.NewSnapshotHandler(historyUC)
	fxHandler := handlers.NewFXHandler(fxProvider)
	dashboardHandler := handlers.NewDashboardHandler(dashboardUC)

	srv := &http.Server{
		Addr:              addr,
		Handler:           httpiface.NewRouter(portfolioHandler, analysisHandler, snapshotHandler, fxHandler, dashboardHandler),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("persistence mode: %s", persistenceMode)
		log.Printf("api listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envDurationOr(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func mustBuildPersistence(databaseURL string) (
	portfolioRepo domainportfolio.Repository,
	analysisRepo domainanalysis.Repository,
	snapshotRepo domainsnapshot.Repository,
	unitOfWork appuow.UnitOfWork,
	mode string,
	cleanup func(),
) {
	if databaseURL == "" {
		p := persistence.NewInMemoryPortfolioRepository()
		a := persistence.NewInMemoryAnalysisRepository()
		s := persistence.NewInMemorySnapshotRepository()

		return p, a, s, infrauow.NewInMemoryUoW(p, a, s), "in-memory", func() {}
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("open postgres connection: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		log.Fatalf("ping postgres: %v", err)
	}

	return pgstore.NewPortfolioRepo(db),
		pgstore.NewAnalysisRepo(db),
		pgstore.NewSnapshotRepo(db),
		infrauow.NewPostgresUoW(db, nil),
		"postgres",
		func() {
			if err := db.Close(); err != nil {
				log.Printf("close postgres connection: %v", err)
			}
		}
}
