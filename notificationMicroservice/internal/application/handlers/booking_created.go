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

type bookingCreatedData struct {
	BookingID              string `json:"bookingId"`
	CustomerID             string `json:"customerId"`
	CustomerExternalUserID string `json:"customerExternalUserId"`
	CarID                  string `json:"carId"`
	StartDate              string `json:"startDate"`
	EndDate                string `json:"endDate"`
	TotalTargetAmount      string `json:"totalTargetAmount"`
	TotalTargetCurrency    string `json:"totalTargetCurrency"`
}

// HandleBookingCreated looks up the customer's email by the auth user_id
// (carried in customerExternalUserId) and sends a confirmation. The
// internal customerId is kept for diagnostics but is NOT used for the
// lookup — the read model is keyed on auth user_id, populated from
// user.registered events.
//
// Defensive cases:
//   - empty externalUserId  → old-format event (pre-rolling-deploy):
//     warn + ack, no requeue.
//   - read-model miss       → user not yet seen via user.registered:
//     warn + ack, no requeue (would loop forever otherwise).
func HandleBookingCreated(ctx context.Context, env notification.EventEnvelope, deps *Deps) error {
	var d bookingCreatedData
	if err := json.Unmarshal(env.Data, &d); err != nil {
		return fmt.Errorf("parse booking.created data: %w", err)
	}

	if d.CustomerExternalUserID == "" {
		slog.Warn("booking.created: missing customerExternalUserId, skipping email",
			"booking_id", d.BookingID, "customer_id", d.CustomerID)
		return nil
	}

	email, err := deps.UserRepo.GetEmail(ctx, d.CustomerExternalUserID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			slog.Warn("booking.created: customer email not in read model, skipping email",
				"booking_id", d.BookingID,
				"customer_id", d.CustomerID,
				"external_user_id", d.CustomerExternalUserID)
			return nil
		}
		return fmt.Errorf("booking.created: lookup email: %w", err)
	}

	htmlBody, textBody, err := deps.Renderer.Render("booking_confirmation", map[string]any{
		"BookingID":     d.BookingID,
		"CarID":         d.CarID,
		"StartDate":     d.StartDate,
		"EndDate":       d.EndDate,
		"TotalAmount":   d.TotalTargetAmount,
		"TotalCurrency": d.TotalTargetCurrency,
	})
	if err != nil {
		return err
	}

	return deps.Notifier.Send(ctx, notifier.Message{
		To:       email,
		Subject:  fmt.Sprintf("Booking %s confirmed", d.BookingID),
		HTMLBody: htmlBody,
		TextBody: textBody,
	})
}
