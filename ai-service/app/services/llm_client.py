"""LLM client abstraction.

Use cases depend on the `LLMClient` protocol; concrete implementations live
behind it. Today the only implementation is Anthropic; swapping providers
later means writing a new class — call sites do not change.
"""
from __future__ import annotations

from dataclasses import dataclass
from typing import Literal, Protocol

import anthropic
from anthropic import AsyncAnthropic
from anthropic.types import MessageParam


Role = Literal["user", "assistant"]


@dataclass(frozen=True)
class ChatMessage:
    """Domain-shaped chat turn — keeps the API surface provider-neutral."""

    role: Role
    content: str


@dataclass(frozen=True)
class TokenUsage:
    """Token accounting per request, including cache attribution."""

    input_tokens: int
    output_tokens: int
    cache_read_input_tokens: int = 0
    cache_creation_input_tokens: int = 0


@dataclass(frozen=True)
class CompletionResult:
    text: str
    usage: TokenUsage
    model: str
    stop_reason: str


class LLMClient(Protocol):
    async def complete(
        self,
        messages: list[ChatMessage],
        *,
        system: str,
    ) -> CompletionResult: ...


class AnthropicLLMClient:
    """Anthropic-backed LLM client.

    Notable choices:
    - System prompt is sent as a cacheable text block (`cache_control: ephemeral`).
      Today the prompt is small so caching won't kick in (Opus 4.7 needs ~4K
      tokens of prefix), but Slice 2 will inject the wiki index/schema here and
      caching will become highly impactful.
    - No sampling parameters — Opus 4.7 removed `temperature`, `top_p`, `top_k`.
    - Adaptive thinking is left off by default for chat; the query pipeline in
      Slice 3 will opt into it for harder questions.
    """

    def __init__(
        self,
        *,
        api_key: str,
        model: str,
        max_tokens: int,
    ) -> None:
        if not api_key:
            raise ValueError("AnthropicLLMClient requires a non-empty api_key")
        self._client = AsyncAnthropic(api_key=api_key)
        self._model = model
        self._max_tokens = max_tokens

    async def complete(
        self,
        messages: list[ChatMessage],
        *,
        system: str,
    ) -> CompletionResult:
        # Map domain messages to the SDK shape. Trust internal callers.
        api_messages: list[MessageParam] = [
            {"role": m.role, "content": m.content} for m in messages
        ]

        # Cacheable system block; the marker is harmless when the prefix is
        # too short to cache, and pays off once the wiki schema lands.
        system_blocks = [
            {
                "type": "text",
                "text": system,
                "cache_control": {"type": "ephemeral"},
            }
        ]

        response = await self._client.messages.create(
            model=self._model,
            max_tokens=self._max_tokens,
            system=system_blocks,
            messages=api_messages,
        )

        text_parts = [block.text for block in response.content if block.type == "text"]
        usage = response.usage

        return CompletionResult(
            text="".join(text_parts),
            usage=TokenUsage(
                input_tokens=usage.input_tokens,
                output_tokens=usage.output_tokens,
                cache_read_input_tokens=getattr(usage, "cache_read_input_tokens", 0) or 0,
                cache_creation_input_tokens=getattr(usage, "cache_creation_input_tokens", 0) or 0,
            ),
            model=response.model,
            stop_reason=response.stop_reason or "unknown",
        )


__all__ = [
    "ChatMessage",
    "CompletionResult",
    "LLMClient",
    "AnthropicLLMClient",
    "TokenUsage",
    "anthropic",
]
