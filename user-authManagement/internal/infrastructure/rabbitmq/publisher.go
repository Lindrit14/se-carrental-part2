package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
)

// Publisher adapts *Connection to the ports.EventPublisher interface.
type Publisher struct{ conn *Connection }

func NewPublisher(conn *Connection) *Publisher { return &Publisher{conn: conn} }

type envelope struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	Version    string `json:"version"`
	OccurredAt string `json:"occurred_at"`
	Data       any    `json:"data"`
}

func (p *Publisher) Publish(ctx context.Context, e ports.Event) error {
	ch := p.conn.Channel()
	if ch == nil {
		return fmt.Errorf("rabbitmq: channel unavailable")
	}
	body, err := json.Marshal(envelope{
		EventID:    e.ID,
		EventType:  e.Type,
		Version:    e.Version,
		OccurredAt: e.OccurredAt,
		Data:       e.Data,
	})
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	return ch.PublishWithContext(ctx,
		p.conn.Exchange(),
		e.RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    e.ID,
			Type:         e.Type,
			Timestamp:    timeFromRFC3339(e.OccurredAt),
			Body:         body,
		},
	)
}
