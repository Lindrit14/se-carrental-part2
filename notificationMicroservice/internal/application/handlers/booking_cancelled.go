package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lindritprekaj/notification-service/internal/domain/notification"
)

type bookingCancelledData struct {
	BookingID  string `json:"bookingId"`
	CustomerID string `json:"customerId"`
}

func HandleBookingCancelled(ctx context.Context, env notification.EventEnvelope, send SendFn) error {
	var d bookingCancelledData
	if err := json.Unmarshal(env.Data, &d); err != nil {
		return fmt.Errorf("parse booking.cancelled data: %w", err)
	}
	send(ctx, "booking_cancellation",
		fmt.Sprintf("Booking cancellation sent for booking_id=%s customer_id=%s", d.BookingID, d.CustomerID))
	return nil
}
