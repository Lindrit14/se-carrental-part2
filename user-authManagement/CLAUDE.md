# CLAUDE.md

Dieses Dokument briefiet Claude Code über das Repo. Es ist absichtlich kurz —
ausführliche Doku gehört in `README.md` / `docs/`.

## Projekt-Überblick

User/Auth-Microservice in **Go (1.23+)**. REST-API für Authentifizierung
und Profil-Management, Domain-Events über **RabbitMQ**. Persistenz in
**MongoDB**. Deployment: **Docker → Azure Container Apps**.

Der Service ist Teil einer Microservice-Landschaft. Er publisht Events,
auf die andere Services reagieren — er kennt diese anderen Services nicht.

## Tech-Stack (verbindlich)

| Bereich     | Wahl                                                   |
|-------------|--------------------------------------------------------|
| HTTP/Router | `github.com/go-chi/chi/v5` + `net/http`                |
| Auth        | JWT (RS256), Access 15 min / Refresh 14 d, Rotation   |
| Hashing     | `golang.org/x/crypto/bcrypt` (Cost 12)                 |
| DB          | MongoDB, Driver `go.mongodb.org/mongo-driver/v2`       |
| Messaging   | RabbitMQ, Client `github.com/rabbitmq/amqp091-go`      |
| Logging     | `log/slog` (Stdlib), JSON-Output                       |
| Config      | ENV via `github.com/kelseyhightower/envconfig`         |
| Validation  | `github.com/go-playground/validator/v10`               |
| Tests       | `testing` + `testify` + `testcontainers-go`            |
| Container   | Multi-Stage Dockerfile, Final-Image `distroless/static`|

Keine alternativen Bibliotheken einführen ohne Diskussion.

## Architektur — Clean Architecture, 4 Layer

Imports zeigen **immer nach innen**:

```
interfaces/http  ─┐
                  ├──►  application  ──►  domain
infrastructure  ──┘
```

1. `internal/domain/`         — Entities + Repository-Interfaces.
                                 **Keine Framework-Imports**, nur Stdlib.
2. `internal/application/`    — Use Cases. Importiert nur `domain` + eigene `ports/`.
3. `internal/infrastructure/` — Mongo, RabbitMQ, JWT, Bcrypt, Config.
                                 Implementiert Interfaces aus `domain` + `ports`.
4. `internal/interfaces/http/`— chi-Router, Handler, Middleware, DTOs.

**Dependency Injection ausschließlich in `cmd/server/main.go`.**
Keine Service-Locator, keine globalen Singletons, kein `init()` für Side-Effects.

## Verzeichnis-Layout

```
cmd/server/                    # Composition Root (main.go)
internal/
  domain/{user,token}/         # Entities, Repo-Interfaces, Sentinel-Errors
  application/
    ports/                     # Outbound-Interfaces (Hasher, TokenService, EventPublisher, Clock)
    auth/                      # Register, Login, Refresh, Logout, PasswordReset
    user/                      # GetProfile, UpdateProfile, DeleteAccount
  infrastructure/
    config/                    # ENV-Loading
    mongodb/                   # Client + Repo-Implementations
    rabbitmq/                  # Connection, Publisher, Event-Definitionen
    crypto/                    # bcrypt_hasher.go, jwt_service.go
    logging/                   # slog Setup
  interfaces/http/
    server.go, router.go
    middleware/                # auth, logging, recovery, request_id, ratelimit
    handlers/                  # auth_handler.go, user_handler.go, health_handler.go
    dto/                       # Request- + Response-DTOs
pkg/apperror/                  # Typisierte Fehler → HTTP-Status-Mapping
api/openapi.yaml               # API-Spec (Source of Truth)
deploy/docker/                 # Dockerfile, docker-compose.yml
deploy/azure/                  # *.bicep, Deploy-Doku
test/{integration,e2e}/        # Tests die echte Mongo/RabbitMQ brauchen
```

## RabbitMQ — Event-Contract

- **Exchange**: `user.events` (`topic`, durable)
- **Routing-Keys & Payloads** (alle JSON, Felder: `event_id`, `event_type`, `version`, `occurred_at`, `data`):
  - `user.registered` — `{ user_id, email }`
  - `user.updated` — `{ user_id, changed_fields }`
  - `user.deleted` — `{ user_id }`
  - `user.password_reset_requested` — `{ user_id, email, reset_token, expires_at }`
  - `user.login` — `{ user_id, ip, user_agent }`

Routing-Keys + Payload-Structs leben als Konstanten/Typen in
`internal/infrastructure/rabbitmq/events.go` — Konsumenten dürfen sich darauf verlassen.
Publisher Confirms sind aktiv, Reconnect mit Exponential Backoff.

## HTTP-Endpunkte

