# CBPI AI Service

FastAPI microservice that hosts the CBPI assistant: a wiki-grounded chat
endpoint backed by Claude (Anthropic SDK), plus the receiving infrastructure
for a weekly Anthropic-driven ingest routine that feeds the wiki.

## Layout

```
app/
  api/            FastAPI routes, dependencies
  core/           settings, wiki path helpers
  schemas/        Pydantic DTOs (chat, ingest, report)
  services/       llm_client, chat_service, wiki_reader, wiki_writer, ingest_pipeline
  main.py         app factory
wiki/             on-disk knowledge base (Karpathy LLM Wiki layout)
  CLAUDE.md       schema + ingest contract (read me first)
  index.md        catalog of pages
  log.md          append-only ingest log
  entities/  concepts/  sources/  syntheses/  playbooks/
scripts/wiki_check.py   standalone linter (used by CI)
tests/                  pytest suite (no Anthropic key required)
```

## Setup

```bash
cd ai-service
pip install -e ".[dev]"
```

## Environment

| Var                  | Purpose                                           | Required for                |
|----------------------|---------------------------------------------------|-----------------------------|
| `ANTHROPIC_API_KEY`  | Anthropic API auth                                | `/v1/chat`                  |
| `AI_MODEL`           | Model id (default `claude-opus-4-7`)              | optional                    |
| `MAX_TOKENS`         | Response cap (default 2048)                       | optional                    |
| `WIKI_DIR`           | Wiki root (default `<repo>/ai-service/wiki`)      | optional                    |
| `WIKI_INGEST_TOKEN`  | Bearer token gating `POST /v1/wiki/ingest`        | weekly routine              |

## Running

```bash
uvicorn app.main:app --reload --port 8000
```

Smoke checks:

```bash
curl localhost:8000/healthz
curl localhost:8000/readyz
curl localhost:8000/v1/wiki/index | jq .pages
python scripts/wiki_check.py
```

## Endpoints

| Method | Path                  | Notes                                                       |
|--------|-----------------------|-------------------------------------------------------------|
| POST   | `/v1/chat`            | Wiki-grounded chat, requires `ANTHROPIC_API_KEY`.           |
| GET    | `/v1/wiki/index`      | Page catalog + raw `index.md`.                              |
| GET    | `/v1/wiki/lint`       | Same checks as `scripts/wiki_check.py`, JSON output.        |
| POST   | `/v1/wiki/ingest`     | **Bearer token required.** Atomic batch apply.              |
| POST   | `/v1/reports`         | Deprecated. Sunset 2026-07-01. Use `/v1/chat`.              |
| GET    | `/healthz`, `/readyz` | Liveness / readiness.                                       |

## The weekly Anthropic ingest routine

The wiki is updated once per week by an Anthropic Managed Agent configured
**outside this repo**. This service only owns the receiving contract.

The agent must:

1. Search the public web for cross-border investing material relevant to a
   Brazilian investor (BRL/USD exposure, IRRF, PFIC, ETFs on B3, etc.).
2. Draft new `sources/` pages and updates to `concepts/` / `syntheses/`,
   following the schema in [`wiki/CLAUDE.md`](wiki/CLAUDE.md). Front-matter
   (`type`, `status: draft`, `last_reviewed`, `sources`) is mandatory.
3. POST the batch to `/v1/wiki/ingest` with header
   `Authorization: Bearer $WIKI_INGEST_TOKEN`. Payload shape lives in
   `app/schemas/ingest.py::IngestRequest`.

Hard rules enforced by the receiver:

- All edits go through `IngestValidator` before any disk write — path must be
  inside `wiki/`, must be markdown, must declare front-matter.
- `status: stable` pages cannot be overwritten without `force=true`.
  The agent **must not** write `status: stable` itself; promotion to stable
  is human-only.
- Apply is atomic: tempdir-stage → rename swap → per-file `.bak` rollback on
  any failure. No partial state can land on disk.
- Every batch appends one line to `wiki/log.md`.

### Suggested system prompt for the agent

> Você é o curador semanal da wiki CBPI. Procure conteúdo público recente
> sobre investimento internacional para um investidor brasileiro (exposição
> cambial BRL/USD, PFIC, IRRF, ETFs na B3). Para cada item relevante, gere
> uma página `sources/<slug>.md` em PT-BR com front-matter completo
> (`type: source`, `status: draft`, `last_reviewed`, lista `sources` com a
> URL original). Atualize páginas existentes em `concepts/` ou `syntheses/`
> apenas quando houver fato novo verificável. Nunca defina `status: stable`.
> Submeta tudo numa única chamada a `POST /v1/wiki/ingest`.

### Suggested tool list for the agent

- Web search / fetch.
- HTTP POST with bearer auth (one call per week to `/v1/wiki/ingest`).
- A read-only `GET /v1/wiki/index` call at the start of each run so it knows
  what already exists.

## Tests

```bash
pytest                    # full suite
pytest -m "not integration"   # skip the live-API test
ANTHROPIC_API_KEY=sk-... pytest -m integration   # exercise the real client
```

The autouse fixture in `tests/conftest.py` swaps the chat dependency for a
fake LLM, so the suite runs hermetically. The integration test (when added)
must call `app.dependency_overrides.clear()` first.
