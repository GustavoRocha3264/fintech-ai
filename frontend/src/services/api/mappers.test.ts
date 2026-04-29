import { describe, expect, it } from 'vitest';
import type { DashboardResponseDto, PortfolioResponseDto } from './dto';
import { toAnalysisReport, toCreatePortfolioRequestDto, toDashboard, toPortfolio } from './mappers';

describe('api mappers', () => {
  it('maps portfolio dto to frontend portfolio model', () => {
    const dto: PortfolioResponseDto = {
      id: 'portfolio-1',
      baseCurrency: 'USD',
      createdAt: '2026-04-29T12:00:00Z',
      positions: [
        {
          id: 'position-1',
          symbol: 'AAPL',
          quantity: 10,
          price: 195,
          currency: 'USD',
        },
      ],
    };

    const portfolio = toPortfolio(dto);

    expect(portfolio).toEqual({
      id: 'portfolio-1',
      baseCurrency: 'USD',
      createdAt: '2026-04-29T12:00:00Z',
      positions: [
        {
          id: 'position-1',
          symbol: 'AAPL',
          quantity: 10,
          price: 195,
          currency: 'USD',
        },
      ],
    });
  });

  it('clones insights when mapping analysis dto', () => {
    const dto = {
      id: 'analysis-1',
      portfolioId: 'portfolio-1',
      createdAt: '2026-04-29T12:00:00Z',
      totalValueBRL: 1000,
      totalValueUSD: 200,
      brlExposurePercent: 60,
      usdExposurePercent: 40,
      topAssetConcentrationPercent: 30,
      insights: ['Diversified enough'],
    };

    const report = toAnalysisReport(dto);
    dto.insights.push('Mutated later');

    expect(report.insights).toEqual(['Diversified enough']);
  });

  it('maps dashboard dto with nullable latest report', () => {
    const dto: DashboardResponseDto = {
      portfolio: {
        id: 'portfolio-1',
        baseCurrency: 'BRL',
        createdAt: '2026-04-29T12:00:00Z',
        positions: [],
      },
      valuation: {
        totalBRL: { amount: 1000, currency: 'BRL' },
        totalUSD: { amount: 200, currency: 'USD' },
        percentInBRL: 70,
        percentInUSD: 30,
      },
      latestReport: null,
      snapshots: [
        {
          id: 'snapshot-1',
          portfolioId: 'portfolio-1',
          timestamp: '2026-04-29T12:00:00Z',
          totalValueBRL: 1000,
          totalValueUSD: 200,
        },
      ],
      fx: {
        usdToBRL: 5,
        brlToUSD: 0.2,
      },
    };

    const dashboard = toDashboard(dto);

    expect(dashboard.latestReport).toBeNull();
    expect(dashboard.snapshots[0]?.portfolioId).toBe('portfolio-1');
    expect(dashboard.fx.usdToBRL).toBe(5);
  });

  it('maps domain input to request dto', () => {
    expect(toCreatePortfolioRequestDto({ baseCurrency: 'BRL' })).toEqual({
      baseCurrency: 'BRL',
    });
  });
});
