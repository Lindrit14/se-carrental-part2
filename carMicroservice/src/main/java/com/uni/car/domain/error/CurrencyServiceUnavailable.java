package com.uni.car.domain.error;

public class CurrencyServiceUnavailable extends RuntimeException {
    public CurrencyServiceUnavailable(String message) {
        super(message);
    }
}
