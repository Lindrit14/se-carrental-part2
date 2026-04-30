package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lindritprekaj/notification-service/internal/domain/notification"
)

type userRegisteredData struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func HandleUserRegistered(ctx context.Context, env notification.EventEnvelope, send SendFn) error {
	var d userRegisteredData
	if err := json.Unmarshal(env.Data, &d); err != nil {
		return fmt.Errorf("parse user.registered data: %w", err)
	}
	send(ctx, "welcome_email", fmt.Sprintf("Welcome email sent to %s (user_id=%s)", d.Email, d.UserID))
	return nil
}
