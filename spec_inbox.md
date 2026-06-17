# Web :: Inbox Message List
* page= /
* Component= InboxList

## Implementation of UI
1. Use HTML template from file `inbox_page.html`
2. Crete components from html template
3. Compare nuxt's components with the html template and make sure they are the same

## Workflow of web page of inbox message list
1. User visits the dashboard page (`/`).
2. The frontend sends a `GET /api/v1/inbox/messages` request to load the list of inbox messages.
3. Receives a paginated list of conversation threads, each showing the latest message snippet, sender info, channel, status, and unread indicator.
   - Render result list in location => testid=message-item-{id}
4. User can filter by channel/status, search by sender/message, and paginate through results.
   - Elements:
     - Channel filter => testid=filter-channel
     - Status filter => testid=filter-status
     - Search input => testid=search-input
     - Pagination controls => testid=pagination-controls
5. Clicking a thread navigates to the conversation view (not covered in this spec).

## Web test cases (Playwright)
| Test Case                 | Steps                                                                                     | Expected Result                                  |
|---------------------------|-------------------------------------------------------------------------------------------|--------------------------------------------------|
| Load inbox messages       | Visit `/`                                                                                 | 200 OK, list of threads rendered with correct info |
| Filter by channel         | Select `line` in channel filter                                                           |
| 200 OK, only `line` threads rendered                                                   |
| Filter by status          | Select `pending` in status filter                                                         | 200 OK, only `pending` threads rendered                                                  |
| Combined filter + search  | Select `line` + enter `tracking` in search input                                                        | 200 OK, matching subset rendered                                                  |
| Pagination                | Click page 2 in pagination controls                                                        | 200 OK, items 11–20 rendered, correct pagination state                                                  |
| Empty result              | Apply filter/search matching nothing                                                        |
| 200 OK, "No messages found" displayed                                                  |
| Invalid JWT               | Remove/alter `Authorization` header in request                                                 | 401 Unauthorized                                                  |

# REST API — List Inbox Messages (Go)

## Domain Model

### InboxMessage

```json
{
  "_id":         "ObjectID",
  "customer_id": "string",
  "sender_name": "string",
  "avatar_url":  "string",
  "channel":     "facebook | line | instagram",
  "preview":     "string (last message snippet)",
  "status":      "pending | replied",
  "unread":      "bool",
  "created_at":  "time.Time",
  "updated_at":  "time.Time (last activity, used for sort)"
}
```

### Business Rules

- An inbox message represents the **latest state of a conversation thread**, one document per thread.
- `channel` is a closed enum: `facebook`, `line`, `instagram`. Filtering by an unknown channel is a validation error.
- `status` is a closed enum: `pending` (awaiting agent reply) or `replied`.
- Results are **always sorted by `updated_at` descending** (most recent activity first).
- Listing is **read-only** and scoped to the authenticated agent's tenant — never returns other tenants' threads.
- Pagination is mandatory server-side; an unbounded list is never returned.

---

## API Endpoint

### `GET /api/v1/inbox/messages`

**Headers**

| Header          | Required | Description     |
|-----------------|----------|-----------------|
| `Authorization` | Yes      | `Bearer <jwt>`  |

**Query Parameters**

| Param      | Type   | Required | Default | Description                                                |
|------------|--------|----------|---------|------------------------------------------------------------|
| `channel`  | string | No       | (all)   | Filter by channel: `facebook` \| `line` \| `instagram`     |
| `status`   | string | No       | (all)   | Filter by status: `pending` \| `replied`                   |
| `q`        | string | No       | —       | Case-insensitive search over `sender_name` and `preview`   |
| `from`     | string | No       | —       | RFC3339 timestamp; include threads with `updated_at >= from`|
| `to`       | string | No       | —       | RFC3339 timestamp; include threads with `updated_at <= to` |
| `page`     | int    | No       | `1`     | 1-based page number (min 1)                                |
| `page_size`| int    | No       | `20`    | Items per page (min 1, max 100)                            |

**Validation Rules**

| Field       | Rule                                                  |
|-------------|-------------------------------------------------------|
| `channel`   | optional; one of `facebook`, `line`, `instagram`      |
| `status`    | optional; one of `pending`, `replied`                 |
| `from`/`to` | optional; valid RFC3339; if both set, `from` <= `to`  |
| `page`      | optional; integer >= 1                                |
| `page_size` | optional; integer in `[1, 100]`                       |

**Success Response — `200 OK`**

```json
{
  "success": true,
  "data": [
    {
      "id":          "64f1a2b3c4d5e6f7a8b9c0d1",
      "customer_id": "cust_123",
      "sender_name": "Marcus Watanabe",
      "avatar_url":  "https://.../marcus.jpg",
      "channel":     "line",
      "preview":     "Can you confirm the tracking number for order #8812?",
      "status":      "pending",
      "unread":      true,
      "created_at":  "2026-06-17T10:30:00Z",
      "updated_at":  "2026-06-17T10:30:00Z"
    },
    {
      "id":          "64f1a2b3c4d5e6f7a8b9c0d2",
      "customer_id": "cust_456",
      "sender_name": "Sarah Jenkins",
      "avatar_url":  "https://.../sarah.jpg",
      "channel":     "facebook",
      "preview":     "Thank you for the quick resolution!",
      "status":      "replied",
      "unread":      false,
      "created_at":  "2026-06-16T14:05:00Z",
      "updated_at":  "2026-06-16T14:05:00Z"
    }
  ],
  "pagination": {
    "page":        1,
    "page_size":   20,
    "total_items": 128,
    "total_pages": 7
  }
}
```

**Empty Result — `200 OK`**

