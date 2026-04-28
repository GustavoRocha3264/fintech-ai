# Cross-Border Portfolio Intelligence Platform (CBPI)

Skeleton of a multi-currency (BRL/USD) portfolio intelligence platform. The
foundation is in place: bounded contexts, interface seams, and per-service
build/test scripts. Business logic is intentionally minimal so the architecture
stays the focus.

## Architecture

```
+-----------+        REST        +------------+      gRPC       +-------------+
| frontend  | <----------------> | backend-go | <-------------> | cpp-engine  |
| (React)   |                    | (Gin/DDD)  |                 | (sim/risk)  |
+-----------+                    +------------+                 +-------------+
                                       |
                                       | HTTP
                                       v
                                 +------------+
                                 | ai-service |
                                 | (FastAPI)  |
                                 +------------+
```

- **backend-go** — DDD/Clean Architecture API. Owns domain rules; delegates
  heavy compute to `cpp-engine` and narrative generation to `ai-service`.
- **cpp-engine** — pure C++20 compute kernels (Monte Carlo, risk, optimizer)
  exposed via gRPC; a C-ABI facade is included for FFI consumers.
- **frontend** — React + TypeScript + Vite + Zustand. Minimal dashboard.
- **ai-service** — FastAPI microservice that produces daily portfolio reports.

### Bounded contexts (Go backend)

- `portfolio` — `Portfolio`, `Position`, `Asset`
- `currency`  — `ExchangeRate`, `CurrencyPair`
- `market`    — `MarketData`, `AssetType`
- `analysis`  — `RiskMetrics`, `AllocationSuggestion`, `Report`
- `shared`    — value objects: `Money` (currency-aware), `Allocation`

### Key interfaces

| Interface              | Where                                                       |
|------------------------|-------------------------------------------------------------|
| `PortfolioRepository`  | `internal/domain/portfolio/repository.go`                   |
| `MarketDataProvider`   | `internal/domain/market/market_data.go`                     |
| `FXRateProvider`       | `internal/domain/currency/exchange_rate.go`                 |
| `AnalysisServiceClient`| `internal/domain/analysis/report.go`                        |
| `EngineClient`         | `internal/application/analysis/service.go`                  |
| C++ kernels            | `cpp-engine/include/engine/engine.h` + `engine_c_api.h`     |

## API (mock data)

| Method | Path                              | Purpose                        |
|--------|-----------------------------------|--------------------------------|
| POST   | `/api/v1/portfolio`               | Create a portfolio             |
| GET    | `/api/v1/portfolio/{id}`          | Fetch a portfolio              |
| POST   | `/api/v1/analysis/run`            | Trigger analysis run           |
| GET    | `/api/v1/analysis/{portfolioId}`  | Fetch latest report            |
| GET    | `/healthz`                        | Liveness                       |

## Folder structure

```
fintech/
├── backend-go/
│   ├── cmd/api/                  # main.go — composition root
│   ├── internal/
│   │   ├── domain/               # entities, value objects, repo interfaces
│   │   ├── application/          # use cases / app services
│   │   ├── infrastructure/       # adapters: persistence, gRPC, AI, FX, auth, messaging
│   │   └── interfaces/http/      # Gin router, handlers, DTOs, middleware
│   ├── proto/                    # local proto copy (mirrors /proto)
│   └── configs/
├── cpp-engine/
│   ├── include/engine/           # public headers (C++ + C-ABI)
│   ├── src/                      # simulator, risk, optimizer, server entry
│   ├── tests/                    # GoogleTest suite (FetchContent)
│   └── CMakeLists.txt
├── frontend/
│   └── src/{components,pages,services,store,types}/
├── ai-service/
│   └── app/{api,core,services,schemas}/
├── proto/                        # canonical gRPC contracts
├── scripts/
├── docker-compose.yml
└── Makefile
```

## Setup

Prerequisites: Go 1.22+, Node 20+, Python 3.11+, CMake 3.20+, Docker.

```bash
# Backend (Go)
make backend-run          # http://localhost:8080
make backend-test

# Frontend (React + Vite)
cd frontend && npm install
make frontend-run         # http://localhost:5173
make frontend-test

# AI service (FastAPI)
cd ai-service && pip install -e ".[dev]"
make ai-run               # http://localhost:8000
make ai-test

# C++ engine
make engine-build
make engine-test

# Full stack
make compose-up
```

## What's deliberately stubbed

- Trading / order logic
- gRPC server in `cpp-engine` (proto contract is defined, server wiring is not)
- JWT verification (`internal/infrastructure/auth`)
- NATS/Kafka adapters (`internal/infrastructure/messaging` — interfaces only)
- Real LLM / model provider in `ai-service`
- Persistent storage — repo is in-memory

## Engineering principles

Separation of concerns · dependency inversion · interfaces first ·
performance-aware boundaries (compute lives in C++, orchestration in Go,
narrative in Python).
