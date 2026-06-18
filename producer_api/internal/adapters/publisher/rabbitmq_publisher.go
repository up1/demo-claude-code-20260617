package publisher

import (
	"context"
	"fmt"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

// RabbitMQPublisher implements ports.MessagePublisher over the RabbitMQ AMQP 1.0
// client. It declares a single durable (quorum) work queue and publishes
// directly to it (default-exchange / work-queue pattern).
type RabbitMQPublisher struct {
	env       *rmq.Environment
	conn      *rmq.AmqpConnection
	publisher *rmq.Publisher
	queue     string
}

// NewRabbitMQPublisher connects to RabbitMQ, declares the queue, and prepares a
// publisher targeting it.
func NewRabbitMQPublisher(ctx context.Context, uri, queue string) (*RabbitMQPublisher, error) {
	env := rmq.NewEnvironment(uri, nil)

	conn, err := env.NewConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("open rabbitmq connection: %w", err)
	}

	if _, err := conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{Name: queue}); err != nil {
		_ = env.CloseConnections(ctx)
		return nil, fmt.Errorf("declare queue %q: %w", queue, err)
	}

	pub, err := conn.NewPublisher(ctx, &rmq.QueueAddress{Queue: queue}, nil)
	if err != nil {
		_ = env.CloseConnections(ctx)
		return nil, fmt.Errorf("create publisher: %w", err)
	}

	return &RabbitMQPublisher{env: env, conn: conn, publisher: pub, queue: queue}, nil
}

// Publish sends body to the configured work queue and verifies the broker
// accepted it.
func (p *RabbitMQPublisher) Publish(ctx context.Context, queue string, body []byte) error {
	if queue != p.queue {
		return fmt.Errorf("publisher bound to queue %q, got %q", p.queue, queue)
	}

	result, err := p.publisher.Publish(ctx, rmq.NewMessage(body))
	if err != nil {
		return fmt.Errorf("publish to %q: %w", queue, err)
	}

	switch result.Outcome.(type) {
	case *rmq.StateAccepted:
		return nil
	case *rmq.StateReleased:
		return fmt.Errorf("message to %q was not routed (released)", queue)
	case *rmq.StateRejected:
		return fmt.Errorf("message to %q was rejected by the broker", queue)
	default:
		return fmt.Errorf("message to %q has unexpected outcome: %v", queue, result.Outcome)
	}
}

// Close releases the publisher and underlying connection.
func (p *RabbitMQPublisher) Close(ctx context.Context) error {
	if p.publisher != nil {
		_ = p.publisher.Close(ctx)
	}
	return p.env.CloseConnections(ctx)
}
