package com.uni.carbooking.domain.error;

public class CurrencyConversionFailed extends DomainException {
    public CurrencyConversionFailed(String message) {
        super(message);
    }
}
