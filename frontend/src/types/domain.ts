export type Currency = 'BRL' | 'USD';

export interface Portfolio {
  id: string;
  ownerId: string;
  name: string;
  baseCurrency: Currency;
}

export interface RiskMetrics {
  volatility: number;
  beta: number;
  var95: number;
  sharpe: number;
}

export interface Report {
  id: string;
  portfolioId: string;
  generatedAt: string;
  risk: RiskMetrics;
  narrative: string;
}
