# Active Context

_Last updated: 2026-06-18 (initial Memory Bank creation)._

## Current work focus

Memory Bank just initialized. The repo is a spec-driven demo with three components in
varying states of completeness:

- **message_api** — most built out: domain, ports, service (+unit tests), Mongo repo
  (+integration test via testcontainers), handler (+unit tests), middleware, config,
  and a `cmd/seed` seeder. Implements `GET /api/v1/inbox/messages` per `spec_inbox_api.md`.
- **order_api** — full hexagonal skeleton present: domain, ports (order + product repos,
  order service), service (+test), order/product Mongo repos, handler (+test), middleware,
  config. Implements `POST /api/v1/orders` per `spec.md`.
- **web01** — Nuxt inbox UI: pages (`index.vue`, `conversations/[thread_id].vue`),
  components (InboxList, AppSidebar, AppTopbar, ChatPanel, ContactPane), Pinia `inbox`
  store wired to the API, Playwright tests scaffolded.

## Recent changes (from git history)

- `3b6bf5f` Update MongoDB version
- `cf742d0` Add integration testing with testcontainers
- `9e9e7ea` Add more cases
- `4a7757a` Update spec
- `27a30d7` Call api (web → API wiring)
- earlier: nuxt-developer/go-developer skills created, web spec + HTML mockup added.

## Uncommitted working-tree state (snapshot)

- Modified: `.claude/skills/go-developer/SKILL.md`
- Deleted (staged): `inbox_page.html`, `spec.md`, `spec_inbox.md` at repo root — these
  moved into `req/` (now `req/inbox_page.html`, `req/spec.md`, `req/spec_inbox_*.md`).
- Untracked: `.serena/`, `.vscode/`, `message_api/`, `req/`, `.claude/skills/memory-bank/`.

## Next steps (candidates)

1. Address `req/bug.md`: `spec_inbox_api.md` is noted as missing the `POST /api/v1/orders`
   endpoint definition — that endpoint actually lives in `spec.md`. Reconcile which spec
   owns which endpoint, or copy the order endpoint definition into the inbox spec if the
   bug intends both APIs documented together.
2. Wire `web01` `apiBaseUrl` to the running message_api and run Playwright e2e green.
3. Consider implementing the planned async LINE/Messenger → RabbitMQ service
   (`req/demo-async/demo_spec_line_api.md`).

## Active decisions & considerations

- Specs in `req/` are the **source of truth**; code must match spec tables exactly
  (status codes, error codes, middleware order, response envelope).
- Middleware order intentionally differs between the two APIs — don't "normalize" it.
- Dockerfile pins golang:1.23 but go.mod is 1.26.1 — verify before container builds.
