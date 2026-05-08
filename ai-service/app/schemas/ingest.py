from typing import Literal

from pydantic import BaseModel, Field


class WikiEdit(BaseModel):
    op: Literal["create", "update"]
    path: str = Field(
        min_length=1,
        description="Wiki-relative path, e.g. 'entities/USD.md'.",
    )
    content: str = Field(min_length=1, description="Full markdown body, including front-matter.")
    force: bool = Field(
        default=False,
        description="Allow overwriting a page whose front-matter status is 'stable'.",
    )


class IngestSource(BaseModel):
    name: str = Field(min_length=1, description="Provenance label, e.g. 'anthropic-weekly-2026-W18'.")
    fetched_at: str = Field(min_length=1, description="ISO-8601 timestamp.")
    notes: str = ""


class IngestRequest(BaseModel):
    source: IngestSource
    edits: list[WikiEdit] = Field(min_length=1)
    log_entry: str = Field(
        min_length=1,
        description="One-line operator log message; appended to wiki/log.md.",
    )


class IngestResultEntry(BaseModel):
    path: str
    op: Literal["create", "update"]
    bytes_written: int


class IngestResponse(BaseModel):
    applied: list[IngestResultEntry]
    source: str
