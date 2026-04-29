package analysis

import (
	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

type RunAnalysis struct {
	portfolios portfolio.Repository
	reports    analysis.Repository
	snapshots  snapshot.Repository
	valuation  apportfolio.ValuationService
}

func NewRunAnalysis(
	p portfolio.Repository,
	r analysis.Repository,
	s snapshot.Repository,
	v apportfolio.ValuationService,
) *RunAnalysis {
	return &RunAnalysis{portfolios: p, reports: r, snapshots: s, valuation: v}
}

func (uc *RunAnalysis) Execute(portfolioID string) (*analysis.AnalysisReport, error) {
	p, err := uc.portfolios.FindByID(portfolioID)
	if err != nil {
		return nil, err
	}

	res, err := uc.valuation.Calculate(*p)
	if err != nil {
		return nil, err
	}
	concentration := portfolio.TopAssetConcentration(p.Positions, res.Prices, res.USDToBRL)

	report := analysis.NewReport(analysis.Input{
		PortfolioID:                  p.ID,
		TotalValueBRL:                res.Valuation.TotalBRL.Amount,
		TotalValueUSD:                res.Valuation.TotalUSD.Amount,
		BRLExposurePercent:           res.Valuation.PercentInBRL,
		USDExposurePercent:           res.Valuation.PercentInUSD,
		TopAssetConcentrationPercent: concentration,
		PositionCount:                len(p.Positions),
	})

	if err := uc.reports.Save(*report); err != nil {
		return nil, err
	}
	if err := uc.snapshots.Save(*snapshot.New(p.ID, res.Valuation.TotalBRL.Amount, res.Valuation.TotalUSD.Amount)); err != nil {
		return nil, err
	}
	return report, nil
}

type GetLatestAnalysis struct {
	reports analysis.Repository
}

func NewGetLatestAnalysis(r analysis.Repository) *GetLatestAnalysis {
	return &GetLatestAnalysis{reports: r}
}

func (uc *GetLatestAnalysis) Execute(portfolioID string) (*analysis.AnalysisReport, error) {
	return uc.reports.GetLatestByPortfolioID(portfolioID)
}
