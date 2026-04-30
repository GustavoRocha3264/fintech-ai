import { useEffect } from 'react';
import { usePortfolioStore } from '../store/portfolioStore';
import { CreatePortfolioForm } from '../components/CreatePortfolioForm';
import { AddPositionForm } from '../components/AddPositionForm';
import { PositionsTable } from '../components/PositionsTable';
import { ValuationPanel } from '../components/ValuationPanel';
import { InsightsPanel } from '../components/InsightsPanel';
import { SnapshotsChart } from '../components/SnapshotsChart';

export function DashboardPage() {
  const {
    portfolioId,
    dashboard,
    loading,
    busy,
    error,
    createPortfolio,
    loadDashboard,
    addPosition,
    runAnalysis,
  } = usePortfolioStore();

  // Fetch the dashboard once per portfolio change. We intentionally exclude
  // `loading` and `error` from deps so a failed fetch doesn't auto-retry in a
  // loop — the user can press "Switch portfolio" or fix the problem and the
  // next portfolioId change will trigger a new attempt.
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => {
    if (portfolioId && !dashboard) {
      void loadDashboard();
    }
  }, [portfolioId]);

  return (
    <main className="app-main">
      {error && <div className="banner error">{error}</div>}

      {!portfolioId && <CreatePortfolioForm onSubmit={createPortfolio} busy={busy} />}

      {portfolioId && !dashboard && loading && (
        <div className="card">
          <div className="actions">
            <span className="spinner" />
            <span className="hint">Loading dashboard…</span>
          </div>
        </div>
      )}

      {dashboard && (
        <>
          <section className="card">
            <div className="card-header">
              <div>
                <div className="card-title">Portfolio</div>
                <div className="hint">
                  Base currency <strong>{dashboard.portfolio.baseCurrency}</strong> · 1 USD ={' '}
                  {dashboard.fx.usdToBRL.toFixed(4)} BRL
                </div>
              </div>
              <div className="card-subtitle">{dashboard.portfolio.id.slice(0, 8)}…</div>
            </div>
            <PositionsTable positions={dashboard.portfolio.positions} />
            <AddPositionForm onSubmit={addPosition} busy={busy} />
          </section>

          <ValuationPanel valuation={dashboard.valuation} />

          <section className="card">
            <div className="card-header">
              <div className="card-title">Analysis</div>
              <button
                className="primary"
                onClick={runAnalysis}
                disabled={busy || dashboard.portfolio.positions.length === 0}
              >
                {busy ? (
                  <span className="actions">
                    <span className="spinner" /> Running…
                  </span>
                ) : (
                  'Run analysis'
                )}
              </button>
            </div>
            <InsightsPanel report={dashboard.latestReport} />
          </section>

          <SnapshotsChart snapshots={dashboard.snapshots} />
        </>
      )}
    </main>
  );
}
