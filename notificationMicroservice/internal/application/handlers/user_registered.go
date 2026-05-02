package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lindritprekaj/notification-service/internal/domain/notification"
	"github.com/lindritprekaj/notification-service/internal/domain/user"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/notifier"
)

type userRegisteredData struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// HandleUserRegistered persists the user_id→email mapping into the local
// read model (so booking events can look up the email later) and then
// sends a welcome email.
func HandleUserRegistered(ctx context.Context, env notification.EventEnvelope, deps *Deps) error {
	var d userRegisteredData
	if err := json.Unmarshal(env.Data, &d); err != nil {
		return fmt.Errorf("parse user.registered data: %w", err)
	}
	if d.UserID == "" || d.Email == "" {
		return fmt.Errorf("user.registered: missing user_id or email")
	}

	if err := deps.UserRepo.Upsert(ctx, user.User{
		ID:        d.UserID,
		Email:     d.Email,
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		return fmt.Errorf("user.registered: upsert read model: %w", err)
	}

	htmlBody, textBody, err := deps.Renderer.Render("welcome", map[string]any{
		"Email": d.Email,
	})
	if err != nil {
		return err
	}

	return deps.Notifier.Send(ctx, notifier.Message{
		To:       d.Email,
		Subject:  "Welcome to " + deps.FromName,
		HTMLBody: htmlBody,
		TextBody: textBody,
	})
}
