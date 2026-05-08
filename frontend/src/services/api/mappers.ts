import type {
  AddPositionInput,
  AnalysisReport,
  CreatePortfolioInput,
  Dashboard,
  FXRate,
  MarketQuote,
  MarketSymbol,
  Money,
  Portfolio,
  PortfolioSnapshot,
  PortfolioWithValuation,
  Position,
  Valuation,
} from '../../domain/models';
import type {
  AddPositionRequestDto,
  AnalysisResponseDto,
  CreatePortfolioRequestDto,
  DashboardResponseDto,
  FXResponseDto,
  MarketQuoteDto,
  MarketSymbolDto,
  MoneyResponseDto,
  PortfolioResponseDto,
  PortfolioWithValuationResponseDto,
  PositionResponseDto,
  SnapshotResponseDto,
  ValuationResponseDto,
} from './dto';

export function toCreatePortfolioRequestDto(input: CreatePortfolioInput): CreatePortfolioRequestDto {
  return { baseCurrency: input.baseCurrency };
}

export function toAddPositionRequestDto(input: AddPositionInput): AddPositionRequestDto {
  return {
    symbol: input.symbol,
    quantity: input.quantity,
    price: input.price,
    currency: input.currency,
  };
}

export function toPosition(dto: PositionResponseDto): Position {
  return {
    id: dto.id,
    symbol: dto.symbol,
    quantity: dto.quantity,
    price: dto.price,
    currency: dto.currency,
  };
}

export function toPortfolio(dto: PortfolioResponseDto): Portfolio {
  return {
    id: dto.id,
    baseCurrency: dto.baseCurrency,
    createdAt: dto.createdAt,
    positions: dto.positions.map(toPosition),
  };
}

export function toMoney(dto: MoneyResponseDto): Money {
  return {
    amount: dto.amount,
    currency: dto.currency,
  };
}

export function toValuation(dto: ValuationResponseDto): Valuation {
  return {
    totalBRL: toMoney(dto.totalBRL),
    totalUSD: toMoney(dto.totalUSD),
    percentInBRL: dto.percentInBRL,
    percentInUSD: dto.percentInUSD,
  };
}

export function toPortfolioWithValuation(dto: PortfolioWithValuationResponseDto): PortfolioWithValuation {
  return {
    portfolio: toPortfolio(dto.portfolio),
    valuation: toValuation(dto.valuation),
  };
}

export function toAnalysisReport(dto: AnalysisResponseDto): AnalysisReport {
  return {
    id: dto.id,
    portfolioId: dto.portfolioId,
    createdAt: dto.createdAt,
    totalValueBRL: dto.totalValueBRL,
    totalValueUSD: dto.totalValueUSD,
    brlExposurePercent: dto.brlExposurePercent,
    usdExposurePercent: dto.usdExposurePercent,
    topAssetConcentrationPercent: dto.topAssetConcentrationPercent,
    insights: [...dto.insights],
  };
}

export function toPortfolioSnapshot(dto: SnapshotResponseDto): PortfolioSnapshot {
  return {
    id: dto.id,
    portfolioId: dto.portfolioId,
    timestamp: dto.timestamp,
    totalValueBRL: dto.totalValueBRL,
    totalValueUSD: dto.totalValueUSD,
  };
}

export function toFXRate(dto: FXResponseDto): FXRate {
  return {
    from: dto.from,
    to: dto.to,
    rate: dto.rate,
  };
}

export function toMarketSymbol(dto: MarketSymbolDto): MarketSymbol {
  return { ticker: dto.ticker, currency: dto.currency, name: dto.name };
}

export function toMarketQuote(dto: MarketQuoteDto): MarketQuote {
  return { symbol: dto.symbol, price: dto.price, currency: dto.currency };
}

export function toDashboard(dto: DashboardResponseDto): Dashboard {
  return {
    portfolio: toPortfolio(dto.portfolio),
    valuation: toValuation(dto.valuation),
    latestReport: dto.latestReport ? toAnalysisReport(dto.latestReport) : null,
    snapshots: dto.snapshots.map(toPortfolioSnapshot),
    fx: {
      usdToBRL: dto.fx.usdToBRL,
      brlToUSD: dto.fx.brlToUSD,
    },
  };
}
