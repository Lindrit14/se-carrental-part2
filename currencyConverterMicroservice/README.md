# currency-converter

Stateless currency conversion microservice. Pulls daily rates from the
European Central Bank XML feed, exposes them via REST and RabbitMQ-RPC.

## Endpoints

REST (public):

- `GET  /healthz` — liveness
- `GET  /api/v1/rates` — current rate set (EUR-based)
- `POST /api/v1/convert` — convert an amount

RabbitMQ-RPC:

- Queue `currency.requests` — same payload as `POST /api/v1/convert`,
  reply published to caller's `reply_to` with the matching `correlation_id`.

## Run

```bash
# from the platform root
make up
# then:
curl -X POST localhost:8000/api/v1/convert \
  -H 'Content-Type: application/json' \
  -d '{"amount":"100.00","from":"USD","to":"EUR"}'
```

## Architecture

Clean Architecture, four layers:

```
app/
  domain/         # Money, RateSet — pure Python, no framework imports
  application/    # ConvertUseCase + RatesProvider Protocol
  infrastructure/ # ECB-XML fetcher, RabbitMQ, config, logging
  interfaces/http/# FastAPI router + DTOs
```
