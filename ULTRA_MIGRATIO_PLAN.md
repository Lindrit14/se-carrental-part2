# ULTRA PLAN: Migration to Full Microservices Architecture + CI/CD + Azure Deployment

I've reviewed your current 4-service setup (user-auth, booking, currency, frontend). This plan migrates you to **5 backend microservices + API Gateway**, addresses your critical architecture issues, and adds full CI/CD + Azure deployment.

The plan is split into **8 phases** so you can ship incrementally. Each phase ends in a deployable, working state — never break the system.

---

## Target Architecture

```
                          ┌────────────────────────────────────┐
                          │   Frontend (React, nginx)          │
                          │   - HttpOnly cookies (no localSt.) │
                          └─────────────────┬──────────────────┘
                                            │ HTTPS
                                            ▼
                          ┌────────────────────────────────────┐
                          │   API Gateway (Spring Cloud GW)    │
                          │   - JWT validation                 │
                          │   - Rate limiting (Redis)          │
                          │   - CORS, correlation IDs          │
                          │   - /api/v1/** routing             │
                          │   Port: 8080                       │
                          └─────────────────┬──────────────────┘
                                            │
   ┌──────────────┬───────────────┬─────────┼──────────┬──────────────┐
   ▼              ▼               ▼         ▼          ▼              ▼
┌──────────┐ ┌──────────┐ ┌────────────┐ ┌────────┐ ┌──────────┐
│ user-auth│ │   car    │ │  booking   │ │currency│ │notificat.│
│   (Go)   │ │  (Java)  │ │  (Java)    │ │(Python)│ │   (Go)   │
│  REST    │ │  REST    │ │ REST+gRPC  │ │  gRPC  │ │ Consumer │
│  :8081   │ │  :8082   │ │   :8083    │ │ :50051 │ │ :8084 hc │
└────┬─────┘ └────┬─────┘ └─────┬──────┘ └───┬────┘ └────┬─────┘
     │            │             │            │           │
     ▼            ▼             ▼            ▼           │
┌─────────┐  ┌──────────┐  ┌──────────┐  ┌───────┐       │
│ MongoDB │  │PostgreSQL│  │PostgreSQL│  │ Redis │       │
│users_db │  │ cars_db  │  │bookings  │  │ TTL   │       │
└─────────┘  └──────────┘  └──────────┘  └───────┘       │
     │            │             │            │           │
     └────────────┴─────────────┴────────────┴───────────┘
                                │
                       ┌────────▼─────────┐
                       │     RabbitMQ     │
                       │  carrental.evt   │
                       │  (topic exch.)   │
                       └──────────────────┘
```

### Key Changes from Current State

| Aspect | Current | Target |
|---|---|---|
| Services | 3 backend (auth, booking, currency) | **5 backend + Gateway** (split car-CRUD out of booking; add notification) |
| Frontend → backend | 3 different origins | **1 origin** (Gateway) |
| Token storage | localStorage (XSS-vulnerable) | **HttpOnly cookies + CSRF** |
| Booking → events | Direct publish (lossy) | **Transactional outbox** |
| gRPC resilience | None | **Circuit breaker + cached rate fallback** |
| Notifications | None | **Dedicated Go consumer service** |
| Pagination | Missing | **Cursor-based on all list endpoints** |
| Tracing | Missing | **W3C traceparent across HTTP, gRPC, RabbitMQ** |
| CI/CD | Only user-auth | **All services + frontend** |
| Cloud | Only user-auth on ACA | **Full stack on Azure Container Apps** |

---

## Phase 0 — Pre-Flight (1–2 days) ✅ DONE (2026-04-30)

**Goal:** Set up the workspace so all subsequent phases are unblockable.

### Tasks
1. **Create a `arch-change` long-lived branch** off main; you'll merge phase-by-phase.
2. **Document baseline**: capture current `docker-compose up` happy path with screenshots/curl scripts as your regression suite.
3. **Add structured logging baseline**:
   - Booking service: switch Spring default → **Logback JSON encoder** (`logstash-logback-encoder`)
   - Confirm user-auth, currency already JSON
