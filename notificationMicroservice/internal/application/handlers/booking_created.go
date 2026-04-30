package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lindritprekaj/notification-service/internal/domain/notification"
)

type bookingCreatedData struct {
	BookingID          string `json:"bookingId"`
	CustomerID         string `json:"customerId"`
	CarID              string `json:"carId"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	TotalTargetAmount   string `json:"totalTargetAmount"`
	TotalTargetCurrency string `json:"totalTargetCurrency"`
}

func HandleBookingCreated(ctx context.Context, env notification.EventEnvelope, send SendFn) error {
	var d bookingCreatedData
	if err := json.Unmarshal(env.Data, &d); err != nil {
		return fmt.Errorf("parse booking.created data: %w", err)
	}
	send(ctx, "booking_confirmation",
		fmt.Sprintf("Booking confirmation sent for booking_id=%s customer_id=%s (%s → %s, %s %s)",
			d.BookingID, d.CustomerID, d.StartDate, d.EndDate, d.TotalTargetAmount, d.TotalTargetCurrency))
	return nil
}
