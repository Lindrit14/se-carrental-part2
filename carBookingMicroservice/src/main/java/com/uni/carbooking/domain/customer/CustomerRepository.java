package com.uni.carbooking.domain.customer;

import java.util.Optional;

public interface CustomerRepository {
    void save(Customer customer);
    Optional<Customer> findByExternalUserId(String externalUserId);
}
