package com.uni.carbooking.application.booking;

import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.booking.BookingRepository;

import java.util.List;

/**
 * Admin-only: returns every booking in the system, sorted by creation time
 * (descending). Caller-side authorization is enforced at the REST layer via
 * {@code @PreAuthorize("hasRole('admin')")}.
 */
public class ListAllBookingsUseCase {

    private final BookingRepository bookings;

    public ListAllBookingsUseCase(BookingRepository bookings) {
        this.bookings = bookings;
    }

    public List<Booking> execute() {
        return bookings.findAll();
    }
}