4. **Add health endpoints everywhere** (`/health` on all HTTP services, gRPC health protocol on currency).
5. **Lock down toolchain versions** in a root `.tool-versions` or `versions.md` (Go 1.26, Java 21, Python 3.12, Node 20).
6. **Create `docs/adr/` folder** for Architecture Decision Records — write ADR-0001 capturing the migration goals.

### Exit Criteria
- All current services run via `docker-compose up` cleanly.
- All services emit structured JSON logs.
- Baseline regression curl script passes.

---

## Phase 1 — Split `car` out of `booking` Service (3–5 days) ✅ DONE (2026-04-30)

**Goal:** Achieve true 5-service split. The `car` aggregate becomes its own service.

### Why first?
The biggest structural change. Doing it before the gateway means the gateway only needs to be configured once with the final route map.

### Tasks
1. **Create `carBookingMicroservice/` → split into two**:
   - `carMicroservice/` (Java/Spring Boot, PostgreSQL `cars_db`)
   - `bookingMicroservice/` (Java/Spring Boot, PostgreSQL `bookings_db`)
2. **Move car CRUD, search, listing** to `carMicroservice`.
3. **Booking service becomes a client of car service** (REST for availability checks).
4. **Database split**:
   - Create new database `cars_db`, migrate `cars` table via Flyway in `carMicroservice`.
   - In `bookingMicroservice`, replace the `cars` table with a denormalized **car snapshot in the `bookings` row** (`car_brand`, `car_model`, `price_per_day_usd_at_booking`) — never query cars cross-service for historical bookings.
   - Keep the `customers` read model in `bookingMicroservice` (it's already a read model owned by booking).
5. **Define booking → car contract**:
   - `GET /api/v1/cars/{id}` for lookup at booking creation
   - `GET /api/v1/cars/{id}/availability?from=&to=` if you implement availability pre-check
6. **Update `docker-compose.yml`**: add `postgres-cars`, `car-service` (port 8082), shift booking to **8083**.
7. **Update frontend API client**: car endpoints → `:8082`, booking → `:8083`. (This is temporary; gateway in Phase 3 unifies it.)
8. **Tests**: copy car-related tests to new service; add a contract test for booking → car.

### Risks
- Data migration: if you have production data, write a one-shot SQL migration to copy `cars` → new DB. Otherwise, recreate.
- Foreign key removal: bookings.car_id is no longer a real FK; document this in ADR.

### Exit Criteria
- `carMicroservice` deployed and reachable on `:8082`.
- `bookingMicroservice` makes successful gRPC call to currency AND REST call to car service when creating a booking.
- All existing user flows still work end-to-end.

---

## Phase 2 — Notification Service (2–3 days) ✅ DONE (2026-04-30)

**Goal:** Add the 5th service. Pure RabbitMQ consumer in Go.

### Tasks
1. **Scaffold `notificationMicroservice/`** — Go 1.26, Clean Architecture (consistent with user-auth).
2. **RabbitMQ bindings**: queue `notifications.queue` bound to:
   - `user.events` exchange: `user.registered`, `user.password_reset_requested`
   - `booking.events` exchange: `booking.created`, `booking.cancelled`
3. **Notifier interface** with two implementations:
   - `MockNotifier` — logs to stdout (default for dev/grading).
   - `SmtpNotifier` — stub with real SMTP, optional via env flag.
4. **Manual ack** + **dead-letter exchange** `carrental.events.dlx` (after 3 failed redeliveries).
5. **Idempotency**: track processed `event_id`s in an in-memory LRU (size 10k) — sufficient for project scope, document as known limit.
6. **Health endpoint** on `:8084` checking RabbitMQ connection.

### Why now?
- Forces you to formalize the event envelope contract before more producers exist.
- Validates that user-auth's existing events are consumable by something other than booking.

### Exit Criteria
- Register a user → see `[NOTIFICATION] Welcome email sent to ...` in notification service logs.
- Create a booking → see `[NOTIFICATION] Booking confirmation ...`.
- Kill RabbitMQ → restart it → notification service reconnects automatically.

---

## Phase 3 — API Gateway (Spring Cloud Gateway) (3–4 days) ✅ DONE (2026-04-30)

**Goal:** Single entry point. Frontend talks to **one** URL.

### Tasks
1. **Scaffold `apiGateway/`** — Spring Boot 3.3 + Spring Cloud Gateway 2024.x (reactive).
2. **Routes** (in `application.yml`):
   ```yaml
   - id: user-auth
     uri: http://user-auth:8081
     predicates: [Path=/api/v1/auth/**, /api/v1/users/**]
   - id: cars
     uri: http://car-service:8082
     predicates: [Path=/api/v1/cars/**]
   - id: bookings
     uri: http://booking-service:8083
     predicates: [Path=/api/v1/bookings/**]
   ```
3. **Global filters**:
   - **JWT validation** using shared public key (RS256) — same key user-auth signs with. Add `aud` claim check (closes issue #11).
   - **CORS** (single origin configured here, remove from individual services).
   - **Rate limiting** via `RequestRateLimiter` + Redis (100 req/min/IP for general, 5/min for `/auth/login` and `/auth/register` — preserves user-auth's existing limits but centralized).
   - **Correlation ID filter**: generate `X-Correlation-Id` if absent, propagate downstream.
   - **JWT claim forwarding**: add `X-User-Id`, `X-User-Role` headers from validated JWT for downstream services.
4. **Bypass JWT** for: `/api/v1/auth/login`, `/api/v1/auth/register`, `/api/v1/auth/refresh`, `/health`, `/api/v1/cars` (public listing).
5. **Update frontend `lib/api/client.ts`**: base URL → `http://localhost:8080`. All path patterns stay the same.
6. **Remove CORS configs** from user-auth, booking, car, currency-REST.
7. **Update `docker-compose.yml`**: add `redis` (for rate limit counters) and `api-gateway` (port 8080).

### Why Spring Cloud Gateway over Nginx/Kong?
- You're already in the Spring ecosystem; minimal new tooling.
- JWT filter and Redis rate limiter are first-class features.
- Reactive = handles many concurrent connections cheaply.

### Exit Criteria
- Frontend works using only `http://localhost:8080` as the API origin.
- Hitting `:8081`, `:8082`, `:8083` directly from outside is blocked (or only for dev, document this).
- Rate limit triggers correctly: 6 logins in a minute → 429.
- Correlation ID appears in logs of every service for a single request.

---

## Phase 4 — Transactional Outbox + Event Persistence (3–4 days)

**Goal:** Fix the two highest-severity issues (#1 and #2).

### Tasks

#### 4a. Outbox in booking service
1. **Add `outbox` table** in `bookings_db` (Flyway migration):
   ```sql
   CREATE TABLE outbox (
     id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     aggregate_id UUID NOT NULL,
     event_type   VARCHAR(100) NOT NULL,
     payload      JSONB NOT NULL,
     occurred_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     published_at TIMESTAMPTZ,
     attempts     INT NOT NULL DEFAULT 0
   );
   CREATE INDEX idx_outbox_unpublished ON outbox(occurred_at) WHERE published_at IS NULL;
   ```
2. **In booking creation transaction**: `INSERT INTO bookings + INSERT INTO outbox` atomically. Remove direct RabbitMQ publish from the request thread.
3. **Outbox publisher**: Spring `@Scheduled` task (every 1s) — selects unpublished events with `FOR UPDATE SKIP LOCKED`, publishes to RabbitMQ with publisher confirms, marks `published_at` on success. Increment `attempts` on failure; alert (log ERROR) after 5 attempts.

#### 4b. Outbox in user-auth service
4. Same pattern in MongoDB: an `outbox` collection, atomic insert with the user document via session/transaction.
5. Background goroutine polls and publishes.

#### 4c. Inbox / Event store in booking
6. **Add `event_log` table** to persist all consumed `user.*` events before processing the customer read model. Replay capability: a CLI command `replay-events --from-timestamp=...` rebuilds the customer table.
7. Same for notification service consumed events (smaller scope; just an `event_log` collection optional).

### Exit Criteria
- Kill RabbitMQ mid-booking-creation → booking is saved → on RabbitMQ recovery, event auto-publishes.
- `event_log` table contains every consumed event; manual replay rebuilds `customers` table identically.

---

## Phase 5 — Resilience, Observability, Security Hardening (3–5 days)

**Goal:** Address mid-severity issues #3, #5, #6, #7, #9, #10, #11, #12.

### 5a. Circuit Breaker on gRPC (issue #5)
- Add **Resilience4j** to booking service.
- `@CircuitBreaker(name = "currency", fallbackMethod = "convertFallback")` around the gRPC client call.
- Fallback strategy: use last-known cached rate from a local Caffeine cache (TTL 1h); if no cache, reject booking with `503 Currency service unavailable`.
- Configure: 5 failures in 10s → open circuit for 30s.

### 5b. HttpOnly Cookies (issue #3)
- user-auth `/login` and `/refresh` set:
  - `Set-Cookie: access_token=...; HttpOnly; Secure; SameSite=Strict; Path=/api/v1; Max-Age=900`
  - `Set-Cookie: refresh_token=...; HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth/refresh; Max-Age=1209600`
- Add **CSRF token** endpoint `/api/v1/auth/csrf` returning a token; frontend stores in memory + sends as `X-CSRF-Token` header.
- Gateway validates CSRF on state-changing requests (POST/PUT/DELETE/PATCH) by comparing header to a server-side store (Redis).
- Frontend: remove `localStorage.setItem('access_token', ...)`. Cookie auto-sent. Use `credentials: 'include'` in fetch.

### 5c. Pagination (issue #6)
- `GET /api/v1/cars?cursor=...&limit=20`
- `GET /api/v1/bookings/me?cursor=...&limit=20`
- Cursor = base64-encoded `(id, created_at)` of last item; stable under inserts.
- Default limit 20, max 100.

### 5d. Rate limiting on booking (issue #7)
- Already covered by gateway global limit, but add stricter per-user limits at gateway: `POST /api/v1/bookings` → 10/min/user.

### 5e. Distributed Tracing (issue #9)
- Adopt **W3C Trace Context (`traceparent`)** as the standard.
- Each service: extract `traceparent` from inbound HTTP, gRPC metadata, RabbitMQ message headers; propagate to all outbound calls.
- Java: Spring Boot 3 has Micrometer Tracing built in — enable with OTLP exporter.
- Go: `go.opentelemetry.io/otel` SDK.
- Python: `opentelemetry-instrumentation-fastapi` + `opentelemetry-instrumentation-grpc`.
- Local dev: run **Jaeger all-in-one** in `docker-compose.yml`; UI on :16686.

### 5f. Prometheus Metrics (issue #12)
- Java: Micrometer + `micrometer-registry-prometheus` → `/actuator/prometheus`.
- Go: `github.com/prometheus/client_golang` → `/metrics` on a separate port (or same).
- Python: `prometheus-client` → `/metrics`.
- Add Prometheus container to compose; scrape configs.
- Define core metrics: HTTP request count/duration/error rate, RabbitMQ publish/consume rate, gRPC call duration, circuit breaker state.

### 5g. JWT `aud` claim (issue #11)
- user-auth issues JWTs with `aud: "carrental-api"`.
- Gateway and booking validate audience.

### 5h. Standardized error response (issue #8)
- Define `ProblemDetails` (RFC 7807) format across all services:
  ```json
  {
    "type": "https://carrental/errors/booking-conflict",
    "title": "Car not available",
    "status": 409,
    "detail": "Car is already booked for this date range",
    "instance": "/api/v1/bookings",
    "correlationId": "..."
  }
  ```
- Java: `@RestControllerAdvice` returning `ProblemDetail` (Spring 6 native).
- Go: shared `pkg/httperr` package writes the same shape.
- Python: FastAPI exception handler.

### Exit Criteria
- `docker-compose stop currency-converter` → bookings still succeed using cached rate (or fail with proper 503 + ProblemDetails).
- Network tab in browser: cookies are HttpOnly; no token in localStorage.
- Hitting any endpoint, you can trace the full request through Jaeger UI across all services.
- Hitting `/metrics` on every service returns Prometheus format.

---

## Phase 6 — CI/CD on GitHub Actions (4–6 days)

**Goal:** Every service has automated build, test, lint, scan, and deploy.

### Repository Layout Decision
**Recommendation: keep monorepo, use path-filtered workflows.** Avoids version skew and shared proto changes across multiple repos.

### Per-Service Workflow Pattern

For each service, create `.github/workflows/<service>-ci.yml` with **path filter** (`paths: ['carMicroservice/**']`).

#### Common Stages (every service)
```
lint → unit-test → integration-test (Testcontainers) → 
build-image → trivy-scan → push-to-ACR (on main) → 
deploy-to-ACA (on main, OIDC)
```

### Service-Specific Pipelines

#### 6a. user-auth (Go) — already exists, audit & enhance
- Add: golangci-lint, gosec, race detector, integration test stage.
- Already has OIDC → ACA deploy; reuse the pattern.

#### 6b. car-service & booking-service (Java)
```yaml
- mvn -B verify (unit + integration via Testcontainers)
- mvn jacoco:report → upload coverage
- sonarcloud scan
- docker build → trivy → push to ACR
- az containerapp update (OIDC)
```

#### 6c. currency-converter (Python)
```yaml
- ruff check + ruff format --check
- mypy src/
- pytest --cov
- bandit -r src/
- docker build → trivy → push → deploy
```

#### 6d. notification-service (Go)
- Mirror user-auth pipeline.

#### 6e. api-gateway (Java)
- Same as Java services.

#### 6f. frontend
```yaml
- npm ci
- npm run lint (eslint)
- npm run typecheck (tsc --noEmit)
- npm run test (vitest)
- npm run build
- docker build (nginx static) → push → deploy
```

### Cross-Cutting Workflows
1. **`pr-checks.yml`** — runs on every PR: changed-services-only matrix build, requires all green to merge.
2. **`security-scan.yml`** — scheduled weekly: `dependabot` config + Trivy on all images + `gitleaks` for secret scanning.
3. **`integration-e2e.yml`** — on merges to main: spins up full docker-compose, runs Postman/Newman or Playwright E2E suite.
4. **`release.yml`** — manual trigger: tags all services with same version, deploys to ACA.

### Secrets Management in GitHub
- **OIDC federated credentials** (no long-lived secrets) — create one Azure AD app registration, grant `acr-push` and `containerapps-contributor` on the resource group, federate to your repo.
- Repository secrets: `AZURE_CLIENT_ID`, `AZURE_TENANT_ID`, `AZURE_SUBSCRIPTION_ID`, `ACR_NAME`, `RESOURCE_GROUP`.
- Per-environment secrets via **GitHub Environments** (`dev`, `staging`, `prod`) with required reviewers on prod.

### Exit Criteria
- Push to `arch-change` → only changed services rebuild.
- Merge to `main` → all changed services deploy automatically.
- A failing test blocks merge.
- Trivy CRITICAL findings block deploy.

---

## Phase 7 — Azure Deployment for All Services (5–7 days)

**Goal:** Production-ready deployment of the full stack.

### Target: Azure Container Apps + Managed Dependencies

| Component | Azure Service | Tier (project) |
|---|---|---|
| All 6 services (gateway, user-auth, car, booking, currency, notification) | **Azure Container Apps** | Consumption |
| Frontend | **Azure Static Web Apps** (or Container App with nginx) | Free tier |
| MongoDB (`users_db`) | **Azure Cosmos DB for MongoDB** | Serverless |
| PostgreSQL (cars + bookings) | **Azure Database for PostgreSQL Flexible Server** | Burstable B1ms, 2 DBs on 1 server |
| Redis (rate limit + currency cache + CSRF) | **Azure Cache for Redis** | Basic C0 |
| RabbitMQ | **Self-hosted** as a Container App (Azure Service Bus would require code changes) | Single instance |
| Container Registry | **Azure Container Registry** | Basic |
| Secrets | **Azure Key Vault** | Standard |
| Observability | **Azure Monitor + Log Analytics + Application Insights** | Pay-as-you-go |
| TLS / Custom Domain | Container Apps managed certificates | Free |

### 7a. Infrastructure as Code
- **Terraform** (recommended) or **Bicep** in `infra/` folder.
- Modules: `network`, `databases`, `cache`, `messaging`, `containerapps`, `monitoring`.
- Two environments: `dev` and `prod`, separate resource groups, parameterized via `*.tfvars`.

### 7b. Network & Security
- All Container Apps in a single **Container Apps Environment** with internal-only addressing for backend services.
- **Only the API Gateway is publicly exposed** (Container Apps ingress = external).
- All other services: internal ingress only — they're reachable from the gateway via DNS (`car-service.internal.<env>.azurecontainerapps.io`).
- PostgreSQL & Cosmos & Redis: **VNet integration**, private endpoints, no public access.
- Key Vault: managed identity per Container App, granted `get/list` on secrets.

### 7c. Migration Strategy
1. **Dev environment first** — full deploy, verify everything works.
2. **Schema migrations**:
   - Java services run Flyway on startup (already configured).
   - user-auth Mongo: ensure indexes created on startup.
3. **DNS cutover**: point `api.<yourdomain>` → Gateway's Container App FQDN.
4. **Frontend** deploys to Static Web Apps with env var `VITE_API_BASE_URL=https://api.<yourdomain>`.

### 7d. Observability in Azure
- **Application Insights** auto-attached to each Container App (set `APPLICATIONINSIGHTS_CONNECTION_STRING`).
- Log Analytics workspace receives stdout from all containers.
- KQL dashboards: error rate per service, p95 latency per endpoint, RabbitMQ queue depth.
- Alert rules: gateway 5xx > 5%/5min, booking circuit breaker open, RabbitMQ consumer lag.

### 7e. Cost Controls
- Cosmos DB serverless: scales to zero usage.
- Container Apps: min replicas = 0 for non-gateway services, gateway min = 1 (avoid cold start on first request).
- Postgres B1ms ≈ €13/mo, Redis C0 ≈ €15/mo, ACR ≈ €5/mo, Static Web Apps free, Container Apps consumption-based.
- Total dev environment: **~€40–60/month** with idle scaling.

### 7f. Disaster Recovery (lightweight, project scope)
- Cosmos: continuous backup (30 days).
- Postgres: automated backups (7 days retention).
- IaC = entire infra rebuildable in <30min.

### Exit Criteria
- `https://api.<yourdomain>/health` returns 200 from all services via gateway aggregation.
- Frontend deployed to Static Web Apps fully functional.
- Push to main → CI/CD deploys to dev → manual approval → deploys to prod.
- Application Insights dashboards show traffic across all services.
- Cost dashboard shows expected ranges.

---

## Phase Summary & Timeline

| Phase | Topic | Duration | Risk | Reversible? |
|---|---|---|---|---|
| 0 | Pre-flight | 1–2 d | Low | N/A |
| 1 | Split car out of booking | 3–5 d | **High** (data migration) | Difficult |
| 2 | Notification service | 2–3 d | Low | Easy |
| 3 | API Gateway | 3–4 d | Medium | Easy |
| 4 | Outbox + event store | 3–4 d | Medium | Easy |
| 5 | Resilience + observability + security | 3–5 d | Medium | Per-feature |
| 6 | CI/CD | 4–6 d | Low | N/A |
| 7 | Azure full deployment | 5–7 d | Medium | Easy (IaC) |

**Total: ~24–36 working days** (5–7 weeks if focused).

---

## Critical Path Dependencies

```
Phase 0 ──► Phase 1 ──► Phase 3 ──► Phase 5 ──► Phase 7
              │           │
              ▼           ▼
           Phase 2 ──► Phase 4
                          │
                          ▼
                       Phase 6 (parallel with 5)
```

You can run **Phase 6 (CI/CD) in parallel with Phase 5** once Phase 4 is in.

---

## Recommendations Before Starting

A few things I'd want confirmed before generating any actual code or detailed task breakdowns:

1. **Repository structure**: stay monorepo, or split? (My strong recommendation: monorepo with path-filtered CI.)
2. **Production traffic**: is this purely an academic project, or will it have real users? (Affects the rigor needed in Phase 4 and 5b.)
3. **Existing data**: any production data in your current Postgres/Mongo that needs migrating in Phase 1, or is it disposable dev data?
4. **Azure subscription**: do you have student credits / Azure for Students? (Affects the tier choices in Phase 7.)
5. **Time budget**: realistic weekly hours? (Lets me adjust phase scoping.)

Want me to **expand any single phase into a day-by-day task breakdown with specific files to touch and commit messages**? Phase 1 or Phase 4 are the highest-leverage to detail first because they carry the most risk.
