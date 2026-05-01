.PHONY: help up down logs ps build keys clean

COMPOSE := docker compose

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

up: ## Build and start the full platform
	@$(MAKE) keys
	$(COMPOSE) up -d --build
	@echo
	@echo "  api-gateway       : http://localhost:8080  (single entry point)"
	@echo "  frontend          : http://localhost:3000"
	@echo "  RabbitMQ UI       : http://localhost:15672 (guest/guest)"
	@echo
	@echo "  Backends are internal-only (reachable via the gateway):"
	@echo "    user-auth, car-service, booking, currency-converter, notification"

down: ## Stop the platform
	$(COMPOSE) down

logs: ## Tail logs from all services
	$(COMPOSE) logs -f --tail=100

ps: ## Show service status
	$(COMPOSE) ps

build: ## Rebuild images without starting
	$(COMPOSE) build

keys: ## Ensure JWT keypair exists in shared-secrets/
	@if [ ! -f shared-secrets/jwt_private.pem ]; then \
	  echo "Generating JWT keys via scripts/gen-jwt-keys.sh"; \
	  ./scripts/gen-jwt-keys.sh shared-secrets; \
	fi
	@echo "shared-secrets/jwt_{private,public}.pem ready"

clean: down ## Stop platform and remove all volumes (DATA LOSS)
	$(COMPOSE) down -v
