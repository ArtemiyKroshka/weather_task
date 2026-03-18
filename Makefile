.PHONY: init-env up down test lint build

init-env: ## Copy .env.example to .env (only if .env is missing)
	@if [ ! -f .env ]; then cp .env.example .env && echo "Created .env from example"; else echo ".env already exists"; fi

build: ## Build the server binary
	go build -o server ./cmd/server

up: ## Build and start containers (detached)
	docker compose up --build -d
	@echo "Services are up"

down: ## Stop and remove containers
	docker compose down
	@echo "Services are down"

test: ## Run all tests with race detector
	go test -race ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	golangci-lint run ./...
