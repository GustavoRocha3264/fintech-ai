package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
	httpiface "github.com/fintech/cbpi/backend-go/internal/interfaces/http"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/handlers"
)

func main() {
	addr := envOr("HTTP_ADDR", ":8080")

	repo := persistence.NewInMemoryPortfolioRepository()
	marketProvider := market.NewStubMarketDataProvider()
	fxProvider := fx.NewStubFXRateProvider()

	createUC := apportfolio.NewCreatePortfolio(repo)
	getUC := apportfolio.NewGetPortfolio(repo)
	addPosUC := apportfolio.NewAddPosition(repo)
	getValuedUC := apportfolio.NewGetPortfolioWithValuation(repo, marketProvider, fxProvider)
	handler := handlers.NewPortfolioHandler(createUC, getUC, addPosUC, getValuedUC)

	srv := &http.Server{
		Addr:              addr,
		Handler:           httpiface.NewRouter(handler),
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
