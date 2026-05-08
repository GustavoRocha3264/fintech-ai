from __future__ import annotations

from pathlib import Path

import pytest

from app.api import dependencies as deps
from app.core.config import get_settings
from app.main import app
from app.services.chat_service import ChatService
from app.services.llm_client import ChatMessage, CompletionResult, TokenUsage
from app.services.wiki_reader import WikiReader


class FakeLLM:
    def __init__(self, *, reply: str = "ok") -> None:
        self.reply = reply
        self.captured_system: str | None = None
        self.captured_messages: list[ChatMessage] | None = None

    async def complete(self, messages, *, system):
        self.captured_messages = list(messages)
        self.captured_system = system
        return CompletionResult(
            text=self.reply,
            model="fake-model",
            stop_reason="end_turn",
            usage=TokenUsage(input_tokens=1, output_tokens=1),
        )


@pytest.fixture
def wiki_root() -> Path:
    return Path(__file__).resolve().parent.parent / "wiki"


@pytest.fixture
def fake_llm() -> FakeLLM:
    return FakeLLM()


@pytest.fixture(autouse=True)
def _override_chat_service(fake_llm, wiki_root):
    """Replace get_chat_service in every test so /v1/chat doesn't need a real
    Anthropic key. The integration test must call app.dependency_overrides.clear()
    before exercising the real client.
    """
    def _factory() -> ChatService:
        return ChatService(llm=fake_llm, wiki=WikiReader(wiki_root))

    app.dependency_overrides[deps.get_chat_service] = _factory
    yield
    app.dependency_overrides.pop(deps.get_chat_service, None)
    get_settings.cache_clear()
