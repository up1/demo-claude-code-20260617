package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"api/internal/core/domain"
	"api/internal/core/ports"
)

type InboxHandler struct {
	svc ports.InboxService
}

func NewInboxHandler(svc ports.InboxService) *InboxHandler {
	return &InboxHandler{svc: svc}
}

// ListMessages handles GET /api/v1/inbox/messages.
func (h *InboxHandler) ListMessages(c *gin.Context) {
	input := ports.ListInboxInput{
		Channel: c.Query("channel"),
		Status:  c.Query("status"),
		Search:  c.Query("q"),
	}

	if raw := c.Query("from"); raw != "" {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResp("VALIDATION_ERROR", "from must be a valid RFC3339 timestamp"))
			return
		}
		input.From = &t
	}
	if raw := c.Query("to"); raw != "" {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResp("VALIDATION_ERROR", "to must be a valid RFC3339 timestamp"))
			return
		}
		input.To = &t
	}

	if raw := c.Query("page"); raw != "" {
		page, err := strconv.Atoi(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResp("VALIDATION_ERROR", "page must be an integer"))
			return
		}
		input.Page = page
	}
	if raw := c.Query("page_size"); raw != "" {
		pageSize, err := strconv.Atoi(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResp("VALIDATION_ERROR", "page_size must be an integer"))
			return
		}
		input.PageSize = pageSize
	}

	result, err := h.svc.ListMessages(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			c.JSON(http.StatusBadRequest, errorResp("VALIDATION_ERROR", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "internal server error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result.Items,
		"pagination": gin.H{
			"page":        result.Page,
			"page_size":   result.PageSize,
			"total_items": result.TotalItems,
			"total_pages": result.TotalPages,
		},
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
