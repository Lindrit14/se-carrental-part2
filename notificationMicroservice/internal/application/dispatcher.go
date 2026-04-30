package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lindritprekaj/notification-service/internal/application/handlers"
	"github.com/lindritprekaj/notification-service/internal/domain/notification"
)

// handlerFn processes a single event envelope.
type handlerFn func(ctx context.Context, env notification.EventEnvelope, send handlers.SendFn) error

// Dispatcher routes event_type to the correct handler.
type Dispatcher struct {
	routes map[string]handlerFn
	send   handlers.SendFn
}

func NewDispatcher(send handlers.SendFn) *Dispatcher {
	d := &Dispatcher{send: send}
	d.routes = map[string]handlerFn{
		"user.registered":                handlers.HandleUserRegistered,
		"user.password_reset_requested":  handlers.HandlePasswordReset,
		"booking.created":                handlers.HandleBookingCreated,
		"booking.cancelled":              handlers.HandleBookingCancelled,
	}
	return d
}

func (d *Dispatcher) Dispatch(ctx context.Context, env notification.EventEnvelope) error {
	h, ok := d.routes[env.EventType]
	if !ok {
		slog.Debug("no handler for event type, skipping", "event_type", env.EventType)
		return nil
	}
	if err := h(ctx, env, d.send); err != nil {
		return fmt.Errorf("handler %s: %w", env.EventType, err)
	}
	return nil
}
