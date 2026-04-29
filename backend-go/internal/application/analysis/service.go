package analysis

import (
	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

type RunAnalysis struct {
	portfolios portfolio.Repository
	reports    analysis.Repository
	snapshots  snapshot.Repository
	market     portfolio.MarketDataProvider
	fx         portfolio.FXRateProvider
}

func NewRunAnalysis(
	p portfolio.Repository,
	r analysis.Repository,
	s snapshot.Repository,
	m portfolio.MarketDataProvider,
	f portfolio.FXRateProvider,
) *RunAnalysis {
	return &RunAnalysis{portfolios: p, reports: r, snapshots: s, market: m, fx: f}
}

func (uc *RunAnalysis) Execute(portfolioID string) (*analysis.AnalysisReport, error) {
	p, err := uc.portfolios.FindByID(portfolioID)
	if err != nil {
		return nil, err
	}

	prices := make(map[string]portfolio.Money, len(p.Positions))
	for _, pos := range p.Positions {
		price, currency, err := uc.market.GetPrice(pos.Symbol)
		if err != nil {
			return nil, err
		}
		prices[pos.Symbol] = portfolio.NewMoney(price, currency)
	}

	rate, err := uc.fx.GetRate(portfolio.CurrencyUSD, portfolio.CurrencyBRL)
	if err != nil {
		return nil, err
	}

	v := portfolio.Valuate(p.Positions, prices, rate)
	concentration := topAssetConcentration(p.Positions, prices, rate)

	report := analysis.NewReport(analysis.Input{
		PortfolioID:                  p.ID,
		TotalValueBRL:                v.TotalBRL.Amount,
		TotalValueUSD:                v.TotalUSD.Amount,
		BRLExposurePercent:           v.PercentInBRL,
		USDExposurePercent:           v.PercentInUSD,
		TopAssetConcentrationPercent: concentration,
		PositionCount:                len(p.Positions),
	})

	if err := uc.reports.Save(*report); err != nil {
		return nil, err
	}
	if err := uc.snapshots.Save(*snapshot.New(p.ID, v.TotalBRL.Amount, v.TotalUSD.Amount)); err != nil {
		return nil, err
	}
	return report, nil
}

func topAssetConcentration(positions []portfolio.Position, prices map[string]portfolio.Money, fxRate float64) float64 {
	var total, top float64
	for _, pos := range positions {
		price, ok := prices[pos.Symbol]
		if !ok {
			price = portfolio.NewMoney(pos.Price, pos.Currency)
		}
		v := pos.Quantity * price.Amount
		if price.Currency == portfolio.CurrencyUSD {
			v *= fxRate
		}
		total += v
		if v > top {
			top = v
		}
	}
	if total == 0 {
		return 0
	}
	return (top / total) * 100
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
