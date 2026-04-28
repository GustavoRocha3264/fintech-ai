from datetime import datetime, timezone
from uuid import uuid4

from app.schemas.report import GenerateReportRequest, ReportResponse, RiskMetrics


class ReportGenerator:
    """LLM-backed narrative generator. Stubbed — wires to a model provider later."""

    def generate(self, req: GenerateReportRequest) -> ReportResponse:
        return ReportResponse(
            id=str(uuid4()),
            portfolio_id=req.portfolio_id,
            generated_at=datetime.now(timezone.utc),
            risk=RiskMetrics(volatility=0.18, beta=1.05, var_95=0.04, sharpe=1.2),
            narrative="Stub narrative — model integration pending.",
        )
