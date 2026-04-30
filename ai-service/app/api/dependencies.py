"""FastAPI dependency wiring.

Keeps `routes.py` declarative. Tests override these via
`app.dependency_overrides[get_chat_service] = ...`.
"""
from __future__ import annotations

from functools import lru_cache

from fastapi import Depends, HTTPException, status

from app.core.config import Settings, get_settings
from app.services.chat_service import ChatService
from app.services.llm_client import AnthropicLLMClient, LLMClient


@lru_cache(maxsize=1)
def _build_llm_client(api_key: str, model: str, max_tokens: int) -> LLMClient:
    return AnthropicLLMClient(api_key=api_key, model=model, max_tokens=max_tokens)


def get_llm_client(settings: Settings = Depends(get_settings)) -> LLMClient:
    if not settings.anthropic_api_key:
        # Fail loudly rather than emit silent stub responses — Slice 1's
        # contract is "real LLM behind /v1/chat".
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="ANTHROPIC_API_KEY not configured",
        )
    return _build_llm_client(
        api_key=settings.anthropic_api_key,
        model=settings.ai_model,
        max_tokens=settings.max_tokens,
    )


def get_chat_service(llm: LLMClient = Depends(get_llm_client)) -> ChatService:
    return ChatService(llm)
