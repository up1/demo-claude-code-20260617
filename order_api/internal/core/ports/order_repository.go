package ports

import (
	"context"

	"api/internal/core/domain"
)

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindByIdempotencyKey(ctx context.Context, key string) (*domain.Order, error)
}
