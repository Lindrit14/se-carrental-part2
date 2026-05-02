package handlers

import (
	"context"

	"github.com/lindritprekaj/notification-service/internal/domain/notification"
	"github.com/lindritprekaj/notification-service/internal/domain/user"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/email"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/notifier"
)

// Deps groups every collaborator a handler may need.
// Concrete handlers pull what they require and ignore the rest.
type Deps struct {
	Notifier        notifier.Notifier
	UserRepo        user.Repository
	Renderer        *email.Renderer
	FromName        string
	FrontendBaseURL string
}

// HandlerFunc is the new handler signature used by the dispatcher.
type HandlerFunc func(ctx context.Context, env notification.EventEnvelope, deps *Deps) error
