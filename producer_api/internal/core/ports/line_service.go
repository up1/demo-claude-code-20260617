package ports

import (
	"context"

	"producer/internal/core/domain"
)

// SendResult is the outcome of handling an inbound LINE message.
type SendResult struct {
	MessageID string
	Duplicate bool
}

// LineService is the driving port for the LINE message use-case.
type LineService interface {
	// SendMessage deduplicates against the raw body, validates the parsed
	// message, and publishes it to the work queue. raw is the exact request
	// body bytes used to compute the dedup key.
	SendMessage(ctx context.Context, raw []byte, msg domain.LineMessage) (SendResult, error)
}
