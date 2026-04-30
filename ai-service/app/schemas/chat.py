"""Pydantic schemas for the /v1/chat endpoint."""
from __future__ import annotations

from typing import Literal

from pydantic import BaseModel, Field


Role = Literal["user", "assistant"]


class ChatMessageDto(BaseModel):
    role: Role
    content: str = Field(min_length=1)


class ChatRequest(BaseModel):
    conversation_id: str | None = None
    messages: list[ChatMessageDto] = Field(min_length=1)


class CitationDto(BaseModel):
    """Citation pointer — populated by the query pipeline in a later slice."""

    slug: str
    title: str
    anchor: str | None = None


class UsageDto(BaseModel):
    input_tokens: int
    output_tokens: int
    cache_read_input_tokens: int = 0
    cache_creation_input_tokens: int = 0


class ChatResponse(BaseModel):
    answer: str
    citations: list[CitationDto] = Field(default_factory=list)
    usage: UsageDto
    model: str
    stop_reason: str
