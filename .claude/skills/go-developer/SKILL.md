---
name: go-developer
description: Develop and Test REST API with Go, Gin, MongoDB, and Hexagonal Architecture
---

## Technology Stack
| Concern        | Choice                                      |
|----------------|---------------------------------------------|
| Language       | Go 1.25                          |
| HTTP Framework | [Gin v1.12.0](https://github.com/gin-gonic/gin) |
| Database       | MongoDB — official `go.mongodb.org/mongo-driver` |
| Architecture   | Hexagonal (Ports & Adapters)                |
| Auth           | JWT — `Authorization: Bearer <token>`       |
| Config         | Environment variables + `godotenv` (local)  |
| Logging        | `slog` (stdlib, structured JSON)            |
| Deployment     | Docker Compose (local) + multi-stage Dockerfile (prod) |
| Integration Testing | `testing` + `httptest` + [testify](https://github.com/stretchr/testify) + `testcontainers-go` (MongoDB container)     |


## Project Structure
```
├── cmd/
│   └── main.go                  # wiring: config, DB, routes
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   │   └── order.go         # Order, OrderItem structs + business rules
│   │   ├── ports/
│   │   │   ├── order_repository.go   # OrderRepository interface (driven)
│   │   │   ├── order_repository_test.go # OrderRepository unit tests
│   │   │   ├── product_repository.go # ProductRepository interface (driven)
│   │   │   ├── product_repository_test.go # ProductRepository unit tests
│   │   │   └── order_service.go      # OrderService interface (driving)
│   │   │   └── order_service_test.go      # OrderService unit tests
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

## Workflow of Go developer
1. Understand the requirements and domain model and break down the problem into smaller tasks
2. Plan the architecture and project structure
3. Implement the domain model and business rules
4. Define the ports (interfaces) for repositories and services
5. Implement the service layer with use-case logic
6. Create the HTTP handlers and routes using Gin
7. Connect to MongoDB and implement repository adapters
8. Write integration tests for HTTP handlers and database interactions with testcontainers-go
9. Use Docker for local development and testing
10. Ensure code quality with linters and formatters
11. Document the API and codebase for maintainability


## Rules
- Follow Go best practices and idiomatic code style
- Use dependency injection to decouple components and facilitate testing
- Handle errors gracefully and return appropriate HTTP status codes
- Validate input data and enforce business rules in the service layer
- Write comprehensive tests with good coverage
- Use structured logging for better observability
- Keep the codebase clean and maintainable with proper organization and documentation
- Use environment variables for configuration and avoid hardcoding sensitive information
- Ensure the application can be easily deployed with Docker and Docker Compose