"""Configuration loaded from environment variables.

Single source of truth for runtime config. Use `get_settings()` from anywhere —
the result is cached so settings are parsed once per process.
"""
from functools import lru_cache
from pathlib import Path

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
        extra="ignore",
    )

    # LLM provider
    anthropic_api_key: str = Field(default="", description="Anthropic API key — empty disables live calls")
    ai_model: str = Field(default="claude-opus-4-7", description="Model ID for chat completions")
    max_tokens: int = Field(default=2048, ge=1, le=128_000)

    # Wiki layout (Slice 2 wires the actual reads)
    wiki_dir: Path = Field(default=Path("wiki"), description="Path to the wiki directory on disk")

    # Operational
    log_level: str = Field(default="INFO")
    backend_url: str = Field(default="http://localhost:8080", description="Go backend base URL")


@lru_cache(maxsize=1)
def get_settings() -> Settings:
    return Settings()
