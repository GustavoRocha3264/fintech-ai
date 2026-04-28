.PHONY: help backend-run backend-test frontend-run frontend-test ai-run ai-test engine-build engine-test compose-up compose-down proto

help:
	@echo "Targets:"
	@echo "  backend-run     run Go API on :8080"
	@echo "  backend-test    run Go tests"
	@echo "  frontend-run    run Vite dev server on :5173"
	@echo "  frontend-test   run Vitest"
	@echo "  ai-run          run FastAPI AI service on :8000"
	@echo "  ai-test         run pytest"
	@echo "  engine-build    cmake build the C++ engine"
	@echo "  engine-test     run GoogleTest suite"
	@echo "  compose-up      start full stack via docker compose"
	@echo "  compose-down    stop the stack"
	@echo "  proto           regenerate gRPC stubs (requires protoc)"

backend-run:
	cd backend-go && go run ./cmd/api

backend-test:
	cd backend-go && go test ./...

frontend-run:
	cd frontend && npm run dev

frontend-test:
	cd frontend && npm test

ai-run:
	cd ai-service && uvicorn app.main:app --reload --port 8000

ai-test:
	cd ai-service && pytest

engine-build:
	cmake -S cpp-engine -B cpp-engine/build && cmake --build cpp-engine/build

engine-test: engine-build
	ctest --test-dir cpp-engine/build --output-on-failure

compose-up:
	docker compose up --build

compose-down:
	docker compose down

proto:
	@echo "Regenerate gRPC stubs from proto/engine.proto — wire protoc + protoc-gen-go-grpc here."
