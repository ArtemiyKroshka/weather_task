# Makefile for Weather API Application

init-env: ## Copy .env.example to .env (only if .env is missing)
	@if [ ! -f .env ]; then cp .env.example .env && echo "Created .env from example"; else echo ".env already exists"; fi

up: ## Build and start containers (detached)
	docker compose up --build -d
	echo "Services are up"

down: ## Stop and remove containers
	docker compose down
	echo "Services are down"

test: ## Run Go unit tests
	go test ./...