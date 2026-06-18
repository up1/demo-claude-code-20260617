# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/main.go

# API runtime
FROM gcr.io/distroless/static:nonroot AS api
COPY --from=builder /server /server
ENTRYPOINT ["/server"]
