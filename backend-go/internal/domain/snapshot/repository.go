package snapshot

type Repository interface {
	Save(snapshot PortfolioSnapshot) error
	FindByPortfolioID(portfolioID string) ([]PortfolioSnapshot, error)
}
