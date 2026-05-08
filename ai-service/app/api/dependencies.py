from functools import lru_cache

from fastapi import Depends, Header, HTTPException, status

from app.core.config import Settings, get_settings
from app.services.chat_service import ChatService
from app.services.ingest_pipeline import IngestValidator
from app.services.llm_client import AnthropicLLMClient, LLMClient
from app.services.wiki_reader import WikiReader
from app.services.wiki_writer import WikiWriter


@lru_cache(maxsize=1)
def _build_llm_client(api_key: str, model: str, max_tokens: int) -> LLMClient:
    return AnthropicLLMClient(api_key=api_key, model=model, max_tokens=max_tokens)


def get_llm_client(settings: Settings = Depends(get_settings)) -> LLMClient:
    if not settings.anthropic_api_key:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="ANTHROPIC_API_KEY is not configured",
        )
    return _build_llm_client(settings.anthropic_api_key, settings.ai_model, settings.max_tokens)


def get_wiki_reader(settings: Settings = Depends(get_settings)) -> WikiReader:
    return WikiReader(settings.wiki_dir)


def get_wiki_writer(settings: Settings = Depends(get_settings)) -> WikiWriter:
    return WikiWriter(wiki_root=settings.wiki_dir)


def get_ingest_validator(settings: Settings = Depends(get_settings)) -> IngestValidator:
    return IngestValidator(wiki_root=settings.wiki_dir)


def get_chat_service(
    llm: LLMClient = Depends(get_llm_client),
    wiki: WikiReader = Depends(get_wiki_reader),
) -> ChatService:
    return ChatService(llm=llm, wiki=wiki)


def require_ingest_token(
    settings: Settings = Depends(get_settings),
    authorization: str | None = Header(default=None),
) -> None:
    if not settings.wiki_ingest_token:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="WIKI_INGEST_TOKEN is not configured",
        )
    expected = f"Bearer {settings.wiki_ingest_token}"
    if authorization != expected:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="invalid or missing bearer token",
        )
