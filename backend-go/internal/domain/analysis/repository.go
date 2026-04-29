package analysis

import "errors"

var ErrNotFound = errors.New("analysis report not found")

type Repository interface {
	Save(report AnalysisReport) error
	GetLatestByPortfolioID(portfolioID string) (*AnalysisReport, error)
}
