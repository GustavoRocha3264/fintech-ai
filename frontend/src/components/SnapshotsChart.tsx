import type { PortfolioSnapshot } from '../domain/models';

interface Props {
  snapshots: PortfolioSnapshot[];
}

const W = 720;
const H = 220;
const PAD_L = 56;
const PAD_R = 16;
const PAD_T = 16;
const PAD_B = 32;

const fmtCurrency = (n: number) =>
  new Intl.NumberFormat('en-US', {
    notation: 'compact',
    maximumFractionDigits: 1,
  }).format(n);

const fmtDate = (t: number) => {
  const d = new Date(t);
  return `${d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`;
};

export function SnapshotsChart({ snapshots }: Props) {
  if (snapshots.length === 0) {
    return (
      <section className="card">
        <div className="card-header">
          <div className="card-title">History</div>
        </div>
        <div className="empty">
          No snapshots yet — each analysis run captures a portfolio value snapshot.
        </div>
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

  const vMinRaw = Math.min(...series.map((p) => p.v));
  const vMaxRaw = Math.max(...series.map((p) => p.v));
  const vPad = (vMaxRaw - vMinRaw) * 0.1 || vMaxRaw * 0.05 || 1;
  const vMin = Math.max(0, vMinRaw - vPad);
  const vMax = vMaxRaw + vPad;
  const vSpan = Math.max(vMax - vMin, 1);

  const x = (t: number) => PAD_L + ((t - tMin) / tSpan) * (W - PAD_L - PAD_R);
  const y = (v: number) => PAD_T + (1 - (v - vMin) / vSpan) * (H - PAD_T - PAD_B);

  const linePath = series
    .map((p, i) => `${i === 0 ? 'M' : 'L'} ${x(p.t).toFixed(1)} ${y(p.v).toFixed(1)}`)
    .join(' ');

  const areaPath =
    series.length > 1
      ? `${linePath} L ${x(series[series.length - 1].t).toFixed(1)} ${y(vMin).toFixed(1)} L ${x(series[0].t).toFixed(1)} ${y(vMin).toFixed(1)} Z`
      : '';

  const yTicks = [vMin, vMin + vSpan / 2, vMax];

  const latest = series[series.length - 1];
  const first = series[0];
  const change = latest.v - first.v;
  const changePct = first.v ? (change / first.v) * 100 : 0;
  const positive = change >= 0;

  return (
    <section className="card">
      <div className="card-header">
        <div>
          <div className="card-title">History</div>
          <p className="hint">Total portfolio value (BRL) per snapshot</p>
        </div>
        {snapshots.length > 1 && (
          <div style={{ textAlign: 'right' }}>
            <div className="metric-label">Δ since first snapshot</div>
            <div
              className="metric-value"
              style={{ fontSize: 16, color: positive ? 'var(--success)' : 'var(--danger)' }}
            >
              {positive ? '+' : ''}
              {change.toFixed(2)} BRL ({positive ? '+' : ''}
              {changePct.toFixed(2)}%)
            </div>
          </div>
        )}
      </div>

      <svg
        viewBox={`0 0 ${W} ${H}`}
        style={{ width: '100%', height: 'auto', display: 'block' }}
        role="img"
        aria-label="Portfolio value history"
      >
        <defs>
          <linearGradient id="snapshot-area" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="var(--primary)" stopOpacity="0.18" />
            <stop offset="100%" stopColor="var(--primary)" stopOpacity="0" />
          </linearGradient>
        </defs>

        {/* Y gridlines */}
        {yTicks.map((tick) => (
          <g key={tick}>
            <line
              x1={PAD_L}
              x2={W - PAD_R}
              y1={y(tick)}
              y2={y(tick)}
              stroke="var(--border)"
              strokeDasharray="3 3"
            />
            <text
              x={PAD_L - 8}
              y={y(tick) + 3}
              fontSize="10"
              fill="var(--text-faint)"
              textAnchor="end"
            >
              {fmtCurrency(tick)}
            </text>
          </g>
        ))}

        {/* X axis dates */}
        {[first, latest].map((p, i) => (
          <text
            key={i}
            x={i === 0 ? PAD_L : W - PAD_R}
            y={H - 12}
            fontSize="10"
            fill="var(--text-faint)"
            textAnchor={i === 0 ? 'start' : 'end'}
          >
            {fmtDate(p.t)}
          </text>
        ))}

        {/* Area fill */}
        {areaPath && <path d={areaPath} fill="url(#snapshot-area)" />}

        {/* Line */}
        {series.length > 1 && (
          <path d={linePath} fill="none" stroke="var(--primary)" strokeWidth={2} strokeLinejoin="round" />
        )}

        {/* Points */}
        {series.map((p, i) => (
          <circle
            key={i}
            cx={x(p.t)}
            cy={y(p.v)}
            r={3.5}
            fill="white"
            stroke="var(--primary)"
            strokeWidth={2}
          >
            <title>{`${fmtDate(p.t)} — ${p.v.toFixed(2)} BRL`}</title>
          </circle>
        ))}
      </svg>

      <div className="hint" style={{ marginTop: 8 }}>
        {snapshots.length} snapshot{snapshots.length === 1 ? '' : 's'}
      </div>
    </section>
  );
}
