package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"api/internal/core/domain"
	"api/internal/core/ports"
)

type OrderHandler struct {
	svc ports.OrderService
}

func NewOrderHandler(svc ports.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

type createOrderRequest struct {
	CustomerID string         `json:"customer_id" binding:"required"`
	Items      []orderItemReq `json:"items"       binding:"required,min=1,dive"`
}

type orderItemReq struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity"   binding:"required,min=1"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, errorResp("VALIDATION_ERROR", "Idempotency-Key header is required"))
		return
	}

	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("VALIDATION_ERROR", err.Error()))
		return
	}

	items := make([]ports.OrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = ports.OrderItemInput{ProductID: item.ProductID, Quantity: item.Quantity}
	}

	result, err := h.svc.CreateOrder(c.Request.Context(), ports.CreateOrderInput{
		CustomerID:     req.CustomerID,
		Items:          items,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, errorResp("PRODUCT_NOT_FOUND", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "internal server error"))
		return
	}

	status := http.StatusCreated
	if result.IsExisting {
		status = http.StatusOK
	}

	c.JSON(status, gin.H{
		"success": true,
		"data":    result.Order,
	})
}

func errorResp(code, message string) gin.H {
	return gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	}
}
