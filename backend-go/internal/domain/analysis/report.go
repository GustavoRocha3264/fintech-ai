package analysis

import (
	"time"

	"github.com/google/uuid"
)

type AnalysisReport struct {
	ID          string
	PortfolioID string
	CreatedAt   time.Time

	TotalValueBRL float64
	TotalValueUSD float64

	BRLExposurePercent float64
	USDExposurePercent float64

	TopAssetConcentrationPercent float64

	Insights []string
}

type Input struct {
	PortfolioID                  string
	TotalValueBRL                float64
	TotalValueUSD                float64
	BRLExposurePercent           float64
	USDExposurePercent           float64
	TopAssetConcentrationPercent float64
	PositionCount                int
}

func NewReport(in Input) *AnalysisReport {
	return &AnalysisReport{
		ID:                           uuid.NewString(),
		PortfolioID:                  in.PortfolioID,
		CreatedAt:                    time.Now().UTC(),
		TotalValueBRL:                in.TotalValueBRL,
		TotalValueUSD:                in.TotalValueUSD,
		BRLExposurePercent:           in.BRLExposurePercent,
		USDExposurePercent:           in.USDExposurePercent,
		TopAssetConcentrationPercent: in.TopAssetConcentrationPercent,
		Insights:                     GenerateInsights(in),
	}
}
