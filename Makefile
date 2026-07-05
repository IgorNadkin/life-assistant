.PHONY: help build run test db-up db-down migrate migrate-status clean deps

help:
	@echo "Life Assistant - Available commands:"
	@echo ""
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run in development mode (with hot reload)"
	@echo "  make test           - Run tests"
	@echo "  make db-up          - Start PostgreSQL with Docker Compose"
	@echo "  make db-down        - Stop PostgreSQL"
	@echo "  make db-logs        - Show database logs"
	@echo "  make migrate        - Run database migrations (idempotent, tracked)"
	@echo "  make migrate-status - Show which migrations are applied/pending"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Download dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run in Docker"
	@echo ""

build:
	@echo "Building application..."
	cd backend && go build -o bin/api cmd/api/main.go
	@echo "Build complete: backend/bin/api"

run: build
	@echo "Starting server..."
	cd backend && ./bin/api

dev:
	@echo "Running in development mode..."
	cd backend && go run cmd/api/main.go

test:
	@echo "Running tests..."
	cd backend && go test -v ./...

db-up:
	@echo "Starting PostgreSQL..."
	cd backend && docker-compose up -d
	@echo "PostgreSQL is running on localhost:5432"

db-down:
	@echo "Stopping PostgreSQL..."
	cd backend && docker-compose down

db-logs:
	cd backend && docker-compose logs -f postgres

migrate:
	@bash scripts/migrate.sh

migrate-status:
	@bash scripts/migrate.sh status

clean:
	@echo "Cleaning build artifacts..."
	cd backend && rm -rf bin/
	cd backend && go clean

fmt:
	@echo "Formatting code..."
	cd backend && go fmt ./...

lint:
	@echo "Running linter..."
	cd backend && golangci-lint run ./... || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

deps:
	@echo "Downloading dependencies..."
	cd backend && go mod download
	cd backend && go mod tidy

docker-build:
	@echo "Building Docker image..."
	docker build -f backend/Dockerfile -t life-assistant:latest .
	@echo "Image built: life-assistant:latest"

docker-run: docker-build
	@echo "Running application in Docker..."
	docker-compose -f backend/docker-compose.yml up

postman:
	@echo "Postman collection available at: postman_collection.json"
	@echo "Import into Postman for API testing"
	chmod +x scripts/test_api.sh
	./scripts/test_api.sh

.DEFAULT_GOAL := help
