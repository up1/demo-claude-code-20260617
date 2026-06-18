package handler_test

import (
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

type mockInboxService struct{ mock.Mock }

func (m *mockInboxService) ListMessages(ctx context.Context, req ports.ListInboxInput) (*ports.ListInboxResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.ListInboxResult), args.Error(1)
}

func newRouter(svc ports.InboxService) *gin.Engine {
	r := gin.New()
	h := handler.NewInboxHandler(svc)
	r.GET("/api/v1/inbox/messages", h.ListMessages)
	return r
}

func getJSON(r *gin.Engine, target string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func sampleMessages() []domain.InboxMessage {
	return []domain.InboxMessage{
		{
			ID:         primitive.NewObjectID(),
			CustomerID: "cust_123",
			SenderName: "Marcus Watanabe",
			Channel:    domain.ChannelLine,
			Preview:    "Can you confirm the tracking number for order #8812?",
			Status:     domain.StatusPending,
			Unread:     true,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		},
	}
}

func TestListMessages_DefaultOK(t *testing.T) {
	svc := new(mockInboxService)
	r := newRouter(svc)

	svc.On("ListMessages", mock.Anything, ports.ListInboxInput{}).
		Return(&ports.ListInboxResult{
			Items: sampleMessages(), Page: 1, PageSize: 20, TotalItems: 1, TotalPages: 1,
		}, nil)

	w := getJSON(r, "/api/v1/inbox/messages")

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Len(t, resp["data"].([]interface{}), 1)
	pg := resp["pagination"].(map[string]interface{})
	assert.Equal(t, float64(1), pg["page"])
	assert.Equal(t, float64(20), pg["page_size"])
	svc.AssertExpectations(t)
}

func TestListMessages_FilterByChannel(t *testing.T) {
	svc := new(mockInboxService)
	r := newRouter(svc)

	svc.On("ListMessages", mock.Anything, ports.ListInboxInput{Channel: "line"}).
		Return(&ports.ListInboxResult{Items: sampleMessages(), Page: 1, PageSize: 20, TotalItems: 1, TotalPages: 1}, nil)

	w := getJSON(r, "/api/v1/inbox/messages?channel=line")

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestListMessages_FilterByStatus(t *testing.T) {
	svc := new(mockInboxService)
	r := newRouter(svc)

	svc.On("ListMessages", mock.Anything, ports.ListInboxInput{Status: "pending"}).
		Return(&ports.ListInboxResult{Items: sampleMessages(), Page: 1, PageSize: 20, TotalItems: 1, TotalPages: 1}, nil)

	w := getJSON(r, "/api/v1/inbox/messages?status=pending")

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestListMessages_CombinedFilterAndSearch(t *testing.T) {
	svc := new(mockInboxService)
	r := newRouter(svc)

	svc.On("ListMessages", mock.Anything, ports.ListInboxInput{Channel: "line", Search: "tracking"}).
		Return(&ports.ListInboxResult{Items: sampleMessages(), Page: 1, PageSize: 20, TotalItems: 1, TotalPages: 1}, nil)

	w := getJSON(r, "/api/v1/inbox/messages?channel=line&q=tracking")

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestListMessages_Pagination(t *testing.T) {
	svc := new(mockInboxService)
	r := newRouter(svc)

	svc.On("ListMessages", mock.Anything, ports.ListInboxInput{Page: 2, PageSize: 10}).
		Return(&ports.ListInboxResult{Items: sampleMessages(), Page: 2, PageSize: 10, TotalItems: 25, TotalPages: 3}, nil)

	w := getJSON(r, "/api/v1/inbox/messages?page=2&page_size=10")

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	pg := resp["pagination"].(map[string]interface{})
	assert.Equal(t, float64(2), pg["page"])
	assert.Equal(t, float64(3), pg["total_pages"])
	svc.AssertExpectations(t)
}

func TestListMessages_EmptyResult(t *testing.T) {
	svc := new(mockInboxService)
	r := newRouter(svc)

	svc.On("ListMessages", mock.Anything, mock.Anything).
		Return(&ports.ListInboxResult{Items: []domain.InboxMessage{}, Page: 1, PageSize: 20, TotalItems: 0, TotalPages: 0}, nil)

	w := getJSON(r, "/api/v1/inbox/messages?channel=line&q=nothingmatches")

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Len(t, resp["data"].([]interface{}), 0)
	pg := resp["pagination"].(map[string]interface{})
	assert.Equal(t, float64(0), pg["total_items"])
}

func TestListMessages_InvalidChannel(t *testing.T) {
	svc := new(mockInboxService)
	r := newRouter(svc)

	svc.On("ListMessages", mock.Anything, ports.ListInboxInput{Channel: "tiktok"}).
		Return(nil, domain.ErrValidation)

	w := getJSON(r, "/api/v1/inbox/messages?channel=tiktok")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp["success"].(bool))
	errField := resp["error"].(map[string]interface{})
	assert.Equal(t, "VALIDATION_ERROR", errField["code"])
}

func TestListMessages_InvalidDateRange(t *testing.T) {
	svc := new(mockInboxService)
	r := newRouter(svc)

	svc.On("ListMessages", mock.Anything, mock.Anything).
		Return(nil, domain.ErrValidation)

	w := getJSON(r, "/api/v1/inbox/messages?from=2026-06-18T00:00:00Z&to=2026-06-17T00:00:00Z")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	errField := resp["error"].(map[string]interface{})
	assert.Equal(t, "VALIDATION_ERROR", errField["code"])
}

func TestListMessages_InvalidDateFormat(t *testing.T) {
	svc := new(mockInboxService)
	r := newRouter(svc)

	w := getJSON(r, "/api/v1/inbox/messages?from=not-a-date")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "ListMessages")
}
