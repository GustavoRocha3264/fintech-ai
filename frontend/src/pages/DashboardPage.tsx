import { useEffect } from 'react';
import { usePortfolioStore } from '../store/portfolioStore';
import { PortfolioCard } from '../components/PortfolioCard';
import { ReportPanel } from '../components/ReportPanel';

const DEMO_ID = 'demo-portfolio-id';

export function DashboardPage() {
  const { current, report, loading, error, load, runAnalysis } = usePortfolioStore();

  useEffect(() => {
    void load(DEMO_ID);
  }, [load]);

  return (
    <main style={{ fontFamily: 'system-ui, sans-serif', maxWidth: 720, margin: '40px auto' }}>
      <h1>Cross-Border Portfolio Intelligence</h1>
      {loading && <p>Loading…</p>}
      {error && <p style={{ color: 'crimson' }}>{error}</p>}
      {current && <PortfolioCard portfolio={current} />}
      <button
        style={{ marginTop: 16 }}
        onClick={() => current && runAnalysis(current.id)}
        disabled={!current || loading}
      >
        Run analysis
      </button>
      <ReportPanel report={report} />
    </main>
  );
}
