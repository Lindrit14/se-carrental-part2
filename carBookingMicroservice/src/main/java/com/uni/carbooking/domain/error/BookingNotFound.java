package com.uni.carbooking.domain.error;

public class BookingNotFound extends DomainException {
    public BookingNotFound(String id) {
        super("booking not found: " + id);
    }
}
