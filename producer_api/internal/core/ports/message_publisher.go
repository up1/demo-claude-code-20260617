package ports

import "context"

// MessagePublisher is a driven port that publishes a raw message body to a
// named queue (RabbitMQ work-queue pattern).
type MessagePublisher interface {
	Publish(ctx context.Context, queue string, body []byte) error
}
