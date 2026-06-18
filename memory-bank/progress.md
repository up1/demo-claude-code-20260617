# Progress

_Status as of 2026-06-18 (initial assessment from file structure + git, not a full test run)._

## What works / is built

### message_api — `GET /api/v1/inbox/messages`
- ✅ domain (`inbox.go`), ports (`inbox_repository.go`, `inbox_service.go`)
- ✅ service + unit tests (`inbox_service.go` / `_test.go`)
- ✅ Mongo repository + integration test (testcontainers) (`inbox_mongo.go` / `_test.go`)
- ✅ handler + unit tests (`inbox_handler.go` / `_test.go`)
- ✅ middleware (jwt, logger, request_id), config, `cmd/seed` seeder
- ✅ Dockerfile + docker-compose (mongo:8)

### order_api — `POST /api/v1/orders`
- ✅ domain (`order.go`), ports (order repo, product repo, order service)
- ✅ service + unit tests (`order_service.go` / `_test.go`)
- ✅ Mongo repos (`order_mongo.go`, `product_mongo.go`)
- ✅ handler + unit tests (`order_handler.go` / `_test.go`)
- ✅ middleware, config, Dockerfile, docker-compose
- ⚠️ No integration test for order/product repos seen (message_api has one).

### web01 — Inbox UI
- ✅ Pages: `index.vue` (inbox), `conversations/[thread_id].vue`
- ✅ Components: InboxList, AppSidebar, AppTopbar, ChatPanel, ContactPane
- ✅ Pinia `inbox` store wired to `GET /api/v1/inbox/messages` (filters, search, pagination, 401 handling)
- ✅ Playwright scaffolding (config, report dirs, junit.xml)

## What's left to build

- ❓ Conversation detail endpoint `GET /api/v1/conversations/{thread_id}` — referenced by
  web spec + `[thread_id].vue`, but no matching Go handler/route confirmed. Verify/implement.
- ❓ Order repo integration tests (parity with message_api).
- ⬜ Async ingestion service (`req/demo-async/demo_spec_line_api.md`): LINE/Messenger →
  Redis dedup → RabbitMQ `line-messages` → OpenTelemetry tracing. Not started.
- ⬜ Wire `web01` `runtimeConfig.public.apiBaseUrl` to a running API and verify e2e green.

## Known issues

- `req/bug.md`: claims `spec_inbox_api.md` is missing the `POST /api/v1/orders`
  endpoint. That endpoint is defined in `spec.md`, not the inbox spec — needs reconciling.
- Dockerfile builder image `golang:1.23-alpine` vs go.mod `go 1.26.1` — version mismatch
  could break container builds.

## Verification still needed

The above ✅ marks reflect presence of files, not a confirmed passing test run. Before
relying on any component: run `go test ./...` in each API dir and `npm test` in web01.

## Evolution of project decisions

- Started as a Go demo (`afdea53 Demo with go`), then HTML mockup → web spec → Nuxt app.
- Specs and the HTML mockup were consolidated from repo root into `req/`.
- Integration testing migrated to **testcontainers** (`cf742d0`) for real-Mongo repo tests.
- MongoDB image/driver version bumped (`3b6bf5f`).
