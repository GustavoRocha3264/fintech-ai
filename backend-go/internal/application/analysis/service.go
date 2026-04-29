package analysis

import (
	"context"

	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/application/uow"
	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

type RunAnalysis struct {
	uow       uow.UnitOfWork
	valuation apportfolio.ValuationService
}

func NewRunAnalysis(u uow.UnitOfWork, v apportfolio.ValuationService) *RunAnalysis {
	return &RunAnalysis{uow: u, valuation: v}
}

// Execute runs an analysis for the given portfolio and atomically persists
// both the resulting report and a value snapshot. If either save fails (or
// any other step inside the UoW), neither is committed.
func (uc *RunAnalysis) Execute(ctx context.Context, portfolioID string) (*analysis.AnalysisReport, error) {
	var report *analysis.AnalysisReport

	err := uc.uow.Do(ctx, func(repos uow.Repositories) error {
		p, err := repos.Portfolio.FindByID(portfolioID)
		if err != nil {
			return err
		}

		res, err := uc.valuation.Calculate(*p)
		if err != nil {
			return err
		}
		concentration := portfolio.TopAssetConcentration(p.Positions, res.Prices, res.USDToBRL)

		r := analysis.NewReport(analysis.Input{
			PortfolioID:                  p.ID,
			TotalValueBRL:                res.Valuation.TotalBRL.Amount,
			TotalValueUSD:                res.Valuation.TotalUSD.Amount,
			BRLExposurePercent:           res.Valuation.PercentInBRL,
			USDExposurePercent:           res.Valuation.PercentInUSD,
			TopAssetConcentrationPercent: concentration,
			PositionCount:                len(p.Positions),
		})

		if err := repos.Analysis.Save(*r); err != nil {
			return err
		}
		if err := repos.Snapshot.Save(*snapshot.New(p.ID, res.Valuation.TotalBRL.Amount, res.Valuation.TotalUSD.Amount)); err != nil {
			return err
		}
		report = r
		return nil
	})
	if err != nil {
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
