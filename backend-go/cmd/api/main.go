package main

import (
	"context"
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
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
	httpiface "github.com/fintech/cbpi/backend-go/internal/interfaces/http"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/handlers"
)

func main() {
	addr := envOr("HTTP_ADDR", ":8080")
	fxAPIURL := envOr("FX_API_URL", "https://open.er-api.com/v6/latest")
	fxTTL := envDurationOr("FX_CACHE_TTL", 5*time.Minute)

	portfolioRepo := persistence.NewInMemoryPortfolioRepository()
	analysisRepo := persistence.NewInMemoryAnalysisRepository()
	snapshotRepo := persistence.NewInMemorySnapshotRepository()
	marketProvider := market.NewStubMarketDataProvider()
	fxProvider := fx.NewFallback(
		fx.NewHTTPProvider(fxAPIURL, fxTTL),
		fx.NewStubFXRateProvider(),
	)

	valuationSvc := apportfolio.NewValuationService(marketProvider, fxProvider)

	createUC := apportfolio.NewCreatePortfolio(portfolioRepo)
	getUC := apportfolio.NewGetPortfolio(portfolioRepo)
	addPosUC := apportfolio.NewAddPosition(portfolioRepo)
	getValuedUC := apportfolio.NewGetPortfolioWithValuation(portfolioRepo, valuationSvc)
	runAnalysisUC := apanalysis.NewRunAnalysis(portfolioRepo, analysisRepo, snapshotRepo, valuationSvc)
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
