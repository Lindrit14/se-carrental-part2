package com.uni.carbooking.interfaces.rest.dto;

import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.money.Money;

import java.math.BigDecimal;

public record BookingResponse(
        String id,
        String carId,
        String carBrand,
        String carModel,
        String carLicensePlate,
        String startDate,
        String endDate,
        String status,
        BigDecimal totalSourceAmount,
        String totalSourceCurrency,
        BigDecimal totalTargetAmount,
        String totalTargetCurrency
) {
    /** Uses the stored locked target. Used for create / cancel responses. */
    public static BookingResponse from(Booking b) {
        return fromWithDisplayTotal(b, b.totalTarget());
    }

    /** Uses {@code displayTotal} as the rendered "Total" — the listing endpoints
     *  pass the booking reconverted into the user's currently selected
     *  display currency, with fallback to the locked target on RPC failure. */
    public static BookingResponse fromWithDisplayTotal(Booking b, Money displayTotal) {
        return new BookingResponse(
                b.id(),
                b.carId(),
                b.carSnapshot().brand(),
                b.carSnapshot().model(),
                b.carSnapshot().licensePlate(),
                b.startDate().toString(),
                b.endDate().toString(),
                b.status().name(),
                b.totalSource().amount(),
                b.totalSource().currency(),
                displayTotal.amount(),
                displayTotal.currency()
        );
    }
}
