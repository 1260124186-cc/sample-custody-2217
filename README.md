# Sample Custody Service

Sample Custody Service is a local HTTP backend for recording the custody chain of
scientific specimens. It is designed for an intake station and a review client
that need a small, auditable workflow without a database dependency.

## Roles and capabilities

- Intake technicians register specimens with origin and quantity details.
- Custody coordinators record handoffs between people and storage locations.
- Review scientists create inspection batches, complete them, and read manifests.

## Structure

- `cmd/server`: process entry point and graceful shutdown.
- `internal/httpapi`: HTTP routing, request decoding, and response mapping.
- `internal/service`: business workflows and cross-entity rules.
- `internal/store`: synchronized in-memory repositories.
- `internal/model`: domain entities, statuses, and typed errors.
- `internal/validation`: request validation and normalization.
- `internal/export`: review manifest rendering.
- `internal/audit`: append-only event recording.

## Run

```bash
go run ./cmd/server
```

The server listens on `127.0.0.1:8090`. Set `SAMPLE_CUSTODY_ADDR` to use a
different address.

## Build and test

```bash
go build ./...
go test ./...
```

## HTTP entry points

- `GET /healthz`
- `POST /samples`
- `GET /samples`
- `GET /samples/{id}`
- `POST /samples/{id}/transfers`
- `POST /batches`
- `GET /batches/{id}`
- `POST /batches/{id}/complete`
- `GET /manifests/{id}`
- `GET /summary`
- `GET /search?q=...`
- `GET /exports/samples`

All request and response bodies use JSON except the manifest endpoint, which
returns plain text.
