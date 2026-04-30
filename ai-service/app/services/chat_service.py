"""Chat use case.

Slice 1 keeps the system prompt static. Slice 3 (`feat(ai)/query-pipeline`)
will replace `_DEFAULT_SYSTEM` with the wiki index and per-question page
selection logic.
"""
from __future__ import annotations

from app.schemas.chat import ChatMessageDto, ChatResponse, UsageDto
from app.services.llm_client import ChatMessage, LLMClient


_DEFAULT_SYSTEM = """\
You are CBPI, an assistant specialized in cross-border investing for users \
holding portfolios in BRL and USD. Answer in Portuguese unless the user \
writes in another language. Be precise about jurisdictions: distinguish \
Brazilian (Receita Federal, B3, IRRF) and US (IRS, NYSE/NASDAQ, PFIC) rules \
when they apply. Always include a brief disclaimer that responses are \
informational, not personalized investment advice. Do not invent regulations, \
tickers, or rates — if you do not know, say so.\
"""


class ChatService:
    def __init__(self, llm: LLMClient) -> None:
        self._llm = llm

    async def reply(self, messages: list[ChatMessageDto]) -> ChatResponse:
        domain_messages = [ChatMessage(role=m.role, content=m.content) for m in messages]
        result = await self._llm.complete(domain_messages, system=_DEFAULT_SYSTEM)

        return ChatResponse(
            answer=result.text,
            citations=[],  # Populated in Slice 3 once the wiki query pipeline lands.
            usage=UsageDto(
                input_tokens=result.usage.input_tokens,
                output_tokens=result.usage.output_tokens,
                cache_read_input_tokens=result.usage.cache_read_input_tokens,
                cache_creation_input_tokens=result.usage.cache_creation_input_tokens,
            ),
            model=result.model,
            stop_reason=result.stop_reason,
        )
