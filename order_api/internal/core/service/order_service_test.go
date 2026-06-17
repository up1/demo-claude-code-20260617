package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api/internal/core/domain"
	"api/internal/core/ports"
	"api/internal/core/service"
)

// --- mocks ---

type mockOrderRepo struct{ mock.Mock }

func (m *mockOrderRepo) Save(ctx context.Context, order *domain.Order) error {
	return m.Called(ctx, order).Error(0)
}

func (m *mockOrderRepo) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Order, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

type mockProductRepo struct{ mock.Mock }

func (m *mockProductRepo) FindByID(ctx context.Context, productID string) (*domain.Product, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

// --- tests ---

func TestCreateOrder_NewOrder(t *testing.T) {
	orderRepo := new(mockOrderRepo)
	productRepo := new(mockProductRepo)
	svc := service.NewOrderService(orderRepo, productRepo)

	ctx := context.Background()
	input := ports.CreateOrderInput{
		CustomerID:     "cust_123",
		Items:          []ports.OrderItemInput{{ProductID: "prod_abc", Quantity: 2}},
		IdempotencyKey: "idem-new",
	}

	orderRepo.On("FindByIdempotencyKey", ctx, "idem-new").Return(nil, nil)
	productRepo.On("FindByID", ctx, "prod_abc").Return(&domain.Product{ID: "prod_abc", Price: 25.00}, nil)
	orderRepo.On("Save", ctx, mock.AnythingOfType("*domain.Order")).Return(nil)

	result, err := svc.CreateOrder(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsExisting)
	assert.Equal(t, 50.00, result.Order.TotalAmount)
	assert.Equal(t, domain.StatusPending, result.Order.Status)
	assert.Equal(t, "cust_123", result.Order.CustomerID)
	orderRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCreateOrder_IdempotentReturn(t *testing.T) {
	orderRepo := new(mockOrderRepo)
	productRepo := new(mockProductRepo)
	svc := service.NewOrderService(orderRepo, productRepo)

	ctx := context.Background()
	existing := &domain.Order{CustomerID: "cust_123", Status: domain.StatusPending}
	input := ports.CreateOrderInput{
		CustomerID:     "cust_123",
		Items:          []ports.OrderItemInput{{ProductID: "prod_abc", Quantity: 1}},
		IdempotencyKey: "idem-existing",
	}

	orderRepo.On("FindByIdempotencyKey", ctx, "idem-existing").Return(existing, nil)

	result, err := svc.CreateOrder(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsExisting)
	assert.Equal(t, existing, result.Order)
	productRepo.AssertNotCalled(t, "FindByID")
	orderRepo.AssertNotCalled(t, "Save")
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
	orderRepo := new(mockOrderRepo)
	productRepo := new(mockProductRepo)
	svc := service.NewOrderService(orderRepo, productRepo)

	ctx := context.Background()
	input := ports.CreateOrderInput{
		CustomerID:     "cust_123",
		Items:          []ports.OrderItemInput{{ProductID: "prod_unknown", Quantity: 1}},
		IdempotencyKey: "idem-noproduct",
	}

	orderRepo.On("FindByIdempotencyKey", ctx, "idem-noproduct").Return(nil, nil)
	productRepo.On("FindByID", ctx, "prod_unknown").Return(nil, nil)

	result, err := svc.CreateOrder(ctx, input)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
	orderRepo.AssertNotCalled(t, "Save")
}

func TestCreateOrder_MultipleItems_TotalCalculated(t *testing.T) {
	orderRepo := new(mockOrderRepo)
	productRepo := new(mockProductRepo)
	svc := service.NewOrderService(orderRepo, productRepo)

	ctx := context.Background()
	input := ports.CreateOrderInput{
		CustomerID: "cust_456",
		Items: []ports.OrderItemInput{
			{ProductID: "prod_abc", Quantity: 2},
			{ProductID: "prod_xyz", Quantity: 1},
		},
		IdempotencyKey: "idem-multi",
	}

	orderRepo.On("FindByIdempotencyKey", ctx, "idem-multi").Return(nil, nil)
	productRepo.On("FindByID", ctx, "prod_abc").Return(&domain.Product{ID: "prod_abc", Price: 25.00}, nil)
	productRepo.On("FindByID", ctx, "prod_xyz").Return(&domain.Product{ID: "prod_xyz", Price: 10.00}, nil)
	orderRepo.On("Save", ctx, mock.AnythingOfType("*domain.Order")).Return(nil)

	result, err := svc.CreateOrder(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, 60.00, result.Order.TotalAmount) // 2*25 + 1*10
	assert.Len(t, result.Order.Items, 2)
}
