package postgres

import (
	"context"

	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

type SnapshotRepo struct {
	tx dbtx
}

func NewSnapshotRepo(tx dbtx) *SnapshotRepo { return &SnapshotRepo{tx: tx} }

const insertSnapshot = `
INSERT INTO portfolio_snapshots (
  id, portfolio_id, timestamp, total_value_brl, total_value_usd
) VALUES ($1,$2,$3,$4,$5)
`

const selectSnapshots = `
SELECT id, portfolio_id, timestamp, total_value_brl, total_value_usd
FROM portfolio_snapshots
WHERE portfolio_id = $1
ORDER BY timestamp ASC
`

func (r *SnapshotRepo) Save(s snapshot.PortfolioSnapshot) error {
	_, err := r.tx.ExecContext(context.Background(), insertSnapshot,
		s.ID, s.PortfolioID, s.Timestamp, s.TotalValueBRL, s.TotalValueUSD,
	)
	return err
}

func (r *SnapshotRepo) FindByPortfolioID(portfolioID string) ([]snapshot.PortfolioSnapshot, error) {
	rows, err := r.tx.QueryContext(context.Background(), selectSnapshots, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []snapshot.PortfolioSnapshot
	for rows.Next() {
		var s snapshot.PortfolioSnapshot
		if err := rows.Scan(&s.ID, &s.PortfolioID, &s.Timestamp, &s.TotalValueBRL, &s.TotalValueUSD); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
