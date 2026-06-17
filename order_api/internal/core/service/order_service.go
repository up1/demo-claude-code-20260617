package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"api/internal/core/domain"
	"api/internal/core/ports"
)

type orderService struct {
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
}

func NewOrderService(orderRepo ports.OrderRepository, productRepo ports.ProductRepository) ports.OrderService {
	return &orderService{orderRepo: orderRepo, productRepo: productRepo}
}

func (s *orderService) CreateOrder(ctx context.Context, req ports.CreateOrderInput) (*ports.CreateOrderResult, error) {
	existing, err := s.orderRepo.FindByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return &ports.CreateOrderResult{Order: existing, IsExisting: true}, nil
	}

	var items []domain.OrderItem
	var total float64

	for _, item := range req.Items {
		product, err := s.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, fmt.Errorf("%w: %s", domain.ErrProductNotFound, item.ProductID)
		}
		items = append(items, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
		})
		total += float64(item.Quantity) * product.Price
	}

	now := time.Now().UTC()
	order := &domain.Order{
		ID:             primitive.NewObjectID(),
		CustomerID:     req.CustomerID,
		Items:          items,
		TotalAmount:    total,
		Status:         domain.StatusPending,
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return &ports.CreateOrderResult{Order: order, IsExisting: false}, nil
}
