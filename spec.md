# REST API — Create Order (Go)

## Stack

| Concern        | Choice                                      |
|----------------|---------------------------------------------|
| Language       | Go (latest stable)                          |
| HTTP Framework | Gin (latest)                                |
| Database       | MongoDB — official `go.mongodb.org/mongo-driver` |
| Architecture   | Hexagonal (Ports & Adapters)                |
| Auth           | JWT — `Authorization: Bearer <token>`       |
| Config         | Environment variables + `godotenv` (local)  |
| Logging        | `slog` (stdlib, structured JSON)            |
| Deployment     | Docker Compose (local) + multi-stage Dockerfile (prod) |

---

## Project Structure

```
.
├── cmd/
│   └── main.go                  # wiring: config, DB, routes
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   │   └── order.go         # Order, OrderItem structs + business rules
│   │   ├── ports/
│   │   │   ├── order_repository.go   # OrderRepository interface (driven)
│   │   │   ├── product_repository.go # ProductRepository interface (driven)
│   │   │   └── order_service.go      # OrderService interface (driving)
│   │   └── service/
│   │       └── order_service.go # use-case logic, implements OrderService port
│   └── adapters/
│       ├── handler/
│       │   └── order_handler.go # Gin HTTP adapter (driving)
│       └── repository/
│           ├── order_mongo.go   # MongoDB impl of OrderRepository
│           └── product_mongo.go # MongoDB impl of ProductRepository
├── docker-compose.yml
├── Dockerfile
└── .env.example
```

---

## Domain Model

### Order

```json
{
  "_id":          "ObjectID",
  "customer_id":  "string (required)",
  "items": [
    {
      "product_id":  "string",
      "quantity":    "int (> 0)",
      "unit_price":  "float64 (fetched server-side)"
    }
  ],
  "total_amount":    "float64 (calculated server-side)",
  "status":          "pending | confirmed | cancelled",
  "idempotency_key": "string (unique index)",
  "created_at":      "time.Time",
  "updated_at":      "time.Time"
}
```

### Business Rules

- `unit_price` is **always fetched** from the `products` collection — never trusted from the client.
- `total_amount` is **always calculated** server-side: `sum(quantity × unit_price)`.
- `idempotency_key` must be unique (MongoDB unique index). Duplicate submission returns the existing order (HTTP 200), not an error.
- Initial status on creation: `pending`.

---

## API Endpoint

### `POST /api/v1/orders`

**Headers**

| Header              | Required | Description                        |
|---------------------|----------|------------------------------------|
| `Authorization`     | Yes      | `Bearer <jwt>`                     |
| `Idempotency-Key`   | Yes      | Client-generated unique string     |
| `Content-Type`      | Yes      | `application/json`                 |

**Request Body**

```json
{
  "customer_id": "cust_123",
  "items": [
    { "product_id": "prod_abc", "quantity": 2 },
    { "product_id": "prod_xyz", "quantity": 1 }
  ]
}
```

**Validation Rules**

| Field         | Rule                          |
|---------------|-------------------------------|
| `customer_id` | required                      |
| `items`       | required, min 1 element       |
| `product_id`  | required per item             |
| `quantity`    | required, integer > 0         |

**Success Response — `201 Created`**

```json
{
  "success": true,
  "data": {
    "id":           "64f1a2b3c4d5e6f7a8b9c0d1",
    "customer_id":  "cust_123",
    "items": [
      { "product_id": "prod_abc", "quantity": 2, "unit_price": 25.00 },
      { "product_id": "prod_xyz", "quantity": 1, "unit_price": 10.00 }
    ],
    "total_amount": 60.00,
    "status":       "pending",
    "created_at":   "2026-06-17T10:00:00Z"
  }
}
```

**Idempotent Repeat — `200 OK`** — same body as above (existing order).

**Error Response**

```json
{
  "success": false,
  "error": {
    "code":    "VALIDATION_ERROR",
    "message": "items must not be empty"
  }
}
```

**Error Codes**

| HTTP | Code                  | Cause                              |
|------|-----------------------|------------------------------------|
| 400  | `VALIDATION_ERROR`    | Missing/invalid request fields     |
| 401  | `UNAUTHORIZED`        | Missing or invalid JWT             |
| 404  | `PRODUCT_NOT_FOUND`   | product_id does not exist          |
| 500  | `INTERNAL_ERROR`      | Unexpected server error            |

---

## Ports (Interfaces)

```go
// OrderRepository — driven port (MongoDB adapter implements this)
type OrderRepository interface {
    Save(ctx context.Context, order *domain.Order) error
    FindByIdempotencyKey(ctx context.Context, key string) (*domain.Order, error)
}

// ProductRepository — driven port
type ProductRepository interface {
    FindByID(ctx context.Context, productID string) (*domain.Product, error)
}

// OrderService — driving port (handler calls this)
type OrderService interface {
    CreateOrder(ctx context.Context, req CreateOrderInput) (*domain.Order, error)
}
```

---

## Middleware Chain

```
JWT Auth → Request ID → Logger → Handler
```

- **JWT Auth**: validates `Authorization` header, injects claims into context.
- **Request ID**: generates `uuid` per request, sets `X-Request-ID` response header, injects into context for logging.
- **Logger**: logs method, path, status, latency, request_id as structured JSON via `slog`.

---

## Configuration

`.env.example`

```env
MONGO_URI=mongodb://localhost:27017
MONGO_DB=orders_db
JWT_SECRET=change_me
PORT=8080
```

Loaded at startup into a `Config` struct. App fails fast if any required variable is missing.

---

## MongoDB Indexes

| Collection | Field             | Type   |
|------------|-------------------|--------|
| `orders`   | `idempotency_key` | Unique |
| `orders`   | `customer_id`     | Normal |
| `products` | `_id`             | (default) |

---

## Testing Strategy

| Layer       | Type        | Approach                                              |
|-------------|-------------|-------------------------------------------------------|
| Service     | Unit        | Mock `OrderRepository` and `ProductRepository` via port interfaces |
| Repository  | Integration | Real MongoDB via Docker Compose in test environment   |
| Handler     | Unit        | `httptest` + mock `OrderService`                      |

---

## Docker Compose

```yaml
services:
  api:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    depends_on:
      - mongodb

  mongodb:
    image: mongo:7
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db

volumes:
  mongo_data:
```

---

## Dockerfile (multi-stage)

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/main.go

# Runtime stage
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/server /server
ENTRYPOINT ["/server"]
```
