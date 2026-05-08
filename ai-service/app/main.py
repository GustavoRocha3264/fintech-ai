from fastapi import Depends, FastAPI

from app.api.routes import router
from app.core.config import Settings, get_settings

app = FastAPI(title="CBPI AI Service", version="0.0.2")
app.include_router(router)


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok"}


@app.get("/readyz")
def readyz(settings: Settings = Depends(get_settings)) -> dict[str, object]:
    return {
        "status": "ok",
        "model": settings.ai_model,
        "anthropic_configured": bool(settings.anthropic_api_key),
        "wiki_dir": str(settings.wiki_dir),
        "ingest_configured": bool(settings.wiki_ingest_token),
    }
