package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lindritprekaj/notification-service/internal/domain/notification"
	"github.com/lindritprekaj/notification-service/internal/domain/user"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/notifier"
)

type bookingCancelledData struct {
	BookingID              string `json:"bookingId"`
	CustomerID             string `json:"customerId"`
	CustomerExternalUserID string `json:"customerExternalUserId"`
}

func HandleBookingCancelled(ctx context.Context, env notification.EventEnvelope, deps *Deps) error {
	var d bookingCancelledData
	if err := json.Unmarshal(env.Data, &d); err != nil {
		return fmt.Errorf("parse booking.cancelled data: %w", err)
	}

	if d.CustomerExternalUserID == "" {
		slog.Warn("booking.cancelled: missing customerExternalUserId, skipping email",
			"booking_id", d.BookingID, "customer_id", d.CustomerID)
		return nil
	}

	email, err := deps.UserRepo.GetEmail(ctx, d.CustomerExternalUserID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			slog.Warn("booking.cancelled: customer email not in read model, skipping email",
				"booking_id", d.BookingID,
				"customer_id", d.CustomerID,
				"external_user_id", d.CustomerExternalUserID)
			return nil
		}
		return fmt.Errorf("booking.cancelled: lookup email: %w", err)
	}

	htmlBody, textBody, err := deps.Renderer.Render("booking_cancellation", map[string]any{
		"BookingID": d.BookingID,
	})
	if err != nil {
		return err
	}

	return deps.Notifier.Send(ctx, notifier.Message{
		To:       email,
		Subject:  fmt.Sprintf("Booking %s cancelled", d.BookingID),
		HTMLBody: htmlBody,
		TextBody: textBody,
	})
}
