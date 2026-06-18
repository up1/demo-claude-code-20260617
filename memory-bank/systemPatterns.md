# System Patterns

## Architecture: Hexagonal (Ports & Adapters) for Go services

Both `order_api` and `message_api` share the same layout under `internal/`:

```
internal/
  core/
    domain/      # entities + enums (order.go / inbox.go) — no framework deps
    ports/       # interfaces: *_repository.go (driven), *_service.go (driving) + I/O value objects
    service/     # business logic implementing the driving port; depends only on driven ports
  adapters/
    handler/     # Gin HTTP handlers — parse request, call service, shape JSON response
    repository/  # MongoDB adapters implementing driven ports
  config/        # env loading into Config struct (fail-fast on missing required vars)
  middleware/    # jwt.go, logger.go, request_id.go
cmd/main.go      # wiring: repo → service → handler → router
```

### Dependency rule
`handler → service (driving port) → repository (driven port)`. The service depends on
**port interfaces**, never on concrete Mongo types. This is what makes service unit
tests possible with mocks.

## Key patterns

- **Response envelope**: success → `{ "success": true, "data": ..., ["pagination": ...] }`;
  error → `{ "success": false, "error": { "code": "...", "message": "..." } }`.
  Error codes: `VALIDATION_ERROR` (400), `UNAUTHORIZED` (401), `PRODUCT_NOT_FOUND` (404),
  `INTERNAL_ERROR` (500).

- **Validation lives in the service**, not the handler. The service normalizes defaults
  (`page=1`, `page_size=20`), clamps (`page_size` max 100), computes `Offset`, and
  derives `TotalPages = ceil(total / pageSize)`.

- **Idempotency (order_api)**: `idempotency_key` carries a unique Mongo index.
  `FindByIdempotencyKey` first; if found return existing order with 200, else save +
  201. Price/total always computed server-side from the `products` collection.

- **Tenant scoping (message_api)**: JWT claims inject tenant/agent into context; the
  service scopes every query to the caller's tenant.

- **Repository returns page + total in one call**: `List(ctx, filter) ([]T, int64, error)`
  — items plus `CountDocuments` total, so the service can compute pagination.

### Middleware order differs between the two services (per spec — keep exact)
- **order_api**: `JWT Auth → Request ID → Logger → Handler`
- **message_api**: `Request ID → Logger → JWT Auth → Handler`

All three middlewares: JWT validates `Authorization: Bearer <jwt>` and injects claims;
Request ID generates a uuid, sets `X-Request-ID`, injects into context; Logger emits
structured JSON (method, path, status, latency, request_id) via `slog`.

## Testing pattern (per layer)
- **Service** — unit; mock the repository port(s). Assert validation, defaults, clamping,
  pagination math, idempotency branch.
- **Handler** — unit; `httptest` + mock service. Assert query/body parsing, status codes,
  JSON shape.
- **Repository** — integration; **real MongoDB via testcontainers-go** (mongodb module).
  Seed fixtures, assert filter/sort/paging.

## Web (web01) patterns
- **Pinia store** (`app/stores/inbox.ts`) is the single source of truth for inbox state
  and the only place that talks to the API (axios instance with `Authorization` header).
- Changing filter/search/status resets `page = 1` then refetches. `goToPage` guards bounds.
- 401 → sets `unauthorized` + friendly message; other errors → generic message; both clear
  the list. `isEmpty` getter drives the "No messages found" state.
- Components mirror `req/inbox_page.html`; interactive elements carry `data-testid`
  (`filter-channel`, `filter-status`, `search-input`, `pagination-controls`,
  `message-item-{id}`, `thread-link-{id}`).
