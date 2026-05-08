import { DashboardPage } from './pages/DashboardPage';
import { usePortfolioStore } from './store/portfolioStore';

export function App() {
  const { portfolioId, reset } = usePortfolioStore();

  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="app-header-inner">
          <div className="app-brand">
            <div className="app-logo" aria-hidden>CB</div>
            <div>
              <h1>Cross-Border Portfolio Intelligence</h1>
              <p className="hint">BRL · USD multi-currency portfolio tracker</p>
            </div>
          </div>
          {portfolioId && (
            <button className="ghost" onClick={reset} title="Switch to a different portfolio">
              Switch portfolio
            </button>
          )}
        </div>
      </header>

      <DashboardPage />
    </div>
  );
}
