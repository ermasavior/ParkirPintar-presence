# presence-service

[![Golang CI/CD](https://github.com/ermasavior/parkirpintar-presence/actions/workflows/cicd.yml/badge.svg)](https://github.com/ermasavior/parkirpintar-presence/actions/workflows/cicd.yml)

Manages driver check-in and check-out. Bridges the reservation lifecycle into a billing session.

## Responsibilities

- `CheckIn` — validates reservation is `CONFIRMED` and not expired, creates a session (`ACTIVE`), marks reservation `CHECKED_IN`
- `CheckOut` — completes the session, releases the spot to `AVAILABLE`, marks reservation `COMPLETED`, calls Billing Service to calculate the invoice and get a parking-fee QRIS code
- `GetSession` — returns current session state

## gRPC API

```
service PresenceService {
  rpc CheckIn    (CheckInRequest)    returns (CheckInResponse);
  rpc CheckOut   (CheckOutRequest)   returns (CheckOutResponse);
  rpc GetSession (GetSessionRequest) returns (GetSessionResponse);
}
```

Proto: [`proto/presence/v1/presence.proto`](proto/presence/v1/presence.proto)

## Dependencies

| Dependency | Purpose |
|---|---|
| PostgreSQL | Sessions, reservations, spots |
| Billing Service (gRPC) | Calculate invoice + create parking-fee payment |

## Configuration

```bash
cp .env.example .env
```

Key variables: `POSTGRES_DSN`, `BILLING_SERVICE_URL`

## Development

```bash
make run              # run locally
make build            # compile binary → bin/presence
make test             # all tests
make test-unit        # unit tests only
make unit-test-coverage
make proto            # regenerate gRPC code from .proto
make mock             # regenerate mocks
```
