package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

type PortfolioRepo struct {
	tx dbtx
}

func NewPortfolioRepo(tx dbtx) *PortfolioRepo {
	return &PortfolioRepo{tx: tx}
}

const upsertPortfolio = `
INSERT INTO portfolios (id, base_currency, created_at)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE SET base_currency = EXCLUDED.base_currency
`

const insertPosition = `
INSERT INTO positions (id, portfolio_id, symbol, quantity, price, currency)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO NOTHING
`

const deletePositions = `DELETE FROM positions WHERE portfolio_id = $1`

const selectPortfolio = `
SELECT id, base_currency, created_at FROM portfolios WHERE id = $1
`

const selectPositions = `
SELECT id, portfolio_id, symbol, quantity, price, currency
FROM positions WHERE portfolio_id = $1
`

func (r *PortfolioRepo) Save(p portfolio.Portfolio) error {
	ctx := context.Background()
	if _, err := r.tx.ExecContext(ctx, upsertPortfolio, p.ID, p.BaseCurrency, p.CreatedAt); err != nil {
		return err
	}
	// Replace positions wholesale — keeps the implementation simple. A
	// production version would diff and update in place.
	if _, err := r.tx.ExecContext(ctx, deletePositions, p.ID); err != nil {
		return err
	}
	for _, pos := range p.Positions {
		if _, err := r.tx.ExecContext(ctx, insertPosition,
			pos.ID, pos.PortfolioID, pos.Symbol, pos.Quantity, pos.Price, pos.Currency,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *PortfolioRepo) FindByID(id string) (*portfolio.Portfolio, error) {
	ctx := context.Background()
	var p portfolio.Portfolio
	row := r.tx.QueryRowContext(ctx, selectPortfolio, id)
	if err := row.Scan(&p.ID, &p.BaseCurrency, &p.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, portfolio.ErrNotFound
		}
		return nil, err
	}

	rows, err := r.tx.QueryContext(ctx, selectPositions, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pos portfolio.Position
		if err := rows.Scan(&pos.ID, &pos.PortfolioID, &pos.Symbol, &pos.Quantity, &pos.Price, &pos.Currency); err != nil {
			return nil, err
		}
		p.Positions = append(p.Positions, pos)
	}
	return &p, rows.Err()
}
