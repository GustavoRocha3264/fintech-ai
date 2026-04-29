import { useState } from 'react';
import type { AddPositionInput, Currency } from '../domain/models';

interface Props {
  onSubmit: (input: AddPositionInput) => void;
  busy: boolean;
}

export function AddPositionForm({ onSubmit, busy }: Props) {
  const [symbol, setSymbol] = useState('');
  const [quantity, setQuantity] = useState('');
  const [price, setPrice] = useState('');
  const [currency, setCurrency] = useState<Currency>('USD');

  const reset = () => {
    setSymbol('');
    setQuantity('');
    setPrice('');
  };

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    const q = Number(quantity);
    const p = Number(price);
    if (!symbol.trim() || !Number.isFinite(q) || q <= 0 || !Number.isFinite(p) || p <= 0) return;
    onSubmit({ symbol: symbol.trim().toUpperCase(), quantity: q, price: p, currency });
    reset();
  };

  return (
    <form
      onSubmit={submit}
      style={{
        border: '1px solid #ddd',
        padding: 16,
        borderRadius: 8,
        marginTop: 16,
        display: 'grid',
        gridTemplateColumns: '1fr 1fr 1fr 1fr auto',
        gap: 8,
        alignItems: 'end',
      }}
    >
      <label>
        Symbol
        <input
          value={symbol}
          onChange={(e) => setSymbol(e.target.value)}
          placeholder="AAPL"
          disabled={busy}
          required
          style={{ width: '100%' }}
        />
      </label>
      <label>
        Quantity
        <input
          value={quantity}
          onChange={(e) => setQuantity(e.target.value)}
          type="number"
          step="any"
          min="0"
          disabled={busy}
          required
          style={{ width: '100%' }}
        />
      </label>
      <label>
        Price
        <input
          value={price}
          onChange={(e) => setPrice(e.target.value)}
          type="number"
          step="any"
          min="0"
          disabled={busy}
          required
          style={{ width: '100%' }}
        />
      </label>
      <label>
        Currency
        <select
          value={currency}
          onChange={(e) => setCurrency(e.target.value as Currency)}
          disabled={busy}
          style={{ width: '100%' }}
        >
          <option value="USD">USD</option>
          <option value="BRL">BRL</option>
        </select>
      </label>
      <button type="submit" disabled={busy}>
        {busy ? '…' : 'Add'}
      </button>
    </form>
  );
}
