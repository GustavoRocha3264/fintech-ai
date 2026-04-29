import type { Valuation } from '../domain/models';

const fmt = (n: number, currency: string) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency, maximumFractionDigits: 2 }).format(n);

export function ValuationPanel({ valuation }: { valuation: Valuation }) {
  const { totalBRL, totalUSD, percentInBRL, percentInUSD } = valuation;
  return (
    <section style={{ marginTop: 16, border: '1px solid #ddd', padding: 16, borderRadius: 8 }}>
      <h3 style={{ marginTop: 0 }}>Valuation</h3>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
        <div>
          <div style={{ fontSize: 12, color: '#666' }}>Total BRL</div>
          <div style={{ fontSize: 20, fontWeight: 600 }}>{fmt(totalBRL.amount, 'BRL')}</div>
        </div>
        <div>
          <div style={{ fontSize: 12, color: '#666' }}>Total USD</div>
          <div style={{ fontSize: 20, fontWeight: 600 }}>{fmt(totalUSD.amount, 'USD')}</div>
        </div>
      </div>
      <div style={{ marginTop: 12 }}>
        <div style={{ fontSize: 12, color: '#666', marginBottom: 4 }}>Allocation</div>
        <div style={{ display: 'flex', height: 12, borderRadius: 6, overflow: 'hidden', background: '#eee' }}>
          <div
            style={{ width: `${percentInBRL}%`, background: '#4f8' }}
            title={`BRL ${percentInBRL.toFixed(1)}%`}
          />
          <div
            style={{ width: `${percentInUSD}%`, background: '#48f' }}
            title={`USD ${percentInUSD.toFixed(1)}%`}
          />
        </div>
        <div style={{ fontSize: 12, color: '#666', marginTop: 4 }}>
          BRL {percentInBRL.toFixed(1)}% · USD {percentInUSD.toFixed(1)}%
        </div>
      </div>
    </section>
  );
}
