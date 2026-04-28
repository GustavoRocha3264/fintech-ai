package portfolio_test

import (
	"errors"
	"testing"

	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
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
	if p.ID == "" {
		t.Fatalf("expected generated id")
	}
	if p.BaseCurrency != "USD" {
		t.Fatalf("base currency = %q", p.BaseCurrency)
	}
	if p.CreatedAt.IsZero() {
		t.Fatalf("created at must be set")
	}

	got, err := get.Execute(p.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != p.ID {
		t.Fatalf("id mismatch: %s vs %s", got.ID, p.ID)
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

func TestGetPortfolio_NotFound(t *testing.T) {
	repo := persistence.NewInMemoryPortfolioRepository()
	get := apportfolio.NewGetPortfolio(repo)

	_, err := get.Execute("missing")
	if !errors.Is(err, portfolio.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
