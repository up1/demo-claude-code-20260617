package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"producer/internal/adapters/handler"
	"producer/internal/core/domain"
	"producer/internal/core/ports"
)

func init() { gin.SetMode(gin.TestMode) }

type mockService struct{ mock.Mock }

func (m *mockService) SendMessage(ctx context.Context, raw []byte, msg domain.LineMessage) (ports.SendResult, error) {
	args := m.Called(ctx, raw, msg)
	return args.Get(0).(ports.SendResult), args.Error(1)
}

func newRouter(svc ports.LineService) *gin.Engine {
	r := gin.New()
	h := handler.NewLineHandler(svc)
	r.POST("/api/v1/line/messages", h.ReceiveMessage)
	return r
}

func postJSON(r *gin.Engine, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/line/messages", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

const validBody = `{"to":"U123","messages":[{"type":"text","text":"Hello, world!"}]}`

func TestReceiveMessage_Success(t *testing.T) {
	svc := new(mockService)
	svc.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(ports.SendResult{MessageID: "1234567890"}, nil)

	w := postJSON(newRouter(svc), validBody)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "success", resp["status"])
	assert.Equal(t, "1234567890", resp["messageId"])
	svc.AssertExpectations(t)
}

func TestReceiveMessage_Duplicate(t *testing.T) {
	svc := new(mockService)
	svc.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(ports.SendResult{Duplicate: true}, nil)

	w := postJSON(newRouter(svc), validBody)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "duplicate", resp["status"])
	svc.AssertExpectations(t)
}

func TestReceiveMessage_InvalidJSON(t *testing.T) {
	svc := new(mockService)
	w := postJSON(newRouter(svc), `{not json`)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "error", resp["status"])
	assert.Equal(t, "Invalid message format", resp["message"])
	svc.AssertNotCalled(t, "SendMessage")
}

func TestReceiveMessage_ValidationError(t *testing.T) {
	svc := new(mockService)
	svc.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(ports.SendResult{}, domain.ErrValidation)

	w := postJSON(newRouter(svc), `{"to":"","messages":[]}`)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "error", resp["status"])
	assert.Equal(t, "Invalid message format", resp["message"])
	svc.AssertExpectations(t)
}

func TestReceiveMessage_SystemDown(t *testing.T) {
	svc := new(mockService)
	svc.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(ports.SendResult{}, errors.New("broker down"))

	w := postJSON(newRouter(svc), validBody)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "error", resp["status"])
	assert.Equal(t, "System down", resp["message"])
	svc.AssertExpectations(t)
}
