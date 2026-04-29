import type {
  AnalysisReport,
  CreatePortfolioInput,
  Dashboard,
  Portfolio,
  PortfolioWithValuation,
} from '../domain/models';
import type {
  AnalysisResponseDto,
  CreatePortfolioRequestDto,
  DashboardResponseDto,
  PortfolioResponseDto,
  PortfolioWithValuationResponseDto,
} from './api/dto';
import {
  toAnalysisReport,
  toCreatePortfolioRequestDto,
  toDashboard,
  toPortfolio,
  toPortfolioWithValuation,
} from './api/mappers';

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
  async createPortfolio(input: CreatePortfolioInput): Promise<Portfolio> {
    const request: CreatePortfolioRequestDto = toCreatePortfolioRequestDto(input);
    const response = await http<PortfolioResponseDto>('/portfolios', {
      method: 'POST',
      body: JSON.stringify(request),
    });
    return toPortfolio(response);
  },
  async getPortfolio(id: string): Promise<Portfolio> {
    const response = await http<PortfolioResponseDto>(`/portfolios/${id}`);
    return toPortfolio(response);
  },
  async getPortfolioWithValuation(id: string): Promise<PortfolioWithValuation> {
    const response = await http<PortfolioWithValuationResponseDto>(`/portfolios/${id}/valuation`);
    return toPortfolioWithValuation(response);
  },
  async getDashboard(id: string): Promise<Dashboard> {
    const response = await http<DashboardResponseDto>(`/portfolios/${id}/dashboard`);
    return toDashboard(response);
  },
  async runAnalysis(portfolioId: string): Promise<AnalysisReport> {
    const response = await http<AnalysisResponseDto>(`/portfolios/${portfolioId}/analysis`, {
      method: 'POST',
    });
    return toAnalysisReport(response);
  },
  async getLatestReport(portfolioId: string): Promise<AnalysisReport> {
    const response = await http<AnalysisResponseDto>(`/portfolios/${portfolioId}/analysis/latest`);
    return toAnalysisReport(response);
  },
};
