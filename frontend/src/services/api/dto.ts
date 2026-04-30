import type { Currency } from '../../domain/models';

export interface PositionResponseDto {
  id: string;
  symbol: string;
  quantity: number;
  price: number;
  currency: Currency;
}

export interface PortfolioResponseDto {
  id: string;
  baseCurrency: Currency;
  createdAt: string;
  positions: PositionResponseDto[];
}

export interface MoneyResponseDto {
  amount: number;
  currency: Currency;
}

export interface ValuationResponseDto {
  totalBRL: MoneyResponseDto;
  totalUSD: MoneyResponseDto;
  percentInBRL: number;
  percentInUSD: number;
}

export interface PortfolioWithValuationResponseDto {
  portfolio: PortfolioResponseDto;
  valuation: ValuationResponseDto;
}

export interface AnalysisResponseDto {
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

export interface SnapshotResponseDto {
  id: string;
  portfolioId: string;
  timestamp: string;
  totalValueBRL: number;
  totalValueUSD: number;
}

export interface FXResponseDto {
  from: string;
  to: string;
  rate: number;
}

export interface DashboardFXResponseDto {
  usdToBRL: number;
  brlToUSD: number;
}

export interface DashboardResponseDto {
  portfolio: PortfolioResponseDto;
  valuation: ValuationResponseDto;
  latestReport: AnalysisResponseDto | null;
  snapshots: SnapshotResponseDto[];
  fx: DashboardFXResponseDto;
}

export interface CreatePortfolioRequestDto {
  baseCurrency: Currency;
}

export interface AddPositionRequestDto {
  symbol: string;
  quantity: number;
  price: number;
  currency: Currency;
}

export interface MarketSymbolDto {
  ticker: string;
  currency: Currency;
  name: string;
}

export interface MarketQuoteDto {
  symbol: string;
  price: number;
  currency: Currency;
}
