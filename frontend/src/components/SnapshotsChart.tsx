import type { PortfolioSnapshot } from '../domain/models';

interface Props {
  snapshots: PortfolioSnapshot[];
}

const W = 600;
const H = 160;
const PAD = 24;

// Lightweight inline SVG chart — keeps the dependency footprint at zero.
// Plots TotalValueBRL over time. Falls back to a "needs more data" hint until
// at least two snapshots are available.
export function SnapshotsChart({ snapshots }: Props) {
  if (snapshots.length === 0) {
    return (
      <section style={{ marginTop: 16, border: '1px dashed #ccc', padding: 16, borderRadius: 8 }}>
        <h3 style={{ marginTop: 0 }}>History</h3>
        <p style={{ color: '#888' }}>No snapshots yet — run analysis to capture one.</p>
      </section>
    );
  }

  const series = snapshots.map((s) => ({
    t: new Date(s.timestamp).getTime(),
    v: s.totalValueBRL,
  }));
  const tMin = series[0].t;
  const tMax = series[series.length - 1].t;
  const tSpan = Math.max(tMax - tMin, 1);
  const vMin = Math.min(...series.map((p) => p.v));
  const vMax = Math.max(...series.map((p) => p.v));
  const vSpan = Math.max(vMax - vMin, 1);

  const x = (t: number) => PAD + ((t - tMin) / tSpan) * (W - 2 * PAD);
  const y = (v: number) => H - PAD - ((v - vMin) / vSpan) * (H - 2 * PAD);

  const path =
    series.length === 1
      ? ''
      : series.map((p, i) => `${i === 0 ? 'M' : 'L'} ${x(p.t).toFixed(1)} ${y(p.v).toFixed(1)}`).join(' ');

  return (
    <section style={{ marginTop: 16, border: '1px solid #ddd', padding: 16, borderRadius: 8 }}>
      <h3 style={{ marginTop: 0 }}>History (Total Value, BRL)</h3>
      <svg viewBox={`0 0 ${W} ${H}`} style={{ width: '100%', height: 'auto' }}>
        <line x1={PAD} y1={H - PAD} x2={W - PAD} y2={H - PAD} stroke="#ddd" />
        <line x1={PAD} y1={PAD} x2={PAD} y2={H - PAD} stroke="#ddd" />
        {path && <path d={path} fill="none" stroke="#48f" strokeWidth={2} />}
        {series.map((p, i) => (
          <circle key={i} cx={x(p.t)} cy={y(p.v)} r={3} fill="#48f" />
        ))}
        <text x={PAD} y={PAD - 8} fontSize={10} fill="#666">
          {vMax.toFixed(0)}
        </text>
        <text x={PAD} y={H - PAD + 14} fontSize={10} fill="#666">
          {vMin.toFixed(0)}
        </text>
      </svg>
      <div style={{ fontSize: 12, color: '#666' }}>{snapshots.length} snapshot(s)</div>
    </section>
  );
}
