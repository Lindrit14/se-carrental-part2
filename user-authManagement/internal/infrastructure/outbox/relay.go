package outbox

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lindritprekaj/user-authmanagement/internal/infrastructure/mongodb"
	"github.com/lindritprekaj/user-authmanagement/internal/infrastructure/rabbitmq"
)

// Relay drains the outbox collection to RabbitMQ. Runs as a goroutine; on
// PollInterval ticks it fetches up to BatchSize pending events and publishes
// each one. Failures are recorded on the row (attempts++, last_error) and
// retried next tick.
type Relay struct {
	repo         *mongodb.OutboxRepository
	conn         *rabbitmq.Connection
	logger       *slog.Logger
	pollInterval time.Duration
	batchSize    int
	maxAttempts  int
}

type Config struct {
	PollInterval time.Duration
	BatchSize    int
	MaxAttempts  int
}

func NewRelay(repo *mongodb.OutboxRepository, conn *rabbitmq.Connection, logger *slog.Logger, cfg Config) *Relay {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = time.Second
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 50
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 5
	}
	return &Relay{
		repo:         repo,
		conn:         conn,
		logger:       logger,
		pollInterval: cfg.PollInterval,
		batchSize:    cfg.BatchSize,
		maxAttempts:  cfg.MaxAttempts,
	}
}

// Run blocks until ctx is cancelled. Safe to start before RabbitMQ is fully
// ready — the publish call will simply fail and the row stays unpublished
// for the next tick.
func (r *Relay) Run(ctx context.Context) {
	t := time.NewTicker(r.pollInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			r.drainOnce(ctx)
		}
	}
}

func (r *Relay) drainOnce(ctx context.Context) {
	events, err := r.repo.FindPending(ctx, r.batchSize)
	if err != nil {
		r.logger.Error("outbox find pending failed", slog.Any("error", err))
		return
	}
	if len(events) == 0 {
		return
	}

	ch := r.conn.Channel()
	if ch == nil {
		r.logger.Warn("outbox: rabbit channel not ready, retrying next tick", slog.Int("pending", len(events)))
		return
	}

	for i := range events {
		e := &events[i]
		if err := r.publish(ctx, ch, e); err != nil {
			level := slog.LevelWarn
			if e.Attempts+1 >= r.maxAttempts {
				level = slog.LevelError
			}
			r.logger.Log(ctx, level, "outbox publish failed",
				slog.String("event_id", e.EventID),
				slog.String("event_type", e.EventType),
				slog.Int("attempts", e.Attempts+1),
				slog.Any("error", err),
			)
			_ = r.repo.MarkFailed(ctx, e.ID, err.Error())
			continue
		}
		if err := r.repo.MarkPublished(ctx, e.ID, time.Now().UTC()); err != nil {
			r.logger.Error("outbox mark-published failed",
				slog.String("event_id", e.EventID), slog.Any("error", err))
		}
	}
}

func (r *Relay) publish(ctx context.Context, ch *amqp.Channel, e *mongodb.OutboxEvent) error {
	envelope := map[string]any{
		"event_id":    e.EventID,
		"event_type":  e.EventType,
		"version":     e.Version,
		"occurred_at": e.OccurredAt.UTC().Format(time.RFC3339),
		"data":        json.RawMessage(toJSON(e.Payload)),
	}
	body, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(ctx,
		r.conn.Exchange(),
		e.RoutingKey,
		false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    e.EventID,
			Type:         e.EventType,
			Timestamp:    e.OccurredAt,
			Body:         body,
		},
	)
}

// toJSON converts a bson.Raw payload back to JSON bytes by round-tripping
// through a generic map. The outbox stores the original Data — not the full
// envelope — so we need to wrap it in the envelope at publish time.
func toJSON(raw []byte) []byte {
	var doc map[string]any
	if err := bsonUnmarshal(raw, &doc); err != nil {
		// Should never happen — payload was inserted as bson.Raw of valid bson.
		return []byte("{}")
	}
	out, err := json.Marshal(doc)
	if err != nil {
		return []byte("{}")
	}
	return out
}
