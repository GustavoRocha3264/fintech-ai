import { useEffect, useRef, useState } from 'react';
import type { AddPositionInput, MarketQuote, MarketSymbol } from '../domain/models';
import { api } from '../services/api';

interface Props {
  onSubmit: (input: AddPositionInput) => void;
  busy: boolean;
}

type ResolveState = 'idle' | 'loading' | 'resolved' | 'error';

export function AddPositionForm({ onSubmit, busy }: Props) {
  const [symbols, setSymbols] = useState<MarketSymbol[]>([]);
  const [symbolInput, setSymbolInput] = useState('');
  const [quantity, setQuantity] = useState('');
  const [price, setPrice] = useState('');
  const [resolveState, setResolveState] = useState<ResolveState>('idle');
  const [resolvedQuote, setResolvedQuote] = useState<MarketQuote | null>(null);
  const [resolveError, setResolveError] = useState<string | null>(null);
  const [symbolsError, setSymbolsError] = useState<string | null>(null);

  // Tracks which symbol the in-flight request was made for, so stale results
  // from a previous symbol are discarded if the user typed something else.
  const pendingSymbol = useRef('');

  useEffect(() => {
    api
      .getMarketSymbols()
      .then((s) => {
        setSymbols(s);
        setSymbolsError(null);
      })
      .catch((e: Error) => setSymbolsError(e.message));
  }, []);

  const resolveFor = async (sym: string) => {
    if (!sym || resolvedQuote?.symbol === sym) return;

    pendingSymbol.current = sym;
    setResolveState('loading');
    setResolveError(null);

    try {
      const quote = await api.getMarketQuote(sym);
      if (pendingSymbol.current !== sym) return;
      setResolvedQuote(quote);
      setResolveState('resolved');
      setPrice((p) => p || quote.price.toFixed(2));
    } catch (e) {
      if (pendingSymbol.current !== sym) return;
      setResolveState('error');
      setResolveError((e as Error).message);
    }
  };

  const handleSymbolChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value.toUpperCase();
    setSymbolInput(val);
    setResolvedQuote(null);
    setResolveState('idle');
    setResolveError(null);
    pendingSymbol.current = '';

    if (symbols.some((s) => s.ticker === val)) {
      void resolveFor(val);
    }
  };

  const handleSymbolBlur = () => {
    const sym = symbolInput.trim().toUpperCase();
    if (sym && resolveState === 'idle') void resolveFor(sym);
  };

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!resolvedQuote) return;
    const q = Number(quantity);
    const p = Number(price);
    if (!Number.isFinite(q) || q <= 0 || !Number.isFinite(p) || p <= 0) return;
    onSubmit({
      symbol: resolvedQuote.symbol,
      quantity: q,
      price: p,
      currency: resolvedQuote.currency,
    });
    reset();
  };

  const reset = () => {
    setSymbolInput('');
    setQuantity('');
    setPrice('');
    setResolvedQuote(null);
    setResolveState('idle');
    setResolveError(null);
    pendingSymbol.current = '';
  };

  return (
    <div style={{ marginTop: 16 }}>
      <h3 style={{ marginBottom: 10 }}>Add position</h3>
      <form onSubmit={submit} className="form-row cols-position">
        <label className="field">
          Symbol
          <input
            list="market-symbols-list"
            value={symbolInput}
            onChange={handleSymbolChange}
            onBlur={handleSymbolBlur}
            placeholder="AAPL, PETR4…"
            disabled={busy}
            required
            autoComplete="off"
            style={{ textTransform: 'uppercase' }}
          />
          <datalist id="market-symbols-list">
            {symbols.map((s) => (
              <option key={s.ticker} value={s.ticker}>
                {s.name} — {s.currency}
              </option>
            ))}
          </datalist>
        </label>

        <label className="field">
          Quantity
          <input
            value={quantity}
            onChange={(e) => setQuantity(e.target.value)}
            type="number"
            step="any"
            min="0"
            placeholder="0"
            disabled={busy}
            required
          />
        </label>

        <label className="field">
          Price{resolvedQuote ? ` (${resolvedQuote.currency})` : ''}
          <input
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            type="number"
            step="any"
            min="0"
            placeholder="0.00"
            disabled={busy || resolveState === 'loading'}
            required
          />
        </label>

        <div className="currency-badge-wrap">
          <label className="field" style={{ marginBottom: 0 }}>
            Currency
            <CurrencyBadge state={resolveState} currency={resolvedQuote?.currency ?? null} />
          </label>
        </div>

        <button
          type="submit"
          className="primary"
          disabled={busy || resolveState !== 'resolved'}
          style={{ alignSelf: 'end', height: 38 }}
        >
          {busy ? '…' : 'Add'}
        </button>
      </form>

      {resolveError && (
        <p className="hint error" style={{ marginTop: 8 }}>
          Couldn’t resolve <strong>{symbolInput}</strong>: {resolveError}
        </p>
      )}
      {symbolsError && (
        <p className="hint error" style={{ marginTop: 8 }}>
          Symbol catalog unavailable: {symbolsError}
        </p>
      )}
      {!symbolsError && symbols.length === 0 && (
        <p className="hint" style={{ marginTop: 8 }}>Loading symbol catalog…</p>
      )}
    </div>
  );
}

function CurrencyBadge({ state, currency }: { state: ResolveState; currency: string | null }) {
  if (state === 'loading') {
    return (
      <div className="currency-badge loading" style={{ height: 38 }}>
        <span className="spinner" />
      </div>
    );
  }
  if (state === 'resolved' && currency) {
    return <div className={`currency-badge ${currency.toLowerCase()}`} style={{ height: 38 }}>{currency}</div>;
  }
  return <div className="currency-badge empty" style={{ height: 38 }}>—</div>;
}
