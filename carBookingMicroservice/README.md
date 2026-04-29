# car-booking

Car booking microservice. REST for the frontend, RabbitMQ for talking to
other services in the platform.

- Consumes `user.*` events from `user-auth` to maintain a `Customer` read model.
- Calls `currency-converter` via **RabbitMQ-RPC** (`currency.requests` queue) to
  convert booking totals into the customer's preferred currency.
- Publishes `booking.created` and `booking.cancelled` to its own
  `booking.events` topic exchange.

## Endpoints (REST, JWT-protected)

| Method | Path | Description |
|---|---|---|
| GET | /api/v1/cars | list cars |
| GET | /api/v1/cars/{id} | car details |
| POST | /api/v1/cars | add car (admin) |
| POST | /api/v1/bookings | create booking |
| GET | /api/v1/bookings/me | list own bookings |
| DELETE | /api/v1/bookings/{id} | cancel booking |
| GET | /healthz, /readyz | health (also /actuator/health) |

## Architecture

Clean Architecture — packages enforce the dependency rule:

```
com.uni.carbooking
├── domain              # entities + repository interfaces (no Spring/JPA)
├── application         # use cases + outbound ports
├── infrastructure      # JPA + RabbitMQ + Security + Config
└── interfaces.rest     # controllers + DTOs + exception handler
```

JPA entities live in `infrastructure.persistence.jpa` and are mapped to
domain types in repository adapters — domain stays framework-free.

## Build

Local Maven (requires `mvn` in PATH or use the Docker image):

```bash
mvn -DskipTests package
```

Or via the platform compose:

```bash
cd ..
make up
```
