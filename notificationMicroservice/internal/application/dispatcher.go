package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lindritprekaj/notification-service/internal/application/handlers"
	"github.com/lindritprekaj/notification-service/internal/domain/notification"
)

// Dispatcher routes event_type to the correct handler with shared deps.
type Dispatcher struct {
	routes map[string]handlers.HandlerFunc
	deps   *handlers.Deps
}

func NewDispatcher(deps *handlers.Deps) *Dispatcher {
	return &Dispatcher{
		deps: deps,
		routes: map[string]handlers.HandlerFunc{
			"user.registered":               handlers.HandleUserRegistered,
			"user.password_reset_requested": handlers.HandlePasswordReset,
			"booking.created":               handlers.HandleBookingCreated,
			"booking.cancelled":             handlers.HandleBookingCancelled,
		},
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, env notification.EventEnvelope) error {
	h, ok := d.routes[env.EventType]
	if !ok {
		slog.Debug("no handler for event type, skipping", "event_type", env.EventType)
		return nil
	}
	if err := h(ctx, env, d.deps); err != nil {
		return fmt.Errorf("handler %s: %w", env.EventType, err)
	}
	return nil
}
