# Tech Context

## Go services (order_api, message_api)

- **Go 1.26.1** (module name: `api` in both).
- **Gin** (`github.com/gin-gonic/gin v1.12.0`) — HTTP framework.
- **MongoDB driver** `go.mongodb.org/mongo-driver v1.17.9` (v2 also pulled transitively).
- **JWT** `github.com/golang-jwt/jwt/v5`.
- **UUID** `github.com/google/uuid` (request IDs).
- **godotenv** for `.env` loading.
- **Testing**: `stretchr/testify` (assert/mock) + **testcontainers-go** with the
  `modules/mongodb` package for integration tests against a real Mongo container.
- `log/slog` for structured logging.

### Config (env)
```
MONGO_URI=mongodb://localhost:27017
MONGO_DB=orders_db          # message_api uses its own DB name
JWT_SECRET=change_me
PORT=8080
```
Loaded into a `Config` struct at startup; app fails fast if a required var is missing.
See each service's `.env.example`.

### MongoDB indexes
- order_api: `orders.idempotency_key` (unique), `orders.customer_id` (normal).
- message_api: `inbox.updated_at` (desc), compound `(channel, status, updated_at)`,
  optional text index on `(sender_name, preview)` — `$regex` fallback for search `q`.

### Build / run
- Multi-stage Dockerfile (golang:1.23-alpine builder → distroless static runtime).
  NOTE: Dockerfile pins golang:1.23 while go.mod says 1.26.1 — potential mismatch to watch.
- `docker-compose.yml` runs the api + a `mongo:8` container.
- `message_api/cmd/seed/main.go` seeds inbox fixture data.

## Web (web01)

- **Nuxt ^4.4.8 / Vue ^3.5.35**, TypeScript.
- **Pinia** (`@pinia/nuxt`) for state.
- **@nuxt/ui ^4.8.2** + **Tailwind CSS ^4** for UI.
- **axios** for API calls.
- **Playwright** (`@playwright/test`) for e2e — `npm test` runs `playwright test`;
  config in `playwright.config.ts`, results in `playwright-report/`, `test-results/`,
  `junit.xml`.
- Dev server: `npm run dev` → http://localhost:3000. Build: `npm run build`.
- `nuxt.config.ts` exposes `runtimeConfig.public.apiBaseUrl` (empty by default — set to
  the Go API origin, e.g. http://localhost:8080, when wiring front to back).

## Dev environment notes
- Platform: macOS (darwin), zsh.
- **RTK (Rust Token Killer)** hook rewrites common CLI commands (`git`, `find`, …) to
  `rtk <cmd>`. `rtk find` does NOT support compound predicates (`-not`, `-exec`) — use
  plain `find` or per-dir `ls -R` instead. `rtk proxy <cmd>` runs raw (may prompt).
- **Serena** MCP server available (symbol tools); call `initial_instructions` before
  heavy coding via Serena.

## Skills used in this repo
`go-developer` (build/test the Go APIs, hexagonal + MongoDB), `nuxt-developer`
(Nuxt + Playwright), `memory-bank` (this), plus `grill-me`, `adr-skill`, `tdd`, etc.