| Methode | Pfad                            | Auth |
|---------|---------------------------------|------|
| POST    | /api/v1/auth/register           | —    |
| POST    | /api/v1/auth/login              | —    |
| POST    | /api/v1/auth/refresh            | —    |
| POST    | /api/v1/auth/logout             | JWT  |
| POST    | /api/v1/auth/password/reset     | —    |
| POST    | /api/v1/auth/password/confirm   | —    |
| GET     | /api/v1/users/me                | JWT  |
| PATCH   | /api/v1/users/me                | JWT  |
| DELETE  | /api/v1/users/me                | JWT  |
| GET     | /healthz                        | —    |
| GET     | /readyz                         | —    |

Vollständige Spec: `api/openapi.yaml`. Nicht doppelt dokumentieren.

## Konventionen

- **Logger**: `log/slog` JSON. Jede Log-Line hat `request_id` (aus Middleware-Context)
  und falls vorhanden `user_id`. Niemals `fmt.Println`, `log.Printf`.
- **Fehler**: Sentinel-Errors in `domain/<aggregate>/errors.go`
  (`ErrUserNotFound`, `ErrEmailTaken`, ...). HTTP-Mapping zentral in `pkg/apperror`.
  Use-Cases geben Domain-Fehler zurück; Handler übersetzen in HTTP-Status.
- **Konfiguration**: nur in `internal/infrastructure/config`. Außerhalb davon
  kein direktes `os.Getenv`. Defaults nur dort, niemals hardcoded im Business-Code.
- **DTOs**: Request-/Response-Typen in `interfaces/http/dto/`. **Domain-Entities
  verlassen den HTTP-Layer nicht** — immer in DTOs übersetzen.
- **Validation**: Auf DTO-Ebene mit Validator-Tags. Domain-Layer geht von
  validen Eingaben aus.
- **Context**: `context.Context` ist erstes Argument jeder Use-Case-Methode
  und jedes Repo-Calls. Timeout wird in der Middleware gesetzt.
- **Tests**:
  - Use-Cases: Unit-Tests mit Mocks (Repo + Ports).
  - Repos: Integration-Tests via `testcontainers-go` gegen echte Mongo / RabbitMQ.
  - HTTP-Handler: `httptest` + `chi.Router`.
  - **Keine** Mocks für Mongo selbst — wir mocken nur die Repo-Interfaces.

## Häufige Befehle

Der gesamte Stack wird vom Platform-Compose im Repo-Root gestartet
(`cd .. && make up`). Service-lokal:

| Befehl                    | Zweck                                          |
|---------------------------|------------------------------------------------|
| `make run`                | Binary lokal starten (erwartet Mongo + RabbitMQ auf localhost) |
| `make test`               | Unit-Tests                                     |
| `make test-integration`   | Integration-Tests (testcontainers)             |
| `make lint`               | golangci-lint                                  |
| `make docker-build`       | Multi-Stage Image bauen (sonst über Platform-Compose) |
| `make gen-keys`           | RSA-Keypair für JWT (RS256) erzeugen           |

## Sicherheit

- Passwort-Mindestlänge **12 Zeichen**, Bcrypt-Cost **12**.
- JWT **RS256**, Private Key kommt nur aus ENV / Azure Key Vault — nicht aus Datei im Image.
- Refresh-Tokens werden gehasht in Mongo gespeichert (TTL-Index auf `expires_at`).
- Rate-Limit auf `/login` und `/register` (5/min/IP).
- E-Mails werden vor Speicherung lowercased + getrimmt; Unique-Index auf `email`.
- CORS nur für konfigurierte Origins (ENV `CORS_ALLOWED_ORIGINS`).
- Security-Header: HSTS, X-Content-Type-Options, X-Frame-Options.

## Deployment

Ziel ist **Azure Container Apps**:

- Image wird in **Azure Container Registry (ACR)** gepushed (Tag = Git-SHA).
- Bicep-Templates in `deploy/azure/`.
- Secrets (`MONGO_URI`, `JWT_PRIVATE_KEY`, `RABBITMQ_URL`) als ACA Secret-References,
  optional aus Azure Key Vault.
- KEDA-Scale-Rule auf HTTP-Concurrency, Min 1 / Max 5 Replicas.
- MongoDB in Prod: Azure Cosmos DB for MongoDB API empfohlen.
- RabbitMQ self-hosted als eigene Container App mit Persistent Storage (Azure Files).

## Was NICHT tun

- Keine direkten Mongo-/AMQP-Calls aus Handlern oder Use-Cases — nur über
  Repo- bzw. Port-Interfaces.
- Keine `panic`s in Library-Code — Fehler zurückgeben.
- Keine `fmt.Println` / `log.Printf` — immer `slog`.
- Keine Secrets ins Repo, auch nicht als Default-Werte.
- Domain-Entities nicht mit BSON- oder JSON-Tags annotieren — Mapping
  passiert im Repository (Mongo) bzw. im DTO (HTTP).
- Keine Business-Logik in Handlern — Handler nur DTO ↔ Use-Case-Übersetzung.
- Keine zirkulären Layer-Imports. Wenn `domain` etwas von `infrastructure`
  zu brauchen scheint: Interface in `domain` definieren, in `infrastructure`
  implementieren.
