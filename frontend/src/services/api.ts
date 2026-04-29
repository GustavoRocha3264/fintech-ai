import type {
  AddPositionInput,
  AnalysisReport,
  CreatePortfolioInput,
  Dashboard,
  Portfolio,
  PortfolioSnapshot,
  PortfolioWithValuation,
  Position,
} from '../domain/models';
import type {
  AnalysisResponseDto,
  CreatePortfolioRequestDto,
  DashboardResponseDto,
  PortfolioResponseDto,
  PortfolioWithValuationResponseDto,
  PositionResponseDto,
  SnapshotResponseDto,
} from './api/dto';
import {
  toAddPositionRequestDto,
  toAnalysisReport,
  toCreatePortfolioRequestDto,
  toDashboard,
  toPortfolio,
  toPortfolioSnapshot,
  toPortfolioWithValuation,
  toPosition,
} from './api/mappers';

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1';

async function http<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...init,
  });
  if (!res.ok) {
    const body = await res.text();
    let message = `${res.status} ${res.statusText}`;
    try {
      const parsed = JSON.parse(body) as { error?: string };
      if (parsed.error) message = parsed.error;
    } catch {
      if (body) message = body;
    }
    throw new Error(message);
  }
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
  async addPosition(portfolioId: string, input: AddPositionInput): Promise<Position> {
    const response = await http<PositionResponseDto>(`/portfolios/${portfolioId}/positions`, {
      method: 'POST',
      body: JSON.stringify(toAddPositionRequestDto(input)),
    });
    return toPosition(response);
  },
  async getSnapshots(portfolioId: string): Promise<PortfolioSnapshot[]> {
    const response = await http<SnapshotResponseDto[]>(`/portfolios/${portfolioId}/snapshots`);
    return response.map(toPortfolioSnapshot);
  },
};
