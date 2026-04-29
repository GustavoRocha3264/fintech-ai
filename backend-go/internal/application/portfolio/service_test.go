package portfolio_test

import (
	"errors"
	"testing"

	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
	"github.com/fintech/cbpi/backend-go/internal/infrastructure/persistence"
)

func TestCreateAndGetPortfolio(t *testing.T) {
	repo := persistence.NewInMemoryPortfolioRepository()
	create := apportfolio.NewCreatePortfolio(repo)
	get := apportfolio.NewGetPortfolio(repo)

	p, err := create.Execute("USD")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := get.Execute(p.ID)
	if err != nil || got.ID != p.ID {
		t.Fatalf("get failed: %v / %v", err, got)
	}
}

func TestCreatePortfolio_RejectsInvalidCurrency(t *testing.T) {
	repo := persistence.NewInMemoryPortfolioRepository()
	create := apportfolio.NewCreatePortfolio(repo)

	_, err := create.Execute("EUR")
	if !errors.Is(err, portfolio.ErrInvalidBaseCurrency) {
		t.Fatalf("expected ErrInvalidBaseCurrency, got %v", err)
	}
}

func TestAddPositionAndValuate(t *testing.T) {
	repo := persistence.NewInMemoryPortfolioRepository()
	create := apportfolio.NewCreatePortfolio(repo)
	add := apportfolio.NewAddPosition(repo)
	view := apportfolio.NewGetPortfolioWithValuation(repo, apportfolio.NewValuationService(market.NewStubMarketDataProvider(), fx.NewStubFXRateProvider()))

	p, err := create.Execute("USD")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := add.Execute(apportfolio.AddPositionInput{
		PortfolioID: p.ID, Symbol: "AAPL", Quantity: 10, Price: 195.0, Currency: "USD",
	}); err != nil {
		t.Fatalf("add USD: %v", err)
	}
	if _, err := add.Execute(apportfolio.AddPositionInput{
		PortfolioID: p.ID, Symbol: "PETR4", Quantity: 100, Price: 38.0, Currency: "BRL",
	}); err != nil {
		t.Fatalf("add BRL: %v", err)
	}

	v, err := view.Execute(p.ID)
	if err != nil {
		t.Fatalf("view: %v", err)
	}
	if len(v.Portfolio.Positions) != 2 {
		t.Fatalf("expected 2 positions, got %d", len(v.Portfolio.Positions))
	}
	if v.Valuation.TotalBRL.Amount <= 0 || v.Valuation.TotalUSD.Amount <= 0 {
		t.Fatalf("expected positive totals: %+v", v.Valuation)
	}
	if v.Valuation.PercentInBRL+v.Valuation.PercentInUSD < 99.99 {
		t.Fatalf("percentages don't sum to 100: %+v", v.Valuation)
	}
}

func TestAddPosition_PortfolioNotFound(t *testing.T) {
	repo := persistence.NewInMemoryPortfolioRepository()
	add := apportfolio.NewAddPosition(repo)

	_, err := add.Execute(apportfolio.AddPositionInput{
		PortfolioID: "missing", Symbol: "AAPL", Quantity: 1, Price: 1, Currency: "USD",
	})
	if !errors.Is(err, portfolio.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
