package com.uni.carbooking.application.event;

/**
 * Public event contract published on the booking.events topic exchange.
 * Consumers may rely on these constants and payload shapes.
 */
public final class BookingEvents {

    public static final String VERSION = "v1";

    public static final String BOOKING_CREATED = "booking.created";
    public static final String BOOKING_CANCELLED = "booking.cancelled";

    public record BookingCreated(
            String bookingId,
            String customerId,
            String customerExternalUserId,
            String carId,
            String startDate,
            String endDate,
            String totalSourceAmount,
            String totalSourceCurrency,
            String totalTargetAmount,
            String totalTargetCurrency
    ) {}

    public record BookingCancelled(
            String bookingId,
            String customerId,
            String customerExternalUserId
    ) {}

    private BookingEvents() {}
}
