package ports

import (
	"context"

	"api/internal/core/domain"
)

type OrderItemInput struct {
	ProductID string
	Quantity  int
}

type CreateOrderInput struct {
	CustomerID     string
	Items          []OrderItemInput
	IdempotencyKey string
}

type CreateOrderResult struct {
	Order      *domain.Order
	IsExisting bool
}

type OrderService interface {
	CreateOrder(ctx context.Context, req CreateOrderInput) (*CreateOrderResult, error)
}
