import { useState } from 'react';
import type { Currency } from '../domain/models';

interface Props {
  onSubmit: (baseCurrency: Currency) => void;
  busy: boolean;
}

export function CreatePortfolioForm({ onSubmit, busy }: Props) {
  const [baseCurrency, setBaseCurrency] = useState<Currency>('USD');

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit(baseCurrency);
      }}
      style={{ border: '1px solid #ddd', padding: 16, borderRadius: 8 }}
    >
      <h2>Create a portfolio</h2>
      <label style={{ display: 'block', marginBottom: 12 }}>
        Base currency:&nbsp;
        <select
          value={baseCurrency}
          onChange={(e) => setBaseCurrency(e.target.value as Currency)}
          disabled={busy}
        >
          <option value="USD">USD</option>
          <option value="BRL">BRL</option>
        </select>
      </label>
      <button type="submit" disabled={busy}>
        {busy ? 'Creating…' : 'Create'}
      </button>
    </form>
  );
}
