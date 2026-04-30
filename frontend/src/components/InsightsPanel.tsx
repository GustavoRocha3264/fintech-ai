import type { AnalysisReport } from '../domain/models';

export function InsightsPanel({ report }: { report: AnalysisReport | null }) {
  if (!report) {
    return <div className="empty">No analysis yet — click <strong>Run analysis</strong> to generate insights.</div>;
  }

  const generated = new Date(report.createdAt);

  return (
    <div>
      <div className="metric-grid" style={{ marginBottom: 16 }}>
        <div>
          <div className="metric-label">Top concentration</div>
          <div className="metric-value">{report.topAssetConcentrationPercent.toFixed(1)}%</div>
        </div>
        <div>
          <div className="metric-label">Generated</div>
          <div className="metric-value" style={{ fontSize: 14, fontWeight: 500, color: 'var(--text-muted)' }}>
            {generated.toLocaleString()}
          </div>
        </div>
      </div>

      {report.insights.length === 0 ? (
        <div className="insight-good">
          <span aria-hidden>✓</span>
          <span>Portfolio looks balanced — no flags raised.</span>
        </div>
      ) : (
        <ul className="insight-list">
          {report.insights.map((line, i) => (
            <li key={i} className="insight-item">
              <span className="insight-icon" aria-hidden>!</span>
              <span>{line}</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
