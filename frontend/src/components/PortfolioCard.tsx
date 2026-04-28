import type { Portfolio } from '../types/domain';

export function PortfolioCard({ portfolio }: { portfolio: Portfolio }) {
  return (
    <section style={{ border: '1px solid #ddd', padding: 16, borderRadius: 8 }}>
      <h2>{portfolio.name}</h2>
      <p>Base currency: {portfolio.baseCurrency}</p>
      <p>ID: {portfolio.id}</p>
    </section>
  );
}
