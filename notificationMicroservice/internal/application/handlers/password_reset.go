package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lindritprekaj/notification-service/internal/domain/notification"
)

type passwordResetData struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	ResetToken string `json:"reset_token"`
	ExpiresAt  string `json:"expires_at"`
}

func HandlePasswordReset(ctx context.Context, env notification.EventEnvelope, send SendFn) error {
	var d passwordResetData
	if err := json.Unmarshal(env.Data, &d); err != nil {
		return fmt.Errorf("parse user.password_reset_requested data: %w", err)
	}
	send(ctx, "password_reset_email",
		fmt.Sprintf("Password reset email sent to %s (token=%s expires=%s)", d.Email, d.ResetToken, d.ExpiresAt))
	return nil
}
