package com.uni.carbooking.application.customer;

import com.uni.carbooking.application.port.out.Clock;
import com.uni.carbooking.domain.customer.CustomerRepository;

/**
 * Triggered by the {@code user.updated} event when the email field changes.
 */
public class UpdateCustomerEmailUseCase {

    private final CustomerRepository customers;
    private final Clock clock;

    public UpdateCustomerEmailUseCase(CustomerRepository customers, Clock clock) {
        this.customers = customers;
        this.clock = clock;
    }

    public void execute(String externalUserId, String newEmail) {
        customers.findByExternalUserId(externalUserId).ifPresent(c -> {
            c.updateEmail(newEmail, clock.now());
            customers.save(c);
        });
    }
}
