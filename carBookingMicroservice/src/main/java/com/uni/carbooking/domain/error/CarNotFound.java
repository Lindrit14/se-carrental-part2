package com.uni.carbooking.domain.error;

public class CarNotFound extends DomainException {
    public CarNotFound(String id) {
        super("car not found: " + id);
    }
}
