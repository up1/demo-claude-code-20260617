package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"api/internal/adapters/handler"
	"api/internal/core/domain"
	"api/internal/core/ports"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// --- mock ---

type mockOrderService struct{ mock.Mock }

func (m *mockOrderService) CreateOrder(ctx context.Context, req ports.CreateOrderInput) (*ports.CreateOrderResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.CreateOrderResult), args.Error(1)
}

// --- helpers ---

func newRouter(svc ports.OrderService) *gin.Engine {
	r := gin.New()
	h := handler.NewOrderHandler(svc)
	r.POST("/api/v1/orders", h.CreateOrder)
	return r
}

func postJSON(r *gin.Engine, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func sampleOrder() *domain.Order {
	return &domain.Order{
		ID:         primitive.NewObjectID(),
		CustomerID: "cust_123",
		Items:      []domain.OrderItem{{ProductID: "prod_abc", Quantity: 2, UnitPrice: 25.00}},
		TotalAmount: 50.00,
		Status:     domain.StatusPending,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

// --- tests ---

func TestCreateOrder_Returns201_NewOrder(t *testing.T) {
	svc := new(mockOrderService)
	r := newRouter(svc)

	body := map[string]interface{}{
		"customer_id": "cust_123",
		"items":       []map[string]interface{}{{"product_id": "prod_abc", "quantity": 2}},
	}
	headers := map[string]string{"Idempotency-Key": "key-new-1"}

	svc.On("CreateOrder", mock.Anything, ports.CreateOrderInput{
		CustomerID:     "cust_123",
		Items:          []ports.OrderItemInput{{ProductID: "prod_abc", Quantity: 2}},
		IdempotencyKey: "key-new-1",
	}).Return(&ports.CreateOrderResult{Order: sampleOrder(), IsExisting: false}, nil)

	w := postJSON(r, body, headers)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.NotNil(t, resp["data"])
	svc.AssertExpectations(t)
}

func TestCreateOrder_Returns200_IdempotentRepeat(t *testing.T) {
	svc := new(mockOrderService)
	r := newRouter(svc)

	body := map[string]interface{}{
		"customer_id": "cust_123",
		"items":       []map[string]interface{}{{"product_id": "prod_abc", "quantity": 2}},
	}
	headers := map[string]string{"Idempotency-Key": "key-existing"}

	svc.On("CreateOrder", mock.Anything, mock.Anything).
		Return(&ports.CreateOrderResult{Order: sampleOrder(), IsExisting: true}, nil)

	w := postJSON(r, body, headers)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
}

func TestCreateOrder_Returns400_MissingIdempotencyKey(t *testing.T) {
	svc := new(mockOrderService)
	r := newRouter(svc)

	body := map[string]interface{}{
		"customer_id": "cust_123",
		"items":       []map[string]interface{}{{"product_id": "prod_abc", "quantity": 1}},
	}

	w := postJSON(r, body, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp["success"].(bool))
	errField := resp["error"].(map[string]interface{})
	assert.Equal(t, "VALIDATION_ERROR", errField["code"])
	svc.AssertNotCalled(t, "CreateOrder")
}

func TestCreateOrder_Returns400_MissingCustomerID(t *testing.T) {
	svc := new(mockOrderService)
	r := newRouter(svc)

	body := map[string]interface{}{
		"items": []map[string]interface{}{{"product_id": "prod_abc", "quantity": 1}},
	}
	headers := map[string]string{"Idempotency-Key": "key-1"}

	w := postJSON(r, body, headers)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "CreateOrder")
}

func TestCreateOrder_Returns400_EmptyItems(t *testing.T) {
	svc := new(mockOrderService)
	r := newRouter(svc)

	body := map[string]interface{}{
		"customer_id": "cust_123",
		"items":       []map[string]interface{}{},
	}
	headers := map[string]string{"Idempotency-Key": "key-1"}

	w := postJSON(r, body, headers)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "CreateOrder")
}

func TestCreateOrder_Returns404_ProductNotFound(t *testing.T) {
	svc := new(mockOrderService)
	r := newRouter(svc)

	body := map[string]interface{}{
		"customer_id": "cust_123",
		"items":       []map[string]interface{}{{"product_id": "prod_unknown", "quantity": 1}},
	}
	headers := map[string]string{"Idempotency-Key": "key-1"}

	svc.On("CreateOrder", mock.Anything, mock.Anything).
		Return(nil, domain.ErrProductNotFound)

	w := postJSON(r, body, headers)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	errField := resp["error"].(map[string]interface{})
	assert.Equal(t, "PRODUCT_NOT_FOUND", errField["code"])
}
