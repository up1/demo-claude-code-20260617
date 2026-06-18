package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"producer/internal/core/domain"
	"producer/internal/core/ports"
)

// LineHandler is the driving HTTP adapter for receiving LINE messages.
type LineHandler struct {
	svc ports.LineService
}

func NewLineHandler(svc ports.LineService) *LineHandler {
	return &LineHandler{svc: svc}
}

// ReceiveMessage handles POST /api/v1/line/messages.
func (h *LineHandler) ReceiveMessage(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("System down"))
		return
	}

	var msg domain.LineMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("Invalid message format"))
		return
	}

	res, err := h.svc.SendMessage(c.Request.Context(), raw, msg)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			c.JSON(http.StatusBadRequest, errorResp("Invalid message format"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResp("System down"))
		return
	}

	if res.Duplicate {
		c.JSON(http.StatusOK, gin.H{"status": "duplicate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"messageId": res.MessageID,
	})
}

func errorResp(message string) gin.H {
	return gin.H{
		"status":  "error",
		"message": message,
	}
}
