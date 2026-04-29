-- Initial schema for the Cross-Border Portfolio Intelligence Platform.
-- Apply with `psql -f migrations/001_init.sql` against an empty database.

CREATE TABLE IF NOT EXISTS portfolios (
    id            TEXT PRIMARY KEY,
    base_currency TEXT NOT NULL CHECK (base_currency IN ('BRL', 'USD')),
    created_at    TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS positions (
    id           TEXT PRIMARY KEY,
    portfolio_id TEXT NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    symbol       TEXT NOT NULL,
    quantity     DOUBLE PRECISION NOT NULL CHECK (quantity > 0),
    price        DOUBLE PRECISION NOT NULL CHECK (price > 0),
    currency     TEXT NOT NULL CHECK (currency IN ('BRL', 'USD'))
);
CREATE INDEX IF NOT EXISTS positions_by_portfolio ON positions(portfolio_id);

CREATE TABLE IF NOT EXISTS analysis_reports (
    id                          TEXT PRIMARY KEY,
    portfolio_id                TEXT NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    created_at                  TIMESTAMPTZ NOT NULL,
    total_value_brl             DOUBLE PRECISION NOT NULL,
    total_value_usd             DOUBLE PRECISION NOT NULL,
    brl_exposure_pct            DOUBLE PRECISION NOT NULL,
    usd_exposure_pct            DOUBLE PRECISION NOT NULL,
    top_asset_concentration_pct DOUBLE PRECISION NOT NULL,
    insights                    TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS analysis_reports_latest
    ON analysis_reports(portfolio_id, created_at DESC);

CREATE TABLE IF NOT EXISTS portfolio_snapshots (
    id              TEXT PRIMARY KEY,
    portfolio_id    TEXT NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    timestamp       TIMESTAMPTZ NOT NULL,
    total_value_brl DOUBLE PRECISION NOT NULL,
    total_value_usd DOUBLE PRECISION NOT NULL
);
CREATE INDEX IF NOT EXISTS portfolio_snapshots_by_time
    ON portfolio_snapshots(portfolio_id, timestamp ASC);
