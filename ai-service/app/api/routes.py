from fastapi import APIRouter, Depends, Response

from app.api.dependencies import get_chat_service
from app.schemas.chat import ChatRequest, ChatResponse
from app.schemas.report import GenerateReportRequest, ReportResponse
from app.services.chat_service import ChatService
from app.services.report_generator import ReportGenerator


router = APIRouter(prefix="/v1", tags=["ai"])
_report_generator = ReportGenerator()


@router.post("/chat", response_model=ChatResponse)
async def chat(
    req: ChatRequest,
    service: ChatService = Depends(get_chat_service),
) -> ChatResponse:
    return await service.reply(req.messages)


@router.post(
    "/reports",
    response_model=ReportResponse,
    deprecated=True,
    summary="Generate daily report — DEPRECATED, prefer POST /v1/chat",
)
def generate_report(req: GenerateReportRequest, response: Response) -> ReportResponse:
    """Legacy endpoint kept for the Go backend's analysis flow.

    New clients should use POST /v1/chat. The Go side's
    `infrastructure/ai/http_client.go` will be migrated in a follow-up.
    """
    response.headers["Deprecation"] = "true"
    response.headers["Link"] = '</v1/chat>; rel="successor-version"'
    response.headers["Sunset"] = "Wed, 01 Jul 2026 00:00:00 GMT"
    return _report_generator.generate(req)
