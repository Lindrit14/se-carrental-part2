package com.uni.car.domain.error;

public class DomainException extends RuntimeException {
    public DomainException(String message) {
        super(message);
    }
}
