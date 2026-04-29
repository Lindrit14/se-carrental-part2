# user-auth

User / Auth Microservice in Go. REST + RabbitMQ events, MongoDB persistence,
deployed to Azure Container Apps.

> Architecture & conventions are documented in [`CLAUDE.md`](./CLAUDE.md).
> API spec is the source of truth in [`api/openapi.yaml`](./api/openapi.yaml).

## Running

This service is part of the platform stack. The full stack (RabbitMQ +
Mongo + booking + currency-converter) is brought up from the platform root:

```bash
cd ..
make up           # builds + starts everything
```

This compiles the Docker image from `deploy/docker/Dockerfile` and wires
the service against the shared RabbitMQ + a dedicated Mongo container.
JWT keys are generated on first run into `secrets/`, and the public key
is copied to `../shared-secrets/` so the booking service can verify tokens.

After the stack is up:

| URL                     | What                          |
|-------------------------|-------------------------------|
| http://localhost:8080   | this service (user-auth)      |
| http://localhost:15672  | RabbitMQ Management (guest/guest) |

Smoke test:

```bash
curl -X POST localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"hunter2hunter2"}'

curl -X POST localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"hunter2hunter2"}'
```

Bind a temporary queue to the `user.events` exchange in the RabbitMQ UI to
observe `user.registered` and `user.login` payloads.

## Running standalone (for service-only dev)

If you want to iterate on this service without rebuilding the platform
image each time, run Mongo + RabbitMQ from the platform compose and the
Go binary directly on the host:

```bash
# from the platform root, start only the dependencies
cd ..
docker compose up -d rabbitmq user-auth-mongo

# then back here, run the binary
cd user-authManagement
make gen-keys                          # one-off: generates secrets/jwt_*.pem
cp .env.example .env                   # set MONGO_URI/RABBITMQ_URL to localhost
export $(cat .env | xargs)
export JWT_PRIVATE_KEY_FILE=$(pwd)/secrets/jwt_private.pem
export JWT_PUBLIC_KEY_FILE=$(pwd)/secrets/jwt_public.pem
make run
```

## Layout

See `CLAUDE.md` for the full architectural rules. In short:

```
cmd/server/                Composition root
internal/domain/           Entities + repo interfaces (no framework imports)
internal/application/      Use cases + outbound ports
internal/infrastructure/   Mongo, RabbitMQ, JWT, bcrypt, config
internal/interfaces/http/  chi router, middleware, handlers, DTOs
pkg/apperror/              Error → HTTP mapping
deploy/docker/Dockerfile   Built by the platform docker-compose
deploy/azure/              Bicep templates (production deployment)
api/openapi.yaml           API spec
```

## Common commands

```bash
make help              # list targets
make run               # run locally (needs Mongo + RabbitMQ on localhost)
make test              # unit tests
make test-integration  # integration tests (testcontainers)
make lint              # golangci-lint
make docker-build      # production image (usually built by platform compose)
make gen-keys          # generate RSA keypair for JWT
```

## Deployment

Azure Container Apps via Bicep — see [`deploy/azure/README.md`](./deploy/azure/README.md).
