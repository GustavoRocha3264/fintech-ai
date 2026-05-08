from dataclasses import dataclass
from typing import Literal, Protocol

from anthropic import AsyncAnthropic

Role = Literal["user", "assistant"]


@dataclass(frozen=True)
class ChatMessage:
    role: Role
    content: str


@dataclass(frozen=True)
class TokenUsage:
    input_tokens: int
    output_tokens: int
    cache_read_input_tokens: int = 0
    cache_creation_input_tokens: int = 0


@dataclass(frozen=True)
class CompletionResult:
    text: str
    model: str
    stop_reason: str | None
    usage: TokenUsage


class LLMClient(Protocol):
    async def complete(
        self,
        messages: list[ChatMessage],
        *,
        system: str,
    ) -> CompletionResult: ...


class AnthropicLLMClient:
    """Async Anthropic client. System prompt is sent as a cached block so that
    long wiki context doesn't get re-billed on every turn (5-min ephemeral TTL).
    """

    def __init__(self, *, api_key: str, model: str, max_tokens: int) -> None:
        self._client = AsyncAnthropic(api_key=api_key)
        self._model = model
        self._max_tokens = max_tokens

    async def complete(
        self,
        messages: list[ChatMessage],
        *,
        system: str,
    ) -> CompletionResult:
        api_messages = [{"role": m.role, "content": m.content} for m in messages]
        system_blocks = [
            {"type": "text", "text": system, "cache_control": {"type": "ephemeral"}}
        ]
        response = await self._client.messages.create(
            model=self._model,
            max_tokens=self._max_tokens,
            system=system_blocks,
            messages=api_messages,
        )
        text = "".join(
            block.text for block in response.content if getattr(block, "type", None) == "text"
        )
        usage = response.usage
        return CompletionResult(
            text=text,
            model=response.model,
            stop_reason=response.stop_reason,
            usage=TokenUsage(
                input_tokens=getattr(usage, "input_tokens", 0),
                output_tokens=getattr(usage, "output_tokens", 0),
                cache_read_input_tokens=getattr(usage, "cache_read_input_tokens", 0) or 0,
                cache_creation_input_tokens=getattr(usage, "cache_creation_input_tokens", 0) or 0,
            ),
        )
