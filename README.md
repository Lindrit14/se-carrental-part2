# Microservices Platform

Three-service Uni project: authentication, car booking, currency conversion.
Mixed inter-service communication: **RabbitMQ** for asynchronous events,
**gRPC** for synchronous request/response. Frontend ↔ services is REST + JWT.

## Services

| Service | Stack | Ports | Repo dir |
|---|---|---|---|
| user-auth | Go + MongoDB | 8080 (HTTP) | `user-authManagement/` |
| booking | Java 21 / Spring Boot 3 + PostgreSQL | 8082 (HTTP) | `carBookingMicroservice/` |
| currency-converter | Python 3.12 / FastAPI | 8000 (HTTP), 9000 (gRPC) | `currencyConverterMicroservice/` |
| RabbitMQ | 3.13 management | 5672 / 15672 | (image) |

## Communication

```
                            Frontend
                               │ REST + JWT
       ┌───────────────────────┼─────────────────────────┐
       ▼                       ▼                         ▼
   user-auth                booking ───── gRPC ────► currency-converter
       │                       │                         (Convert)
       │ user.events           │
       │ (publish)             │ booking.events
       └──────────┬────────────┴──────────┐
                  ▼                       ▼
                       RabbitMQ
                  (topic exchanges)
```

- `user.*` events (RabbitMQ): user-auth publishes → booking consumes (creates/updates/anonymizes Customer).
- `currency.v1.CurrencyConverter/Convert` (gRPC): booking calls currency-converter synchronously.
- `booking.*` events (RabbitMQ): booking publishes (booking.created, booking.cancelled).
- currency-converter has no AMQP connection — it's a gRPC server only (plus its public HTTP API for the frontend).

## Quickstart

```bash
make up
```

That:

1. Generates the RSA JWT keypair into `shared-secrets/` (if missing) via `scripts/gen-jwt-keys.sh`. The public key is consumed by every service that verifies JWTs; only user-auth reads the private key.
2. Builds and starts all containers.

After it's healthy:

```bash
# Health
curl localhost:8080/healthz   # user-auth
curl localhost:8082/healthz   # booking
curl localhost:8000/healthz   # currency-converter

# Currency conversion (no auth)
curl -X POST localhost:8000/api/v1/convert \
  -H 'Content-Type: application/json' \
  -d '{"amount":"100.00","from":"USD","to":"EUR"}'

# Register + login (user-auth)
curl -X POST localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"a@b.de","password":"hunter2hunter2"}'

TOKENS=$(curl -s -X POST localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"a@b.de","password":"hunter2hunter2"}')
ACCESS=$(echo "$TOKENS" | python3 -c 'import json,sys;print(json.load(sys.stdin)["access_token"])')

# Cars (booking, JWT-protected)
curl -H "Authorization: Bearer $ACCESS" localhost:8082/api/v1/cars
```

## Common targets

```bash
make help     # list all
make up       # build + start everything
make down     # stop platform
make logs     # tail all logs
make ps       # service status
make clean    # stop + remove volumes (data loss)
```

## Architecture notes

- Each service follows **Clean Architecture** (domain → application → infrastructure → interfaces).
- Each service owns its own database — no shared schemas.
- JWTs are RS256-signed by user-auth with `shared-secrets/jwt_private.pem`.
  Other services verify locally with `shared-secrets/jwt_public.pem` — no
  network call to user-auth needed (stateless verification).
- Booking holds a **read model** of users (`Customer` table), populated and
  kept in sync by consuming `user.*` events from RabbitMQ.

For per-service details, see each service's own README and CLAUDE.md.
