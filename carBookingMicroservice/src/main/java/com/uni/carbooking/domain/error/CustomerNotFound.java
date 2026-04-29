package com.uni.carbooking.domain.error;

public class CustomerNotFound extends DomainException {
    public CustomerNotFound(String externalUserId) {
        super("customer not found for external user: " + externalUserId);
    }
}
