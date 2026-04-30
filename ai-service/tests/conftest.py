"""Test fixtures.

By default the chat service is wired to a fake LLM so the suite runs offline
and pydantic-validation tests don't trip on the missing-API-key 503. The
integration test in `test_chat.py` removes the override and exercises the real
dependency chain.
"""
from __future__ import annotations

import pytest

from app.api.dependencies import get_chat_service
from app.main import app
from app.services.chat_service import ChatService
from app.services.llm_client import (
    ChatMessage,
    CompletionResult,
    TokenUsage,
)


class _DefaultFakeLLM:
    async def complete(
        self, messages: list[ChatMessage], *, system: str
    ) -> CompletionResult:
        return CompletionResult(
            text="ok",
            usage=TokenUsage(input_tokens=1, output_tokens=1),
            model="fake",
            stop_reason="end_turn",
        )


@pytest.fixture(autouse=True)
def _stub_chat_service():
    app.dependency_overrides[get_chat_service] = lambda: ChatService(_DefaultFakeLLM())
    yield
    app.dependency_overrides.clear()
