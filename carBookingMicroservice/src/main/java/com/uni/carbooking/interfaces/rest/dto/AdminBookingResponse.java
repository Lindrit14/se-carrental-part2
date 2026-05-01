package com.uni.carbooking.interfaces.rest.dto;

import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.money.Money;

import java.math.BigDecimal;

/**
 * Admin projection of a Booking. Includes the internal customerId so admins
 * can correlate bookings with customers.
 */
public record AdminBookingResponse(
        String id,
        String customerId,
        String carId,
        String startDate,
        String endDate,
        String status,
        BigDecimal totalSourceAmount,
        String totalSourceCurrency,
        BigDecimal totalTargetAmount,
        String totalTargetCurrency,
        String createdAt,
        String updatedAt
) {
    public static AdminBookingResponse from(Booking b) {
        return fromWithDisplayTotal(b, b.totalTarget());
    }

    public static AdminBookingResponse fromWithDisplayTotal(Booking b, Money displayTotal) {
        return new AdminBookingResponse(
                b.id(),
                b.customerId(),
                b.carId(),
                b.startDate().toString(),
                b.endDate().toString(),
                b.status().name(),
                b.totalSource().amount(),
                b.totalSource().currency(),
                displayTotal.amount(),
                displayTotal.currency(),
                b.createdAt().toString(),
                b.updatedAt().toString()
        );
    }
}
