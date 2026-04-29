package com.uni.carbooking.interfaces.rest.dto;

import com.uni.carbooking.domain.booking.Booking;

import java.math.BigDecimal;

public record BookingResponse(
        String id,
        String carId,
        String startDate,
        String endDate,
        String status,
        BigDecimal totalSourceAmount,
        String totalSourceCurrency,
        BigDecimal totalTargetAmount,
        String totalTargetCurrency
) {
    public static BookingResponse from(Booking b) {
        return new BookingResponse(
                b.id(),
                b.carId(),
                b.startDate().toString(),
                b.endDate().toString(),
                b.status().name(),
                b.totalSource().amount(),
                b.totalSource().currency(),
                b.totalTarget().amount(),
                b.totalTarget().currency()
        );
    }
}
