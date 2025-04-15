APP = mcpserver
MAIN = cmd/api/main.go
MIGRATE = cmd/migrate/main.go
OUT = bin
ENV_FILE = configs/api/.env
TEMP_ENV_FILE = configs/temp/.env
FRONTEND_DIR = frontend

.PHONY: up
up: build-image
	@echo "Starting all services (backend + dependencies)..."
	docker compose up -d
	@echo "Starting frontend development server..."
	cd $(FRONTEND_DIR) && npm run dev

.PHONY: down
down:
	@echo "Stopping all services (backend + dependencies)..."
	docker compose down

.PHONY: build-image
build-image:
	@echo "Building Docker image for backend and dependencies..."
	docker compose build

.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	cp $(ENV_FILE) $(TEMP_ENV_FILE)
	go run $(MIGRATE)
	rm -f $(TEMP_ENV_FILE)

.PHONY: frontend-install
frontend-install:
	@echo "Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install

.PHONY: frontend-dev
frontend-dev:
	@echo "Starting Next.js frontend development server..."
	cd $(FRONTEND_DIR) && npm run dev

.PHONY: frontend-build
frontend-build:
	@echo "Building Next.js frontend for static export..."
	cd $(FRONTEND_DIR) && npm run build

.PHONY: frontend-lint
frontend-lint:
	@echo "Linting frontend code..."
	cd $(FRONTEND_DIR) && npm run lint

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  --- Backend --- "
	@echo "  make migrate         - Run database migrations"
	@echo "  --- Frontend (Next.js) --- "
	@echo "  make frontend-install - Install frontend dependencies (npm install)"
	@echo "  make frontend-dev    - Start Next.js development server (npm run dev)"
	@echo "  make frontend-build  - Build Next.js static export (npm run build)"
	@echo "  make frontend-lint   - Lint frontend code (npm run lint)"
	@echo "  --- Docker --- "
	@echo "  make up              - Start all backend services with Docker and frontend dev server"
	@echo "  make down            - Stop all backend services in Docker"
	@echo "  make build-image     - Build Docker images for backend services"
	@echo "  --- Combined Development --- "
	@echo "  make help            - Show this help message"
