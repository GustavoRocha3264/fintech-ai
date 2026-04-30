export type Currency = 'BRL' | 'USD';

export interface Position {
  id: string;
  symbol: string;
  quantity: number;
  price: number;
  currency: Currency;
}

export interface Portfolio {
  id: string;
  baseCurrency: Currency;
  createdAt: string;
  positions: Position[];
}

export interface Money {
  amount: number;
  currency: Currency;
}

export interface Valuation {
  totalBRL: Money;
  totalUSD: Money;
  percentInBRL: number;
  percentInUSD: number;
}

export interface PortfolioWithValuation {
  portfolio: Portfolio;
  valuation: Valuation;
}

export interface AnalysisReport {
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

export interface PortfolioSnapshot {
  id: string;
  portfolioId: string;
  timestamp: string;
  totalValueBRL: number;
  totalValueUSD: number;
}

export interface FXRate {
  from: string;
  to: string;
  rate: number;
}

export interface DashboardFX {
  usdToBRL: number;
  brlToUSD: number;
}

export interface Dashboard {
  portfolio: Portfolio;
  valuation: Valuation;
  latestReport: AnalysisReport | null;
  snapshots: PortfolioSnapshot[];
  fx: DashboardFX;
}

export interface CreatePortfolioInput {
  baseCurrency: Currency;
}

export interface AddPositionInput {
  symbol: string;
  quantity: number;
  price: number;
  currency: Currency;
}

export interface MarketSymbol {
  ticker: string;
  currency: Currency;
  name: string;
}

export interface MarketQuote {
  symbol: string;
  price: number;
  currency: Currency;
}
