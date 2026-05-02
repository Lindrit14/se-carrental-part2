package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/lindritprekaj/notification-service/internal/domain/notification"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/notifier"
)

type passwordResetData struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	ResetToken string `json:"reset_token"`
	ExpiresAt  string `json:"expires_at"`
}

// HandlePasswordReset builds the reset URL and sends the password-reset
// email. The reset token is sensitive — it must NEVER be logged.
func HandlePasswordReset(ctx context.Context, env notification.EventEnvelope, deps *Deps) error {
	var d passwordResetData
	if err := json.Unmarshal(env.Data, &d); err != nil {
		return fmt.Errorf("parse user.password_reset_requested data: %w", err)
	}
	if d.Email == "" || d.ResetToken == "" {
		return fmt.Errorf("password_reset: missing email or token")
	}

	base := strings.TrimRight(deps.FrontendBaseURL, "/")
	if base == "" {
		return fmt.Errorf("password_reset: FRONTEND_BASE_URL not configured")
	}
	resetURL := fmt.Sprintf("%s/password/confirm?token=%s", base, url.QueryEscape(d.ResetToken))

	htmlBody, textBody, err := deps.Renderer.Render("password_reset", map[string]any{
		"Email":     d.Email,
		"ResetURL":  resetURL,
		"ExpiresAt": d.ExpiresAt,
	})
	if err != nil {
		return err
	}

	return deps.Notifier.Send(ctx, notifier.Message{
		To:       d.Email,
		Subject:  "Reset your password",
		HTMLBody: htmlBody,
		TextBody: textBody,
	})
}
