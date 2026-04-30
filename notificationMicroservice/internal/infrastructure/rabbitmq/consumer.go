package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lindritprekaj/notification-service/internal/domain/notification"
)

// DispatchFn is called once per successfully decoded event.
type DispatchFn func(ctx context.Context, env notification.EventEnvelope) error

// Consumer connects to RabbitMQ, declares topology, and consumes messages.
// On connection loss it reconnects with exponential back-off.
type Consumer struct {
	url      string
	dispatch DispatchFn
	seen     *lruSeen // idempotency cache
	ready    bool
	mu       sync.RWMutex
}

func NewConsumer(url string, dispatch DispatchFn) *Consumer {
	return &Consumer{
		url:      url,
		dispatch: dispatch,
		seen:     newLRUSeen(10_000),
	}
}

func (c *Consumer) IsReady() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ready
}

// Run blocks and reconnects forever until ctx is cancelled.
func (c *Consumer) Run(ctx context.Context) {
	backoff := time.Second
	for {
		if err := c.connectAndConsume(ctx); err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("rabbitmq consumer error, reconnecting", "err", err, "backoff", backoff)
			c.setReady(false)
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
			if backoff < 30*time.Second {
				backoff *= 2
			}
		}
	}
}

func (c *Consumer) connectAndConsume(ctx context.Context) error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// One message at a time per consumer to keep ordering guarantees.
	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}

	if err := declareTopology(ch); err != nil {
		return err
	}

	msgs, err := ch.Consume(NotificationsQueue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	c.setReady(true)
	slog.Info("rabbitmq consumer ready", "queue", NotificationsQueue)

	connErr := make(chan *amqp.Error, 1)
	conn.NotifyClose(connErr)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-connErr:
			if err != nil {
				return err
			}
			return nil
		case d, ok := <-msgs:
			if !ok {
				return nil
			}
			c.handle(ctx, d)
		}
	}
}

func (c *Consumer) handle(ctx context.Context, d amqp.Delivery) {
	var env notification.EventEnvelope
	if err := json.Unmarshal(d.Body, &env); err != nil {
		slog.Error("failed to parse event envelope, nacking", "err", err)
		_ = d.Nack(false, false) // don't requeue malformed messages
		return
	}

	// Idempotency: skip already-processed events.
	if env.EventID != "" && c.seen.has(env.EventID) {
		slog.Debug("duplicate event, skipping", "event_id", env.EventID)
		_ = d.Ack(false)
		return
	}

	if err := c.dispatch(ctx, env); err != nil {
		slog.Error("handler error, nacking for redelivery", "event_type", env.EventType, "err", err)
		_ = d.Nack(false, true) // requeue; DLX takes over after x-delivery-limit
		return
	}

	if env.EventID != "" {
		c.seen.add(env.EventID)
	}
	_ = d.Ack(false)
}

func (c *Consumer) setReady(r bool) {
	c.mu.Lock()
	c.ready = r
	c.mu.Unlock()
}
