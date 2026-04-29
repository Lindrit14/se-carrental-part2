package ports

import "context"

// Event is a transport-agnostic envelope that the publisher serializes.
type Event struct {
	ID         string
	Type       string
	Version    string
	OccurredAt string // RFC3339
	Data       any
	RoutingKey string
}

// EventPublisher publishes domain events to the message bus.
// Adapters (e.g. RabbitMQ) implement this.
type EventPublisher interface {
	Publish(ctx context.Context, e Event) error
}
