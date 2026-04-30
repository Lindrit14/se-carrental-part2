package com.uni.carbooking.domain.error;

public class CarHasBookings extends DomainException {
    public CarHasBookings(String id) {
        super("car has existing bookings: " + id);
    }
}
