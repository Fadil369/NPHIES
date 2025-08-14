# NPHIES Platform Makefile

# Variables
DOCKER_COMPOSE = docker compose
DOCKER_COMPOSE_DEV = $(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.dev.yml
GO_VERSION = 1.21
JAVA_VERSION = 17

# Colors for terminal output
GREEN = \033[0;32m
YELLOW = \033[0;33m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: help dev-up dev-down build test lint clean docker-build docker-push

help: ## Show this help message
	@echo "NPHIES Platform Development Commands"
	@echo "======================================"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development Environment
dev-up: ## Start local development environment
	@echo "$(GREEN)Starting NPHIES development environment...$(NC)"
	$(DOCKER_COMPOSE_DEV) up -d
	@echo "$(GREEN)Development environment started!$(NC)"
	@echo "Services available at:"
	@echo "  - API Gateway: http://localhost:8080"
	@echo "  - Kafka UI: http://localhost:8081"
	@echo "  - Redis Commander: http://localhost:8082"
	@echo "  - PostgreSQL: localhost:5432"

dev-down: ## Stop local development environment
	@echo "$(YELLOW)Stopping NPHIES development environment...$(NC)"
	$(DOCKER_COMPOSE_DEV) down
	@echo "$(GREEN)Development environment stopped!$(NC)"

dev-logs: ## Show logs for development environment
	$(DOCKER_COMPOSE_DEV) logs -f

# Build Commands
build: build-go build-java ## Build all services

build-go: ## Build Go services
	@echo "$(GREEN)Building Go services...$(NC)"
	@cd services/api-gateway && go mod tidy && go build -o bin/api-gateway ./cmd/main.go
	@cd services/eligibility-service && go mod tidy && go build -o bin/eligibility-service ./cmd/main.go
	@cd services/terminology-service && go mod tidy && go build -o bin/terminology-service ./cmd/main.go
	@cd services/wallet-service && go mod tidy && go build -o bin/wallet-service ./cmd/main.go
	@if [ -d "services/analytics-service" ]; then \
		cd services/analytics-service && go mod tidy && go build -o bin/analytics-service ./cmd/main.go; \
	fi
	@echo "$(GREEN)Go services built successfully!$(NC)"

build-java: ## Build Java services
	@echo "$(GREEN)Building Java services...$(NC)"
	@if [ -d "services/claims-service" ]; then \
		cd services/claims-service && ./mvnw clean package -DskipTests; \
	fi
	@echo "$(GREEN)Java services built successfully!$(NC)"

# Testing
test: test-go test-java ## Run all tests

test-go: ## Run Go tests
	@echo "$(GREEN)Running Go tests...$(NC)"
	@cd services/api-gateway && go test ./... -v
	@cd services/eligibility-service && go test ./... -v
	@cd services/terminology-service && go test ./... -v
	@cd services/wallet-service && go test ./... -v
	@if [ -d "services/analytics-service" ]; then \
		cd services/analytics-service && go test ./... -v; \
	fi

test-java: ## Run Java tests
	@echo "$(GREEN)Running Java tests...$(NC)"
	@if [ -d "services/claims-service" ]; then \
		cd services/claims-service && ./mvnw test; \
	fi

test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	@cd tests/integration && go test ./... -v

# Code Quality
lint: lint-go lint-java ## Run linters for all services

lint-go: ## Run Go linter
	@echo "$(GREEN)Running Go linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed$(NC)" && exit 1)
	@cd services/api-gateway && golangci-lint run
	@cd services/eligibility-service && golangci-lint run
	@cd services/terminology-service && golangci-lint run

lint-java: ## Run Java linter
	@echo "$(GREEN)Running Java linter...$(NC)"
	@if [ -d "services/claims-service" ]; then \
		cd services/claims-service && ./mvnw checkstyle:check; \
	fi

# Docker Commands
docker-build: ## Build Docker images for all services
	@echo "$(GREEN)Building Docker images...$(NC)"
	docker build -t nphies/api-gateway:latest services/api-gateway/
	docker build -t nphies/eligibility-service:latest services/eligibility-service/
	docker build -t nphies/terminology-service:latest services/terminology-service/
	@if [ -d "services/claims-service" ]; then \
		docker build -t nphies/claims-service:latest services/claims-service/; \
	fi
	@echo "$(GREEN)Docker images built successfully!$(NC)"

docker-push: docker-build ## Push Docker images to registry
	@echo "$(GREEN)Pushing Docker images...$(NC)"
	docker push nphies/api-gateway:latest
	docker push nphies/eligibility-service:latest
	docker push nphies/terminology-service:latest
	@if [ -d "services/claims-service" ]; then \
		docker push nphies/claims-service:latest; \
	fi

# Database
db-migrate: ## Run database migrations
	@echo "$(GREEN)Running database migrations...$(NC)"
	@cd scripts/database && ./migrate.sh

db-seed: ## Seed database with test data
	@echo "$(GREEN)Seeding database...$(NC)"
	@cd scripts/database && ./seed.sh

# Cleanup
clean: ## Clean build artifacts and dependencies
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@find . -name "bin" -type d -exec rm -rf {} + 2>/dev/null || true
	@find . -name "target" -type d -exec rm -rf {} + 2>/dev/null || true
	@find . -name "node_modules" -type d -exec rm -rf {} + 2>/dev/null || true
	@echo "$(GREEN)Cleanup completed!$(NC)"

clean-docker: ## Clean Docker resources
	@echo "$(YELLOW)Cleaning Docker resources...$(NC)"
	docker system prune -f
	docker volume prune -f

# Infrastructure
infra-plan: ## Plan infrastructure changes
	@echo "$(GREEN)Planning infrastructure changes...$(NC)"
	@cd infrastructure/terraform && terraform plan

infra-apply: ## Apply infrastructure changes
	@echo "$(GREEN)Applying infrastructure changes...$(NC)"
	@cd infrastructure/terraform && terraform apply

# Monitoring
logs: ## Show service logs
	kubectl logs -f -l app=nphies --tail=100

health-check: ## Check service health
	@echo "$(GREEN)Checking service health...$(NC)"
	@curl -f http://localhost:8080/health || echo "$(RED)API Gateway unhealthy$(NC)"
	@curl -f http://localhost:8090/health || echo "$(RED)Eligibility Service unhealthy$(NC)"
	@curl -f http://localhost:8091/health || echo "$(RED)Terminology Service unhealthy$(NC)"

# Development Tools
setup-dev: ## Setup development environment
	@echo "$(GREEN)Setting up development environment...$(NC)"
	@scripts/setup/dev-setup.sh

setup-hooks: ## Setup git hooks
	@echo "$(GREEN)Setting up git hooks...$(NC)"
	@cp scripts/hooks/* .git/hooks/
	@chmod +x .git/hooks/*

# Load Testing
load-test: ## Run load tests
	@echo "$(GREEN)Running load tests...$(NC)"
	@cd tests/load && ./run-load-tests.sh

# Security
security-scan: ## Run security scans
	@echo "$(GREEN)Running security scans...$(NC)"
	@scripts/security/scan.sh

# Generate API documentation
docs-api: ## Generate API documentation
	@echo "$(GREEN)Generating API documentation...$(NC)"
	@cd services/api-gateway && swag init
	@cd services/eligibility-service && swag init
	@cd services/terminology-service && swag init