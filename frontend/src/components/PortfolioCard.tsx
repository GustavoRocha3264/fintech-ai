import type { Portfolio } from '../domain/models';

export function PortfolioCard({ portfolio }: { portfolio: Portfolio }) {
  return (
    <section style={{ border: '1px solid #ddd', padding: 16, borderRadius: 8 }}>
      <h2>Portfolio {portfolio.id}</h2>
      <p>Base currency: {portfolio.baseCurrency}</p>
      <p>Positions: {portfolio.positions.length}</p>
    </section>
  );
}
