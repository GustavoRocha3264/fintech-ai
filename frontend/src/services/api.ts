import type { Portfolio, Report, Currency } from '../types/domain';

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1';

async function http<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...init,
  });
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  return res.json() as Promise<T>;
}

export const api = {
  createPortfolio(input: { ownerId: string; name: string; baseCurrency: Currency }) {
    return http<Portfolio>('/portfolio', { method: 'POST', body: JSON.stringify(input) });
  },
  getPortfolio(id: string) {
    return http<Portfolio>(`/portfolio/${id}`);
  },
  runAnalysis(portfolioId: string) {
    return http<Report>('/analysis/run', {
      method: 'POST',
      body: JSON.stringify({ portfolioId }),
    });
  },
  getLatestReport(portfolioId: string) {
    return http<Report>(`/analysis/${portfolioId}`);
  },
};
