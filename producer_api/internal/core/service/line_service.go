package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"producer/internal/core/domain"
	"producer/internal/core/ports"
)

// QueueName is the RabbitMQ work queue for LINE messages (per spec).
const QueueName = "line-messages"

const dedupKeyPrefix = "line:dedup:"

// queueEnvelope is the exact message shape published to RabbitMQ (per spec
// section "Message format in RabbitMQ queue").
type queueEnvelope struct {
	Name    string             `json:"name"`
	Message domain.LineMessage `json:"message"`
}

// LineService implements ports.LineService.
type LineService struct {
	publisher ports.MessagePublisher
	dedup     ports.DedupStore
	tracer    trace.Tracer
}

// NewLineService wires the use-case with its driven adapters.
func NewLineService(publisher ports.MessagePublisher, dedup ports.DedupStore) *LineService {
	return &LineService{
		publisher: publisher,
		dedup:     dedup,
		tracer:    otel.Tracer("producer/line-service"),
	}
}

// SendMessage runs the dedup-check → validate → publish flow.
func (s *LineService) SendMessage(ctx context.Context, raw []byte, msg domain.LineMessage) (ports.SendResult, error) {
	ctx, span := s.tracer.Start(ctx, "line.send_message")
	defer span.End()

	log := slog.With("to", msg.To, "message_count", len(msg.Messages))

	// 1. Duplicate detection on a hash of the raw body.
	key := dedupKeyPrefix + hashBody(raw)
	span.SetAttributes(attribute.String("dedup.key", key))

	duplicate, err := s.dedup.IsDuplicate(ctx, key)
	if err != nil {
		return s.fail(span, log, "dedup check failed", fmt.Errorf("dedup check: %w", err))
	}
	if duplicate {
		log.Info("duplicate message ignored", "dedup_key", key)
		span.SetAttributes(attribute.Bool("message.duplicate", true))
		return ports.SendResult{Duplicate: true}, nil
	}

	// 2. Validate the message schema.
	if err := msg.Validate(); err != nil {
		log.Warn("message validation failed", "error", err)
		span.SetStatus(codes.Error, "validation failed")
		span.RecordError(err)
		return ports.SendResult{}, err
	}

	// 3. Publish the envelope to the work queue.
	body, err := json.Marshal(queueEnvelope{Name: QueueName, Message: msg})
	if err != nil {
		return s.fail(span, log, "marshal envelope failed", fmt.Errorf("marshal envelope: %w", err))
	}
	if err := s.publisher.Publish(ctx, QueueName, body); err != nil {
		return s.fail(span, log, "publish failed", fmt.Errorf("publish: %w", err))
	}

	// 4. Generate the message id returned to the caller.
	messageID := uuid.NewString()
	span.SetAttributes(attribute.String("message.id", messageID))
	log.Info("message published", "message_id", messageID, "queue", QueueName)

	return ports.SendResult{MessageID: messageID}, nil
}

func (s *LineService) fail(span trace.Span, log *slog.Logger, msg string, err error) (ports.SendResult, error) {
	log.Error(msg, "error", err)
	span.SetStatus(codes.Error, msg)
	span.RecordError(err)
	return ports.SendResult{}, err
}

func hashBody(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
