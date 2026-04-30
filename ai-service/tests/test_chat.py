"""Slice 1 contract tests for /v1/chat.

Two layers:
- Unit/contract: stub LLMClient via FastAPI dependency override. No network.
- Integration: hits real Anthropic API, gated on ANTHROPIC_API_KEY.
"""
from __future__ import annotations

import os

import pytest
from fastapi.testclient import TestClient

from app.api.dependencies import get_chat_service
from app.main import app
from app.schemas.chat import ChatMessageDto, ChatResponse, UsageDto
from app.services.chat_service import ChatService
from app.services.llm_client import (
    ChatMessage,
    CompletionResult,
    LLMClient,
    TokenUsage,
)


class _FakeLLM:
    """In-process LLMClient used by the contract test."""

    def __init__(self) -> None:
        self.calls: list[tuple[str, list[ChatMessage]]] = []

    async def complete(
        self, messages: list[ChatMessage], *, system: str
    ) -> CompletionResult:
        self.calls.append((system, list(messages)))
        return CompletionResult(
            text="resposta de teste",
            usage=TokenUsage(input_tokens=42, output_tokens=7, cache_read_input_tokens=0),
            model="fake-model",
            stop_reason="end_turn",
        )


def _override_with(fake: LLMClient):
    return lambda: ChatService(fake)


def test_chat_contract_returns_answer_and_usage() -> None:
    fake = _FakeLLM()
    app.dependency_overrides[get_chat_service] = _override_with(fake)
    try:
        client = TestClient(app)
        resp = client.post(
            "/v1/chat",
            json={
                "conversation_id": "c1",
                "messages": [{"role": "user", "content": "O que é PFIC?"}],
            },
        )
    finally:
        app.dependency_overrides.clear()

    assert resp.status_code == 200
    body = ChatResponse.model_validate(resp.json())
    assert body.answer == "resposta de teste"
    assert body.citations == []
    assert body.usage == UsageDto(
        input_tokens=42, output_tokens=7, cache_read_input_tokens=0
    )
    assert body.model == "fake-model"
    assert body.stop_reason == "end_turn"

    # Verify the system prompt was actually applied (not empty / not silently swapped).
    assert len(fake.calls) == 1
    system_used, messages_seen = fake.calls[0]
    assert "CBPI" in system_used
    assert messages_seen == [ChatMessage(role="user", content="O que é PFIC?")]


def test_chat_rejects_empty_messages() -> None:
    client = TestClient(app)
    resp = client.post("/v1/chat", json={"messages": []})
    assert resp.status_code == 422  # pydantic validation


def test_chat_rejects_invalid_role() -> None:
    client = TestClient(app)
    resp = client.post(
        "/v1/chat",
        json={"messages": [{"role": "system", "content": "hi"}]},
    )
    assert resp.status_code == 422


@pytest.mark.integration
@pytest.mark.skipif(
    not os.getenv("ANTHROPIC_API_KEY"),
    reason="ANTHROPIC_API_KEY not set; skipping live API call",
)
def test_chat_against_real_anthropic_api() -> None:
    """Smoke test: real provider returns a non-empty answer.

    Run explicitly with: pytest -m integration
    Skipped in normal CI to keep the fast suite hermetic.
    """
    # Drop the autouse stub so the real dependency chain wires up.
    app.dependency_overrides.clear()
    client = TestClient(app)
    resp = client.post(
        "/v1/chat",
        json={
            "messages": [
                {"role": "user", "content": "Diga 'olá' em uma palavra."}
            ]
        },
    )
    assert resp.status_code == 200, resp.text
    body = resp.json()
    assert isinstance(body["answer"], str)
    assert len(body["answer"]) > 0
    assert body["usage"]["output_tokens"] > 0
