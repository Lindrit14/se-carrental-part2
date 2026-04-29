package com.uni.carbooking.application.booking;

import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.booking.BookingRepository;
import com.uni.carbooking.domain.customer.Customer;
import com.uni.carbooking.domain.customer.CustomerRepository;
import com.uni.carbooking.domain.error.CustomerNotFound;

import java.util.List;

public class ListMyBookingsUseCase {

    private final BookingRepository bookings;
    private final CustomerRepository customers;

    public ListMyBookingsUseCase(BookingRepository bookings, CustomerRepository customers) {
        this.bookings = bookings;
        this.customers = customers;
    }

    public List<Booking> execute(String externalUserId) {
        Customer c = customers.findByExternalUserId(externalUserId)
                .orElseThrow(() -> new CustomerNotFound(externalUserId));
        return bookings.findByCustomerId(c.id());
    }
}
