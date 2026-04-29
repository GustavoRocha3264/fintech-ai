import type { AnalysisResponse } from '../types/domain';

export function ReportPanel({ report }: { report: AnalysisResponse | null }) {
  if (!report) return <p>No report yet.</p>;
  return (
    <section style={{ marginTop: 24 }}>
      <h3>Latest report</h3>
      <p>Generated: {report.createdAt}</p>
      <ul>
        <li>Total BRL: {report.totalValueBRL}</li>
        <li>Total USD: {report.totalValueUSD}</li>
        <li>BRL exposure: {report.brlExposurePercent}%</li>
        <li>USD exposure: {report.usdExposurePercent}%</li>
        <li>Top asset concentration: {report.topAssetConcentrationPercent}%</li>
      </ul>
      <ul>
        {report.insights.map((insight) => (
          <li key={insight}>{insight}</li>
        ))}
      </ul>
    </section>
  );
}
