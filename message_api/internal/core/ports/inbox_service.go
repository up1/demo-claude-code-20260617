package ports

import (
	"context"
	"time"

	"api/internal/core/domain"
)

// ListInboxInput is the raw request received by the service before normalization.
type ListInboxInput struct {
	Channel  string
	Status   string
	Search   string
	From     *time.Time
	To       *time.Time
	Page     int
	PageSize int
}

// ListInboxResult is the paginated result returned by the service.
type ListInboxResult struct {
	Items      []domain.InboxMessage
	Page       int
	PageSize   int
	TotalItems int64
	TotalPages int
}

// InboxService is the driving port called by the HTTP handler.
type InboxService interface {
	ListMessages(ctx context.Context, req ListInboxInput) (*ListInboxResult, error)
}
