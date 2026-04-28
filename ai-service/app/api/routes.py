from fastapi import APIRouter

from app.schemas.report import GenerateReportRequest, ReportResponse
from app.services.report_generator import ReportGenerator

router = APIRouter(prefix="/v1", tags=["analysis"])
_generator = ReportGenerator()


@router.post("/reports", response_model=ReportResponse)
def generate_report(req: GenerateReportRequest) -> ReportResponse:
    return _generator.generate(req)
