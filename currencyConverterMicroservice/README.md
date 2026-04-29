# currency-converter

Stateless currency conversion microservice. Pulls daily rates from the
European Central Bank XML feed, exposes them via REST (public, for the
frontend) and gRPC (internal, called by other services).

## Endpoints

REST (public, port 8000):

- `GET  /healthz` — liveness
- `GET  /api/v1/rates` — current rate set (EUR-based)
- `POST /api/v1/convert` — convert an amount

gRPC (internal, port 9000):

- Service `currency.v1.CurrencyConverter` — see `proto/currency/v1/currency.proto`.
- One unary RPC: `Convert(ConvertRequest) returns (ConvertResponse)`.
- Errors map to gRPC status codes:
  | Domain | Status |
  |---|---|
  | unknown ISO-4217 code / malformed `amount` | `INVALID_ARGUMENT` |
  | ECB rates not yet loaded | `FAILED_PRECONDITION` |
  | unexpected | `INTERNAL` |

## Build / codegen

Generated proto stubs are **not** committed (see `.gitignore`). Run:

```bash
pip install -e . grpcio-tools
./scripts/gen-proto.sh   # writes app/infrastructure/grpc/generated/
```

The Dockerfile runs the same script during the build stage. Tests trigger
codegen automatically via `tests/conftest.py` if stubs are missing.

## Run

```bash
# from the platform root
make up
# REST
curl -X POST localhost:8000/api/v1/convert \
  -H 'Content-Type: application/json' \
  -d '{"amount":"100.00","from":"USD","to":"EUR"}'
# gRPC (e.g. with grpcurl)
grpcurl -plaintext -d '{"amount":"100.00","from_currency":"USD","to_currency":"EUR"}' \
  localhost:9000 currency.v1.CurrencyConverter/Convert
```

## Architecture

Clean Architecture, four layers:

```
app/
  domain/          # Money, RateSet — pure Python, no framework imports
  application/     # ConvertUseCase + RatesProvider Protocol
  infrastructure/  # ECB-XML fetcher, gRPC server, config, logging
  interfaces/http/ # FastAPI router + DTOs
proto/
  currency/v1/     # gRPC contract (source of truth)
```
