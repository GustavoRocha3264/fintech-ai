from datetime import datetime
from pydantic import BaseModel


class Holding(BaseModel):
    asset_id: str
    quantity: float
    currency: str


class GenerateReportRequest(BaseModel):
    portfolio_id: str
    holdings: list[Holding] = []


class RiskMetrics(BaseModel):
    volatility: float
    beta: float
    var_95: float
    sharpe: float


class ReportResponse(BaseModel):
    id: str
    portfolio_id: str
    generated_at: datetime
    risk: RiskMetrics
    narrative: str
