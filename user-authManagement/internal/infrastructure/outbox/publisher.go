// Package outbox implements the transactional-outbox pattern for the
// user-auth service. The Publisher writes events to a Mongo collection
// (ports.EventPublisher), while Relay drains that collection to RabbitMQ
// asynchronously.
//
// On standalone Mongo (dev) the user write and outbox write are not
// transactional. ADR-0004 documents that trade-off; in production this
// service runs against a replica set where multi-document transactions
// could be added if stronger atomicity becomes necessary.
package outbox

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	"github.com/lindritprekaj/user-authmanagement/internal/infrastructure/mongodb"
)

// Publisher implements ports.EventPublisher by writing events to the outbox
// collection in Mongo. It is the only way the service produces events —
// nothing publishes directly to RabbitMQ.
type Publisher struct {
	repo *mongodb.OutboxRepository
}

func NewPublisher(repo *mongodb.OutboxRepository) *Publisher {
	return &Publisher{repo: repo}
}

func (p *Publisher) Publish(ctx context.Context, e ports.Event) error {
	payload, err := mongodb.MustMarshalPayload(e.Data)
	if err != nil {
		return fmt.Errorf("outbox publish: %w", err)
	}

	occurredAt, err := time.Parse(time.RFC3339, e.OccurredAt)
	if err != nil {
		occurredAt = time.Now().UTC()
	}

	row := &mongodb.OutboxEvent{
		EventID:     e.ID,
		AggregateID: aggregateIDFromEvent(e),
		EventType:   e.Type,
		RoutingKey:  e.RoutingKey,
		Version:     e.Version,
		Payload:     payload,
		OccurredAt:  occurredAt,
	}
	return p.repo.Insert(ctx, row)
}

// aggregateIDFromEvent best-effort extracts the user_id from the event payload.
// Falls back to "unknown" so the row still inserts — aggregate_id is for
// observability, not correctness.
func aggregateIDFromEvent(e ports.Event) string {
	if m, ok := e.Data.(map[string]any); ok {
		if v, ok := m["user_id"].(string); ok {
			return v
		}
	}
	v := reflect.ValueOf(e.Data)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		f := v.FieldByName("UserID")
		if f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
	}
	return "unknown"
}
