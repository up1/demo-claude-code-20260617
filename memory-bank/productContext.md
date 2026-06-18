# Product Context

## Why this project exists

The repo is a teaching/demo vehicle for a customer-support + commerce platform. It
models the slice of a SaaS helpdesk where agents triage incoming customer
conversations across channels and where orders are placed reliably.

## Problems it solves

### Inbox (message_api + web01)
Support agents need a single inbox that aggregates conversation threads from multiple
channels (Facebook, LINE, Instagram). They must be able to:
- See the latest snippet, sender, channel, status, and unread state per thread.
- Filter by channel and status, search by sender/preview text, and page through results.
- Open a thread to view the full conversation.

One document per thread represents the **latest state** of that conversation, always
sorted by most recent activity (`updated_at` desc). Listing is read-only and scoped
to the authenticated agent's tenant — agents never see other tenants' threads.

### Orders (order_api)
Customers place orders. The system must guarantee:
- **Correct pricing** — `unit_price` is always fetched from the `products` collection;
  `total_amount` is always recomputed server-side. The client cannot influence price.
- **Idempotency** — a client-supplied `Idempotency-Key` (unique index) makes retries
  safe: a duplicate submission returns the existing order (HTTP 200), not a new one or
  an error.
- Orders start in `pending` status.

### Async ingestion (planned — demo-async)
Receive messages from LINE / Facebook Messenger, deduplicate via Redis, and publish to
a RabbitMQ `line-messages` work queue, with full OpenTelemetry tracing. Decouples
message receipt from downstream processing.

## How it should work (user experience)

- **Agent** opens the web app at `/`, sees the inbox list, filters/searches/paginates,
  clicks a thread → navigates to `/conversations/{thread_id}` showing the conversation.
- Empty results show "No messages found"; expired/invalid sessions surface a clear
  re-auth message (401 handling).
- All API responses use a consistent envelope: `{ success, data, ... }` on success and
  `{ success: false, error: { code, message } }` on failure.

## UX goals
- Match the provided HTML mockup pixel-for-pixel (`req/inbox_page.html`).
- Predictable, stable `data-testid` hooks for every interactive element.
- Fast, paginated lists — never an unbounded fetch.
