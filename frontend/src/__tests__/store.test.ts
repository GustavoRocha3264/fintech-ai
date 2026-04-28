import { describe, it, expect } from 'vitest';
import { usePortfolioStore } from '../store/portfolioStore';

describe('portfolio store', () => {
  it('starts with empty state', () => {
    const s = usePortfolioStore.getState();
    expect(s.current).toBeNull();
    expect(s.report).toBeNull();
    expect(s.loading).toBe(false);
  });
});
