export type Currency = 'BRL' | 'USD';

export interface PositionResponse {
  id: string;
  symbol: string;
  quantity: number;
  price: number;
  currency: Currency;
}

export interface PortfolioResponse {
  id: string;
  baseCurrency: Currency;
  createdAt: string;
  positions: PositionResponse[];
}

export interface MoneyResponse {
  amount: number;
  currency: Currency;
}

export interface ValuationResponse {
  totalBRL: MoneyResponse;
  totalUSD: MoneyResponse;
  percentInBRL: number;
  percentInUSD: number;
}

export interface PortfolioWithValuationResponse {
  portfolio: PortfolioResponse;
  valuation: ValuationResponse;
}

export interface AnalysisResponse {
  id: string;
  portfolioId: string;
  createdAt: string;
  totalValueBRL: number;
  totalValueUSD: number;
  brlExposurePercent: number;
  usdExposurePercent: number;
  topAssetConcentrationPercent: number;
  insights: string[];
}

export interface SnapshotResponse {
  id: string;
  portfolioId: string;
  timestamp: string;
  totalValueBRL: number;
  totalValueUSD: number;
}

export interface FXResponse {
  from: string;
  to: string;
  rate: number;
}

export interface DashboardFXResponse {
  usdToBRL: number;
  brlToUSD: number;
}

export interface DashboardResponse {
  portfolio: PortfolioResponse;
  valuation: ValuationResponse;
  latestReport: AnalysisResponse | null;
  snapshots: SnapshotResponse[];
  fx: DashboardFXResponse;
}

export interface CreatePortfolioRequest {
  baseCurrency: Currency;
}

export interface AddPositionRequest {
  symbol: string;
  quantity: number;
  price: number;
  currency: Currency;
}
