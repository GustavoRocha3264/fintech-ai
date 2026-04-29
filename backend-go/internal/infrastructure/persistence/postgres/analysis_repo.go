package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
)

// Insights are persisted as a delimited string to keep this adapter
// stdlib-only. Switch to pq.Array / pgx text[] if a driver dependency is
// later acceptable.
type AnalysisRepo struct {
	tx dbtx
}

func NewAnalysisRepo(tx dbtx) *AnalysisRepo { return &AnalysisRepo{tx: tx} }

const insertAnalysis = `
INSERT INTO analysis_reports (
  id, portfolio_id, created_at,
  total_value_brl, total_value_usd,
  brl_exposure_pct, usd_exposure_pct,
  top_asset_concentration_pct, insights
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
`

const selectLatestAnalysis = `
SELECT id, portfolio_id, created_at,
       total_value_brl, total_value_usd,
       brl_exposure_pct, usd_exposure_pct,
       top_asset_concentration_pct, insights
FROM analysis_reports
WHERE portfolio_id = $1
ORDER BY created_at DESC
LIMIT 1
`

const insightsSep = "\x1f" // ASCII unit separator — won't appear in human text.

func (r *AnalysisRepo) Save(report analysis.AnalysisReport) error {
	_, err := r.tx.ExecContext(context.Background(), insertAnalysis,
		report.ID, report.PortfolioID, report.CreatedAt,
		report.TotalValueBRL, report.TotalValueUSD,
		report.BRLExposurePercent, report.USDExposurePercent,
		report.TopAssetConcentrationPercent,
		strings.Join(report.Insights, insightsSep),
	)
	return err
}

func (r *AnalysisRepo) GetLatestByPortfolioID(portfolioID string) (*analysis.AnalysisReport, error) {
	var rep analysis.AnalysisReport
	var insights string
	row := r.tx.QueryRowContext(context.Background(), selectLatestAnalysis, portfolioID)
	if err := row.Scan(
		&rep.ID, &rep.PortfolioID, &rep.CreatedAt,
		&rep.TotalValueBRL, &rep.TotalValueUSD,
		&rep.BRLExposurePercent, &rep.USDExposurePercent,
		&rep.TopAssetConcentrationPercent,
		&insights,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, analysis.ErrNotFound
		}
		return nil, err
	}
	if insights != "" {
		rep.Insights = strings.Split(insights, insightsSep)
	}
	return &rep, nil
}
