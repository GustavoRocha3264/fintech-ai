import { useState } from 'react';
import type { Currency } from '../domain/models';

interface Props {
  onSubmit: (baseCurrency: Currency) => void;
  busy: boolean;
}

export function CreatePortfolioForm({ onSubmit, busy }: Props) {
  const [baseCurrency, setBaseCurrency] = useState<Currency>('USD');

  return (
    <section className="card">
      <div className="card-header">
        <div>
          <div className="card-title">Get started</div>
          <p className="hint">Create a virtual portfolio to start tracking BRL and USD positions.</p>
        </div>
      </div>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          onSubmit(baseCurrency);
        }}
        className="form-row cols-create"
      >
        <label className="field">
          Base currency
          <select
            value={baseCurrency}
            onChange={(e) => setBaseCurrency(e.target.value as Currency)}
            disabled={busy}
          >
            <option value="USD">USD — US Dollar</option>
            <option value="BRL">BRL — Brazilian Real</option>
          </select>
        </label>
        <button type="submit" className="primary" disabled={busy} style={{ height: 38 }}>
          {busy ? 'Creating…' : 'Create portfolio'}
        </button>
      </form>
    </section>
  );
}
