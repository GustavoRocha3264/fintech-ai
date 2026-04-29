import type {
  AnalysisResponse,
  CreatePortfolioRequest,
  DashboardResponse,
  PortfolioResponse,
  PortfolioWithValuationResponse,
} from '../types/domain';

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
  createPortfolio(input: CreatePortfolioRequest) {
    return http<PortfolioResponse>('/portfolios', { method: 'POST', body: JSON.stringify(input) });
  },
  getPortfolio(id: string) {
    return http<PortfolioResponse>(`/portfolios/${id}`);
  },
  getPortfolioWithValuation(id: string) {
    return http<PortfolioWithValuationResponse>(`/portfolios/${id}/valuation`);
  },
  getDashboard(id: string) {
    return http<DashboardResponse>(`/portfolios/${id}/dashboard`);
  },
  runAnalysis(portfolioId: string) {
    return http<AnalysisResponse>(`/portfolios/${portfolioId}/analysis`, {
      method: 'POST',
    });
  },
  getLatestReport(portfolioId: string) {
    return http<AnalysisResponse>(`/portfolios/${portfolioId}/analysis/latest`);
  },
};
