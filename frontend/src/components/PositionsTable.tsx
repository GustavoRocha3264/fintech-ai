import type { Position } from '../domain/models';

const fmt = (n: number, currency: string) =>
  new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    maximumFractionDigits: 2,
  }).format(n);

const fmtQty = (n: number) =>
  new Intl.NumberFormat('en-US', { maximumFractionDigits: 4 }).format(n);

export function PositionsTable({ positions }: { positions: Position[] }) {
  if (positions.length === 0) {
    return <div className="empty">No positions yet — add one below to get started.</div>;
  }
  return (
    <table className="data">
      <thead>
        <tr>
          <th>Symbol</th>
          <th>Currency</th>
          <th style={{ textAlign: 'right' }}>Quantity</th>
          <th style={{ textAlign: 'right' }}>Price</th>
          <th style={{ textAlign: 'right' }}>Value</th>
        </tr>
      </thead>
      <tbody>
        {positions.map((p) => (
          <tr key={p.id}>
            <td className="symbol">{p.symbol}</td>
            <td>
              <span className={`currency-badge ${p.currency.toLowerCase()}`} style={{ minWidth: 0, height: 22, padding: '0 8px', fontSize: 11 }}>
                {p.currency}
              </span>
            </td>
            <td className="num">{fmtQty(p.quantity)}</td>
            <td className="num">{fmt(p.price, p.currency)}</td>
            <td className="num">
              <strong>{fmt(p.price * p.quantity, p.currency)}</strong>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
