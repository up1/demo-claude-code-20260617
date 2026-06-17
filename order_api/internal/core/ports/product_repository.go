package ports

import (
	"context"

	"api/internal/core/domain"
)

type ProductRepository interface {
	FindByID(ctx context.Context, productID string) (*domain.Product, error)
}
