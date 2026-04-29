package com.uni.carbooking.infrastructure.persistence.jpa;

import com.uni.carbooking.domain.customer.Customer;
import com.uni.carbooking.domain.customer.CustomerRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
class CustomerRepositoryAdapter implements CustomerRepository {

    private final CustomerSpringRepo repo;

    CustomerRepositoryAdapter(CustomerSpringRepo repo) {
        this.repo = repo;
    }

    @Override
    public void save(Customer customer) {
        repo.save(CustomerJpaEntity.fromDomain(customer));
    }

    @Override
    public Optional<Customer> findByExternalUserId(String externalUserId) {
        return repo.findByExternalUserId(externalUserId).map(CustomerJpaEntity::toDomain);
    }
}
