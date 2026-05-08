from fastapi import APIRouter, Depends, HTTPException, Response, status

from app.api.dependencies import (
    get_chat_service,
    get_ingest_validator,
    get_wiki_reader,
    get_wiki_writer,
    require_ingest_token,
)
from app.schemas.chat import ChatRequest, ChatResponse, TokenUsageDto
from app.schemas.ingest import IngestRequest, IngestResponse, IngestResultEntry
from app.schemas.report import GenerateReportRequest, ReportResponse
from app.services.chat_service import ChatService
from app.services.ingest_pipeline import IngestValidator, _parse_front_matter, extract_wikilinks
from app.services.llm_client import ChatMessage
from app.services.report_generator import ReportGenerator
from app.services.wiki_reader import WikiReader
from app.services.wiki_writer import WikiWriter

router = APIRouter(prefix="/v1", tags=["ai"])
_generator = ReportGenerator()


@router.post("/reports", response_model=ReportResponse, deprecated=True)
def generate_report(req: GenerateReportRequest, response: Response) -> ReportResponse:
    response.headers["Deprecation"] = "true"
    response.headers["Link"] = '</v1/chat>; rel="successor-version"'
    response.headers["Sunset"] = "Wed, 01 Jul 2026"
    return _generator.generate(req)


@router.post("/chat", response_model=ChatResponse)
async def chat(
    req: ChatRequest,
    service: ChatService = Depends(get_chat_service),
) -> ChatResponse:
    msgs = [ChatMessage(role=m.role, content=m.content) for m in req.messages]
    result, citations = await service.answer(msgs)
    return ChatResponse(
        answer=result.text,
        citations=citations,
        model=result.model,
        stop_reason=result.stop_reason,
        usage=TokenUsageDto(
            input_tokens=result.usage.input_tokens,
            output_tokens=result.usage.output_tokens,
            cache_read_input_tokens=result.usage.cache_read_input_tokens,
            cache_creation_input_tokens=result.usage.cache_creation_input_tokens,
        ),
    )


@router.get("/wiki/index")
def wiki_index(reader: WikiReader = Depends(get_wiki_reader)) -> dict[str, object]:
    return {"pages": reader.list_pages(), "index": reader.read_index()}


@router.get("/wiki/lint")
def wiki_lint(reader: WikiReader = Depends(get_wiki_reader)) -> dict[str, object]:
    issues: list[dict[str, str]] = []
    pages = reader.list_pages()
    page_set = {p.removesuffix(".md") for p in pages}
    for rel in pages:
        page = reader.read_page(rel)
        if page is None:
            continue
        if _parse_front_matter(page.body) is None:
            issues.append({"path": rel, "issue": "missing front-matter"})
        for link in extract_wikilinks(page.body):
            if link not in page_set:
                issues.append({"path": rel, "issue": f"broken wikilink: {link}"})
    return {"pages": len(pages), "issues": issues}


@router.post(
    "/wiki/ingest",
    response_model=IngestResponse,
    dependencies=[Depends(require_ingest_token)],
)
def wiki_ingest(
    req: IngestRequest,
    validator: IngestValidator = Depends(get_ingest_validator),
    writer: WikiWriter = Depends(get_wiki_writer),
) -> IngestResponse:
    errors = validator.validate(req.edits)
    if errors:
        raise HTTPException(
            status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
            detail=[{"path": e.path, "message": e.message} for e in errors],
        )
    results = writer.apply(req)
    return IngestResponse(
        applied=[
            IngestResultEntry(path=r.path, op=r.op, bytes_written=r.bytes_written) for r in results
        ],
        source=req.source.name,
    )
