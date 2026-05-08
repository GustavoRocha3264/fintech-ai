from typing import Literal

from pydantic import BaseModel, Field

Role = Literal["user", "assistant"]


class ChatMessageDto(BaseModel):
    role: Role
    content: str = Field(min_length=1)


class ChatRequest(BaseModel):
    conversation_id: str | None = None
    messages: list[ChatMessageDto] = Field(min_length=1)


class TokenUsageDto(BaseModel):
    input_tokens: int = 0
    output_tokens: int = 0
    cache_read_input_tokens: int = 0
    cache_creation_input_tokens: int = 0


class ChatResponse(BaseModel):
    answer: str
    citations: list[str] = []
    model: str
    stop_reason: str | None = None
    usage: TokenUsageDto
