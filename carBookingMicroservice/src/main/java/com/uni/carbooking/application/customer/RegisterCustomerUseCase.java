package com.uni.carbooking.application.customer;

import com.uni.carbooking.application.port.out.Clock;
import com.uni.carbooking.application.port.out.IdGenerator;
import com.uni.carbooking.domain.customer.Customer;
import com.uni.carbooking.domain.customer.CustomerRepository;

/**
 * Triggered by the {@code user.registered} event from user-auth.
 * Idempotent: re-running with the same external user id is a no-op.
 */
public class RegisterCustomerUseCase {

    private final CustomerRepository customers;
    private final IdGenerator ids;
    private final Clock clock;

    public RegisterCustomerUseCase(CustomerRepository customers, IdGenerator ids, Clock clock) {
        this.customers = customers;
        this.ids = ids;
        this.clock = clock;
    }

    public void execute(String externalUserId, String email) {
        if (customers.findByExternalUserId(externalUserId).isPresent()) {
            return;
        }
        var now = clock.now();
        var c = new Customer(ids.newId(), externalUserId, email, false, now, now);
        customers.save(c);
    }
}
