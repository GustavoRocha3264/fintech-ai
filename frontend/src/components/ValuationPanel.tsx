import type { Valuation } from '../domain/models';

const fmt = (n: number, currency: string) =>
  new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    maximumFractionDigits: 2,
  }).format(n);

export function ValuationPanel({ valuation }: { valuation: Valuation }) {
  const { totalBRL, totalUSD, percentInBRL, percentInUSD } = valuation;
  const brlWidth = Math.max(0, Math.min(100, percentInBRL));
  const usdWidth = Math.max(0, Math.min(100, percentInUSD));

  return (
    <section className="card">
      <div className="card-header">
        <div className="card-title">Valuation</div>
        <span className="card-subtitle">Live</span>
      </div>

      <div className="metric-grid">
        <div>
          <div className="metric-label">Total in BRL</div>
          <div className="metric-value" style={{ color: 'var(--brl)' }}>
            {fmt(totalBRL.amount, 'BRL')}
          </div>
        </div>
        <div>
          <div className="metric-label">Total in USD</div>
          <div className="metric-value" style={{ color: 'var(--usd)' }}>
            {fmt(totalUSD.amount, 'USD')}
          </div>
        </div>
      </div>

      <div style={{ marginTop: 20 }}>
        <div className="metric-label">Currency allocation</div>
        <div className="alloc-bar" role="img" aria-label={`BRL ${brlWidth.toFixed(0)}%, USD ${usdWidth.toFixed(0)}%`}>
          {brlWidth > 0 && <div className="alloc-bar-brl" style={{ width: `${brlWidth}%` }} />}
          {usdWidth > 0 && <div className="alloc-bar-usd" style={{ width: `${usdWidth}%` }} />}
        </div>
        <div className="alloc-legend">
          <span className="alloc-legend-item">
            <span className="alloc-legend-dot" style={{ background: 'var(--brl)' }} />
            BRL {percentInBRL.toFixed(1)}%
          </span>
          <span className="alloc-legend-item">
            <span className="alloc-legend-dot" style={{ background: 'var(--usd)' }} />
            USD {percentInUSD.toFixed(1)}%
          </span>
        </div>
      </div>
    </section>
  );
}
