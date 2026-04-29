package dashboard

import (
	"errors"

	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

// Dashboard is the aggregated read-model returned by GetDashboard. It
// composes data already produced by other use cases — no new business rules
// live here.
type Dashboard struct {
	Portfolio    *portfolio.Portfolio
	Valuation    portfolio.Valuation
	LatestReport *analysis.AnalysisReport // nil when no analysis has run yet
	Snapshots    []snapshot.PortfolioSnapshot
	USDToBRL     float64
	BRLToUSD     float64
}

type GetDashboard struct {
	portfolios portfolio.Repository
	valuation  apportfolio.ValuationService
	analyses   analysis.Repository
	snapshots  snapshot.Repository
	fx         portfolio.FXRateProvider
}

func NewGetDashboard(
	p portfolio.Repository,
	v apportfolio.ValuationService,
	a analysis.Repository,
	s snapshot.Repository,
	fx portfolio.FXRateProvider,
) *GetDashboard {
	return &GetDashboard{portfolios: p, valuation: v, analyses: a, snapshots: s, fx: fx}
}

func (uc *GetDashboard) Execute(portfolioID string) (*Dashboard, error) {
	p, err := uc.portfolios.FindByID(portfolioID)
	if err != nil {
		return nil, err
	}

	res, err := uc.valuation.Calculate(*p)
	if err != nil {
		return nil, err
	}

	var report *analysis.AnalysisReport
	if r, err := uc.analyses.GetLatestByPortfolioID(p.ID); err == nil {
		report = r
	} else if !errors.Is(err, analysis.ErrNotFound) {
		return nil, err
	}

	snaps, err := uc.snapshots.FindByPortfolioID(p.ID)
	if err != nil {
		return nil, err
	}

	brlToUSD, err := uc.fx.GetRate(portfolio.CurrencyBRL, portfolio.CurrencyUSD)
	if err != nil {
		return nil, err
	}

	return &Dashboard{
		Portfolio:    p,
		Valuation:    res.Valuation,
		LatestReport: report,
		Snapshots:    snaps,
		USDToBRL:     res.USDToBRL,
		BRLToUSD:     brlToUSD,
	}, nil
}
