import type { AnalysisReport } from '../domain/models';

export function InsightsPanel({ report }: { report: AnalysisReport | null }) {
  if (!report) {
    return (
      <section style={{ marginTop: 16, border: '1px dashed #ccc', padding: 16, borderRadius: 8 }}>
        <h3 style={{ marginTop: 0 }}>Insights</h3>
        <p style={{ color: '#888' }}>No analysis yet — run one to generate insights.</p>
      </section>
    );
  }
  return (
    <section style={{ marginTop: 16, border: '1px solid #ddd', padding: 16, borderRadius: 8 }}>
      <h3 style={{ marginTop: 0 }}>Insights</h3>
      <div style={{ fontSize: 12, color: '#666', marginBottom: 8 }}>
        Generated {new Date(report.createdAt).toLocaleString()} · Top concentration{' '}
        {report.topAssetConcentrationPercent.toFixed(1)}%
      </div>
      {report.insights.length === 0 ? (
        <p style={{ color: '#3a7' }}>Portfolio looks balanced — no flags raised.</p>
      ) : (
        <ul>
          {report.insights.map((line, i) => (
            <li key={i}>{line}</li>
          ))}
        </ul>
      )}
    </section>
  );
}
