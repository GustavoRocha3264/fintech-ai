from fastapi import FastAPI
from app.api.routes import router

app = FastAPI(title="CBPI AI Service", version="0.0.1")
app.include_router(router)


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok"}
