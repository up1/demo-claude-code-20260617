# producer_api — LINE message producer

REST API that receives LINE (and, by extension, Facebook Messenger) messages,
rejects duplicates, and publishes valid messages to a RabbitMQ work queue. It is
the **producer** half of an async demo; a separate service consumes the queue.

Spec: `../req/demo-async/demo_spec_line_api.md`.

## Architecture (Hexagonal / Ports & Adapters)

```
cmd/main.go                         wiring: config → otel → redis → rabbitmq → routes
internal/
  config/config.go                  env loading (godotenv)
  core/
    domain/line_message.go          LineMessage + Message, Validate(), ErrValidation
    ports/
      message_publisher.go          MessagePublisher (driven)
      dedup_store.go                DedupStore (driven)
      line_service.go               LineService (driving) + SendResult
    service/line_service.go         use-case: dedup → validate → publish
  adapters/
    handler/line_handler.go         POST /api/v1/line/messages (Gin)
    middleware/{logger,request_id}.go
    cache/redis_dedup.go            Redis SETNX impl of DedupStore
    publisher/rabbitmq_publisher.go AMQP 1.0 impl of MessagePublisher
  observability/otel.go             tracer provider (OTLP/gRPC or stdout fallback)
```

The service is **stateless** — no MongoDB. Driven adapters are Redis and RabbitMQ.

## Request flow (`POST /api/v1/line/messages`)

1. Handler reads the **raw body bytes** (needed for both hashing and parsing) and
   unmarshals into `domain.LineMessage`. Malformed JSON → `400 Invalid message format`.
2. Service computes `SHA-256(raw)` → Redis key `line:dedup:<hash>` and calls
   `DedupStore.IsDuplicate` (atomic `SETNX key 1 EX <ttl>`).
   - already seen → `200 {"status":"duplicate"}` (publish is skipped).
3. `LineMessage.Validate()` enforces the schema (`to` required, ≥1 message, each
   message `type` ∈ {text,image}; text needs `text`; image needs valid http(s)
   `originalContentUrl` + `previewImageUrl`). Failure → `ErrValidation` → `400`.
4. Service marshals the spec envelope `{"name":"line-messages","message":<LineMessage>}`
   and publishes (persistent) to the `line-messages` quorum queue.
5. Success → `200 {"status":"success","messageId":"<uuid>"}`.
6. Any infra error (redis/rabbitmq) → `500 {"status":"error","message":"System down"}`.

Every step is traced (OTel span `line.send_message`) and logged via `slog` (JSON).
Gin HTTP spans come from `otelgin` middleware.

## Key decisions

- **Dedup key** = SHA-256 of the raw request body (no message ID in the LINE payload).
- **Auth** = none — public webhook endpoint (matches the spec).
- **RabbitMQ client** = `github.com/rabbitmq/rabbitmq-amqp-go-client` v1.2.0 (AMQP 1.0,
  package `pkg/rabbitmqamqp`). Requires **RabbitMQ ≥ 4.0**. Queue is a durable
  **quorum** queue; the publisher targets it via `QueueAddress` (default exchange =
  work-queue pattern) and checks the publish `Outcome` is `StateAccepted`.
- **Redis** = `github.com/redis/go-redis/v9`.
- **OTel** = v1.44.0 SDK + `otelgin` v0.69.0; OTLP/gRPC when
  `OTEL_EXPORTER_OTLP_ENDPOINT` is set, otherwise stdout exporter (runs without a collector).

## Configuration (env / `.env`, see `.env.example`)

| Var | Required | Default | Purpose |
|-----|----------|---------|---------|
| `PORT` | no | `8080` | HTTP port |
| `SERVICE_NAME` | no | `producer-api` | OTel service / span attribution |
| `REDIS_ADDR` | **yes** | – | Redis host:port |
| `REDIS_PASSWORD` | no | empty | Redis auth |
| `DEDUP_TTL` | no | `24h` | dedup key TTL (Go duration) |
| `RABBITMQ_URI` | **yes** | – | AMQP URI (e.g. `amqp://guest:guest@host:5672/`) |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | no | empty (stdout) | OTLP/gRPC collector host:port |

## Run & verify

```bash
cp .env.example .env
docker compose up -d --build        # api + redis:7 + rabbitmq:4-management

curl -i -X POST localhost:8080/api/v1/line/messages -H 'Content-Type: application/json' \
  -d '{"to":"U1","messages":[{"type":"text","text":"hi"}]}'   # 200 success + messageId
# repeat same body            → 200 {"status":"duplicate"}
# {"to":"","messages":[]}     → 400 Invalid message format
# malformed JSON              → 400 Invalid message format

# Inspect the queue: management UI http://localhost:15672 (guest/guest), queue "line-messages"
docker compose down -v
```

Tests: `go test ./...` (handler httptest + mock service; service unit tests with mock
publisher + dedup store covering happy path, duplicate short-circuit, validation,
publish error, dedup error). `go build ./...` and `go vet ./...` are clean.

## Notes for future work

- Endpoint is public per spec; if exposed beyond the demo, add LINE signature
  verification (`X-Line-Signature`) or the repo's JWT middleware.
- Spec doesn't define the duplicate response body; we return `200 {"status":"duplicate"}`.
- Facebook Messenger ingestion is mentioned in the spec overview but not yet a
  separate endpoint — the domain/validation would extend here.
