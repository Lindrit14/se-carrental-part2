package com.uni.carbooking.domain.error;

public class BookingNotOwned extends DomainException {
    public BookingNotOwned() {
        super("booking does not belong to caller");
    }
}
