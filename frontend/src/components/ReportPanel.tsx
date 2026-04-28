import type { Report } from '../types/domain';

export function ReportPanel({ report }: { report: Report | null }) {
  if (!report) return <p>No report yet.</p>;
  return (
    <section style={{ marginTop: 24 }}>
      <h3>Latest report</h3>
      <p>Generated: {report.generatedAt}</p>
      <ul>
        <li>Volatility: {report.risk.volatility}</li>
        <li>Beta: {report.risk.beta}</li>
        <li>VaR 95: {report.risk.var95}</li>
        <li>Sharpe: {report.risk.sharpe}</li>
      </ul>
      <p>{report.narrative}</p>
    </section>
  );
}
