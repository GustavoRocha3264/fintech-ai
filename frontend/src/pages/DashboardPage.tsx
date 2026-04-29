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
    reset,
  } = usePortfolioStore();

  useEffect(() => {
    if (portfolioId && !dashboard && !loading) {
      void loadDashboard();
    }
  }, [portfolioId, dashboard, loading, loadDashboard]);

  return (
    <main style={{ fontFamily: 'system-ui, sans-serif', maxWidth: 820, margin: '40px auto', padding: '0 16px' }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0 }}>Cross-Border Portfolio Intelligence</h1>
        {portfolioId && (
          <button onClick={reset} style={{ fontSize: 12 }}>
            Switch portfolio
          </button>
        )}
      </header>

      {error && (
        <p style={{ color: 'crimson', background: '#fee', padding: 8, borderRadius: 4, marginTop: 16 }}>
          {error}
        </p>
      )}

      {!portfolioId && <CreatePortfolioForm onSubmit={createPortfolio} busy={busy} />}

      {portfolioId && !dashboard && loading && <p>Loading dashboard…</p>}

      {dashboard && (
        <>
          <section style={{ marginTop: 16, border: '1px solid #ddd', padding: 16, borderRadius: 8 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline' }}>
              <h2 style={{ margin: 0 }}>Portfolio</h2>
              <span style={{ fontSize: 12, color: '#666' }}>
                {dashboard.portfolio.id} · base {dashboard.portfolio.baseCurrency} · 1 USD ={' '}
                {dashboard.fx.usdToBRL.toFixed(4)} BRL
              </span>
            </div>
            <PositionsTable positions={dashboard.portfolio.positions} />
            <AddPositionForm onSubmit={addPosition} busy={busy} />
          </section>

          <ValuationPanel valuation={dashboard.valuation} />

          <div style={{ marginTop: 16 }}>
            <button onClick={runAnalysis} disabled={busy || dashboard.portfolio.positions.length === 0}>
              {busy ? 'Working…' : 'Run analysis'}
            </button>
          </div>

          <InsightsPanel report={dashboard.latestReport} />
          <SnapshotsChart snapshots={dashboard.snapshots} />
        </>
      )}
    </main>
  );
}
