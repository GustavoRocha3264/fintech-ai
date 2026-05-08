from functools import lru_cache
from pathlib import Path

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", case_sensitive=False, extra="ignore")

    anthropic_api_key: str = Field(default="")
    ai_model: str = Field(default="claude-opus-4-7")
    max_tokens: int = Field(default=2048, ge=1, le=128_000)

    wiki_dir: Path = Field(default=Path(__file__).resolve().parent.parent.parent / "wiki")
    wiki_ingest_token: str = Field(default="")

    log_level: str = Field(default="INFO")
    backend_url: str = Field(default="http://localhost:8080")


@lru_cache(maxsize=1)
def get_settings() -> Settings:
    return Settings()
