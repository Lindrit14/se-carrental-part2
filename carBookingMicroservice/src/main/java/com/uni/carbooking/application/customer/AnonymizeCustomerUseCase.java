package com.uni.carbooking.application.customer;

import com.uni.carbooking.application.port.out.Clock;
import com.uni.carbooking.domain.customer.CustomerRepository;

/**
 * Triggered by the {@code user.deleted} event from user-auth (GDPR cleanup).
 * Bookings are kept but PII is removed.
 */
public class AnonymizeCustomerUseCase {

    private final CustomerRepository customers;
    private final Clock clock;

    public AnonymizeCustomerUseCase(CustomerRepository customers, Clock clock) {
        this.customers = customers;
        this.clock = clock;
    }

    public void execute(String externalUserId) {
        customers.findByExternalUserId(externalUserId).ifPresent(c -> {
            c.anonymize(clock.now());
            customers.save(c);
        });
    }
}
