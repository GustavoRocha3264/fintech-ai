import type { Position } from '../domain/models';

const fmt = (n: number, currency: string) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency, maximumFractionDigits: 2 }).format(n);

export function PositionsTable({ positions }: { positions: Position[] }) {
  if (positions.length === 0) {
    return <p style={{ color: '#888' }}>No positions yet — add one below.</p>;
  }
  return (
    <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: 8 }}>
      <thead>
        <tr style={{ textAlign: 'left', borderBottom: '1px solid #ddd' }}>
          <th>Symbol</th>
          <th>Quantity</th>
          <th>Price</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody>
        {positions.map((p) => (
          <tr key={p.id} style={{ borderBottom: '1px solid #f0f0f0' }}>
            <td>{p.symbol}</td>
            <td>{p.quantity}</td>
            <td>{fmt(p.price, p.currency)}</td>
            <td>{fmt(p.price * p.quantity, p.currency)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
