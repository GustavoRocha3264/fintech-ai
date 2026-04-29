// @vitest-environment jsdom
import { afterEach, describe, expect, it } from 'vitest';
import { usePortfolioStore } from '../store/portfolioStore';

afterEach(() => {
  usePortfolioStore.getState().reset();
});

describe('portfolio store', () => {
  it('starts with empty state', () => {
    const s = usePortfolioStore.getState();
    expect(s.dashboard).toBeNull();
    expect(s.loading).toBe(false);
    expect(s.busy).toBe(false);
    expect(s.error).toBeNull();
  });

  it('reset clears the in-memory portfolio reference', () => {
    usePortfolioStore.setState({ portfolioId: 'abc', dashboard: null });
    expect(usePortfolioStore.getState().portfolioId).toBe('abc');
    usePortfolioStore.getState().reset();
    expect(usePortfolioStore.getState().portfolioId).toBeNull();
  });
});
