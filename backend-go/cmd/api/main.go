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
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
	httpiface "github.com/fintech/cbpi/backend-go/internal/interfaces/http"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/handlers"
)

func main() {
	addr := envOr("HTTP_ADDR", ":8080")

	portfolioRepo := persistence.NewInMemoryPortfolioRepository()
	analysisRepo := persistence.NewInMemoryAnalysisRepository()
	marketProvider := market.NewStubMarketDataProvider()
	fxProvider := fx.NewStubFXRateProvider()

	createUC := apportfolio.NewCreatePortfolio(portfolioRepo)
	getUC := apportfolio.NewGetPortfolio(portfolioRepo)
	addPosUC := apportfolio.NewAddPosition(portfolioRepo)
	getValuedUC := apportfolio.NewGetPortfolioWithValuation(portfolioRepo, marketProvider, fxProvider)
	runAnalysisUC := apanalysis.NewRunAnalysis(portfolioRepo, analysisRepo, marketProvider, fxProvider)
	latestAnalysisUC := apanalysis.NewGetLatestAnalysis(analysisRepo)

	portfolioHandler := handlers.NewPortfolioHandler(createUC, getUC, addPosUC, getValuedUC)
	analysisHandler := handlers.NewAnalysisHandler(runAnalysisUC, latestAnalysisUC)

	srv := &http.Server{
		Addr:              addr,
		Handler:           httpiface.NewRouter(portfolioHandler, analysisHandler),
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
