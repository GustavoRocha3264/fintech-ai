#!/usr/bin/env bash
# Convenience launcher for local development. Starts each service in the
# background; use Ctrl-C to tear them down.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

cleanup() { jobs -p | xargs -r kill 2>/dev/null || true; }
trap cleanup EXIT

(cd "$ROOT/ai-service"  && uvicorn app.main:app --port 8000) &
(cd "$ROOT/backend-go"  && go run ./cmd/api) &
(cd "$ROOT/frontend"    && npm run dev) &

wait
