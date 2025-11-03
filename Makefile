.PHONY: help build run test clean docker-build docker-up docker-down docker-logs

# Variables
APP_NAME=auto-message-dispatcher
DOCKER_COMPOSE=docker-compose
GO=go
BINARY_DIR=bin
BINARY_NAME=server

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	$(GO) build -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Build complete: $(BINARY_DIR)/$(BINARY_NAME)"

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	$(GO) run ./cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	$(GO) test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BINARY_DIR)
	$(GO) clean

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	$(GO) mod tidy

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-up: ## Start all services with Docker Compose
	@echo "Starting Docker Compose services..."
	$(DOCKER_COMPOSE) up -d
	@echo "Services started. View logs with 'make docker-logs'"

docker-down: ## Stop all services
	@echo "Stopping Docker Compose services..."
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker Compose logs
	$(DOCKER_COMPOSE) logs -f app

docker-restart: docker-down docker-up ## Restart all services

docker-clean: ## Remove all containers, volumes, and images
	@echo "Cleaning Docker resources..."
	$(DOCKER_COMPOSE) down -v
	docker rmi $(APP_NAME):latest 2>/dev/null || true

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	psql -h localhost -U postgres -d message_dispatcher -f migrations/001_create_messages_table.sql

dev: docker-up ## Start development environment
	@echo "Development environment started!"
	@echo "API available at: http://localhost:8080"
	@echo "Health check: http://localhost:8080/health"

.DEFAULT_GOAL := help

