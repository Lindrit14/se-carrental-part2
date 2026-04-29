.PHONY: help run build test test-integration lint tidy fmt vet \
        docker-build gen-keys clean

BINARY := server
PKG    := ./cmd/server
IMAGE  := user-auth:dev

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

run: ## Run server locally (expects Mongo + RabbitMQ on localhost; see ../docker-compose.yml)
	go run $(PKG)

build: ## Build binary into ./bin
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/$(BINARY) $(PKG)

test: ## Run unit tests
	go test ./... -short -race -count=1

test-integration: ## Run integration tests (testcontainers)
	go test ./test/integration/... -race -count=1 -timeout 5m

lint: ## Run golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || { \
	  echo "golangci-lint not found. Install: brew install golangci-lint"; exit 1; }
	golangci-lint run ./...

tidy: ## go mod tidy
	go mod tidy

fmt: ## Format code
	gofmt -s -w .

vet: ## go vet
	go vet ./...

docker-build: ## Build production image (manual; usually built by platform compose)
	docker build -f deploy/docker/Dockerfile -t $(IMAGE) .

gen-keys: ## Generate RSA keypair for JWT (RS256)
	@./scripts/gen-jwt-keys.sh

clean: ## Remove build artifacts
	rm -rf bin dist coverage.*
