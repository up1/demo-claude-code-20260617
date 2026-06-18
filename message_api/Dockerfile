# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/main.go
RUN CGO_ENABLED=0 go build -o /seed ./cmd/seed

# API runtime (default target)
FROM gcr.io/distroless/static:nonroot AS api
COPY --from=builder /server /server
ENTRYPOINT ["/server"]

# Seeder runtime (one-off mock-data loader)
FROM gcr.io/distroless/static:nonroot AS seeder
COPY --from=builder /seed /seed
ENTRYPOINT ["/seed"]
