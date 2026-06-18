# Project Brief

## Overview

This is a **workshop / demo monorepo** (`workshop-20260617/demo01`) exploring an
agentic, spec-driven development workflow. It builds a small **customer-messaging /
order platform** split across independent services, each driven by a written spec
in `req/` and implemented with the help of Claude Code skills (go-developer,
nuxt-developer, grill-me, memory-bank, etc.).

The goal is not a single shippable product but a **reference implementation** that
demonstrates clean architecture, test-first development, and tight spec ↔ code
correspondence.

## Components

| Dir            | What it is                                                              | Stack                          |
|----------------|-------------------------------------------------------------------------|--------------------------------|
| `order_api/`   | REST API to **create orders** (idempotent, server-side pricing)         | Go, Gin, MongoDB, Hexagonal    |
| `message_api/` | REST API to **list inbox messages** (paginated, filterable conversations)| Go, Gin, MongoDB, Hexagonal   |
| `web01/`       | **Inbox web UI** — message list, filters, conversation view             | Nuxt 4 / Vue 3, Pinia, Tailwind, Playwright |
| `req/`         | Specs (source of truth): `spec.md`, `spec_inbox_api.md`, `spec_inbox_web.md`, `bug.md`, `demo-async/` | Markdown |

`req/demo-async/demo_spec_line_api.md` describes a **future** service: receiving
LINE / Facebook Messenger messages → Redis dedup → RabbitMQ queue, traced with
OpenTelemetry. Not yet implemented.

## Core Requirements & Goals

1. Each service implements its spec faithfully — request/response shapes, validation
   rules, error codes, and middleware order must match the spec tables exactly.
2. **Hexagonal architecture** (ports & adapters) for the Go services.
3. **Server trust boundary**: prices and totals are computed server-side, never
   trusted from the client (order_api). Listing is tenant-scoped (message_api).
4. Test-first: unit tests with mocked ports, integration tests against real MongoDB
   via testcontainers.
5. The Nuxt UI must match the provided HTML mockup (`req/inbox_page.html`) and expose
   stable `data-testid` hooks for Playwright e2e tests.

## Scope

In scope: the two Go APIs, the Nuxt inbox UI, their specs and tests.
Out of scope (for now): the async LINE/Messenger ingestion service, auth/identity
service (JWT is validated but issued out-of-band), real product catalog management.
