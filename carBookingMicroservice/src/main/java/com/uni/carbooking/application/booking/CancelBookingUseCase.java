package com.uni.carbooking.application.booking;

import com.uni.carbooking.application.event.BookingEvents;
import com.uni.carbooking.application.port.out.Clock;
import com.uni.carbooking.application.port.out.EventPublisher;
import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.booking.BookingRepository;
import com.uni.carbooking.domain.customer.Customer;
import com.uni.carbooking.domain.customer.CustomerRepository;
import com.uni.carbooking.domain.error.BookingNotFound;
import com.uni.carbooking.domain.error.BookingNotOwned;
import com.uni.carbooking.domain.error.CustomerNotFound;

public class CancelBookingUseCase {

    private final BookingRepository bookings;
    private final CustomerRepository customers;
    private final EventPublisher publisher;
    private final Clock clock;

    public CancelBookingUseCase(BookingRepository bookings, CustomerRepository customers,
                                EventPublisher publisher, Clock clock) {
        this.bookings = bookings;
        this.customers = customers;
        this.publisher = publisher;
        this.clock = clock;
    }

    public void execute(String bookingId, String externalUserId) {
        Customer customer = customers.findByExternalUserId(externalUserId)
                .orElseThrow(() -> new CustomerNotFound(externalUserId));
        Booking b = bookings.findById(bookingId)
                .orElseThrow(() -> new BookingNotFound(bookingId));
        if (!b.customerId().equals(customer.id())) {
            throw new BookingNotOwned();
        }
        b.cancel(clock.now());
        bookings.save(b);

        publisher.publish(
                BookingEvents.BOOKING_CANCELLED,
                BookingEvents.BOOKING_CANCELLED,
                new BookingEvents.BookingCancelled(b.id(), b.customerId())
        );
    }
}
