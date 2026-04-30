.PHONY: help up down logs ps build keys clean

COMPOSE := docker compose

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

up: ## Build and start the full platform
	@$(MAKE) keys
	$(COMPOSE) up -d --build
	@echo
	@echo "  user-auth         : http://localhost:8080"
	@echo "  car-service       : http://localhost:8082"
	@echo "  booking           : http://localhost:8083"
	@echo "  currency-converter: http://localhost:8000"
	@echo "  RabbitMQ UI       : http://localhost:15672 (guest/guest)"

down: ## Stop the platform
	$(COMPOSE) down

logs: ## Tail logs from all services
	$(COMPOSE) logs -f --tail=100

ps: ## Show service status
	$(COMPOSE) ps

build: ## Rebuild images without starting
	$(COMPOSE) build

keys: ## Ensure JWT keys exist + public key copied to shared-secrets/
	@if [ ! -f user-authManagement/secrets/jwt_private.pem ]; then \
	  echo "Generating JWT keys via user-authManagement/scripts/gen-jwt-keys.sh"; \
	  cd user-authManagement && ./scripts/gen-jwt-keys.sh ./secrets; \
	fi
	@mkdir -p shared-secrets
	@cp -f user-authManagement/secrets/jwt_public.pem shared-secrets/jwt_public.pem
	@echo "shared-secrets/jwt_public.pem ready"

clean: down ## Stop platform and remove all volumes (DATA LOSS)
	$(COMPOSE) down -v