```json
{
  "success": true,
  "data": [],
  "pagination": { "page": 1, "page_size": 20, "total_items": 0, "total_pages": 0 }
}
```

**Error Response**

```json
{
  "success": false,
  "error": {
    "code":    "VALIDATION_ERROR",
    "message": "channel must be one of facebook, line, instagram"
  }
}
```

**Error Codes**

| HTTP | Code               | Cause                                          |
|------|--------------------|------------------------------------------------|
| 400  | `VALIDATION_ERROR` | Invalid query param (channel/status/date/page) |
| 401  | `UNAUTHORIZED`     | Missing or invalid JWT                          |
| 500  | `INTERNAL_ERROR`   | Unexpected server error                         |

---

## Ports (Interfaces)

```go
// InboxRepository — driven port (MongoDB adapter implements this)
type InboxRepository interface {
    List(ctx context.Context, filter InboxFilter) ([]domain.InboxMessage, int64, error)
}

// InboxService — driving port (handler calls this)
type InboxService interface {
    ListMessages(ctx context.Context, req ListInboxInput) (*ListInboxResult, error)
}
```

```go
// Input/output value objects (ports package)
type InboxFilter struct {
    Channel  string     // "" = all
    Status   string     // "" = all
    Search   string     // "" = none
    From     *time.Time // nil = unbounded
    To       *time.Time // nil = unbounded
    Offset   int        // (page-1) * pageSize
    Limit    int        // pageSize
}

type ListInboxInput struct {
    Channel  string
    Status   string
    Search   string
    From     *time.Time
    To       *time.Time
    Page     int
    PageSize int
}

type ListInboxResult struct {
    Items      []domain.InboxMessage
    Page       int
    PageSize   int
    TotalItems int64
    TotalPages int
}
```

- The **service** normalizes/validates input (defaults `page=1`, `page_size=20`, clamps to max 100), computes `Offset`, and maps the repo's `(items, total)` into `ListInboxResult` (deriving `TotalPages = ceil(total / pageSize)`).
- The **repository** builds the Mongo query (filters + `$regex` for search) and returns the page plus a `CountDocuments` total in a single call.

---

## Middleware Chain

```
Request ID → Logger → JWT Auth → Handler
```

- **JWT Auth**: validates `Authorization` header, injects claims (tenant/agent) into context; the service scopes the query to the caller's tenant.
- **Request ID**: generates `uuid` per request, sets `X-Request-ID` response header, injects into context for logging.
- **Logger**: logs method, path, status, latency, request_id as structured JSON via `slog`.

---

## MongoDB Indexes

| Collection | Field(s)                          | Type     | Purpose                                  |
|------------|-----------------------------------|----------|------------------------------------------|
| `inbox`    | `updated_at`                      | Normal (desc) | Default sort                        |
| `inbox`    | `channel`, `status`, `updated_at` | Compound | Filtered list + sort                     |
| `inbox`    | `sender_name`, `preview`          | Text (optional) | Search via `q` (or `$regex` fallback) |

---

## Testing Strategy

| Layer      | Type        | Approach                                              |
|------------|-------------|-------------------------------------------------------|
| Service    | Unit        | Mock `InboxRepository`; assert defaults, clamping, pagination math, filter pass-through |
| Repository | Integration | Real MongoDB via Docker Compose; seed fixtures, assert filter/sort/paging |
| Handler    | Unit        | `httptest` + mock `InboxService`; assert query parsing, status codes, JSON shape |

## Test cases in table format

| Test Case                 | Input                                              | Expected Output                                  |
|---------------------------|----------------------------------------------------|--------------------------------------------------|
| List default              | No query params                                    | 200 OK, page 1, page_size 20, sorted by recency  |
| Filter by channel         | `?channel=line`                                    | 200 OK, only `line` threads                      |
| Filter by status          | `?status=pending`                                  | 200 OK, only `pending` threads                   |
| Combined filter + search  | `?channel=line&q=tracking`                         | 200 OK, matching subset                          |
| Pagination                | `?page=2&page_size=10`                             | 200 OK, items 11–20, correct `total_pages`       |
| Empty result              | Filter matching nothing                            | 200 OK, `data: []`, `total_items: 0`             |
| Invalid channel           | `?channel=tiktok`                                  | 400 `VALIDATION_ERROR`                            |
| Invalid page_size         | `?page_size=500`                                   | 400 `VALIDATION_ERROR` (or clamped to 100)       |
| Invalid date range        | `?from=2026-06-18T00:00:00Z&to=2026-06-17T00:00:00Z` | 400 `VALIDATION_ERROR`                          |
| Missing/invalid JWT       | No `Authorization` header                          | 401 `UNAUTHORIZED`                               |

---

## Wiring (cmd/main.go)

```go
inboxRepo := repository.NewInboxMongoRepo(db)
inboxSvc := service.NewInboxService(inboxRepo)
inboxHandler := handler.NewInboxHandler(inboxSvc)

api := r.Group("/api/v1")
api.Use(middleware.JWTAuth(cfg.JWTSecret))
{
    api.GET("/inbox/messages", inboxHandler.ListMessages)
}
```

---

## File Layout (new files only)

```
internal/
  core/
    domain/inbox.go                       # InboxMessage, channel/status enums
    ports/inbox_repository.go             # InboxRepository, InboxFilter
    ports/inbox_service.go                # InboxService, ListInboxInput/Result
    service/inbox_service.go              # validation, defaults, pagination math
    service/inbox_service_test.go         # unit tests (mock repo)
  adapters/
    handler/inbox_handler.go              # query parsing, response shaping
    handler/inbox_handler_test.go         # unit tests (mock service)
    repository/inbox_mongo.go             # Mongo query + count
```
