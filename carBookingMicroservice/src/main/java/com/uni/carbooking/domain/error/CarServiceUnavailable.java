package com.uni.carbooking.domain.error;

public class CarServiceUnavailable extends DomainException {
    public CarServiceUnavailable(String message) {
        super(message);
    }
}
