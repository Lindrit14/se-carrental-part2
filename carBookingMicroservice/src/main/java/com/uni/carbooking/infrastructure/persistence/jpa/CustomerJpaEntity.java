package com.uni.carbooking.infrastructure.persistence.jpa;

import com.uni.carbooking.domain.customer.Customer;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

import java.time.Instant;

@Entity
@Table(name = "customers")
class CustomerJpaEntity {

    @Id
    @Column(length = 36)
    private String id;

    @Column(name = "external_user_id", nullable = false, unique = true, length = 36)
    private String externalUserId;

    @Column
    private String email;

    @Column(nullable = false)
    private boolean anonymized;

    @Column(name = "created_at", nullable = false)
    private Instant createdAt;

    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt;

    protected CustomerJpaEntity() {}

    static CustomerJpaEntity fromDomain(Customer c) {
        var e = new CustomerJpaEntity();
        e.id = c.id();
        e.externalUserId = c.externalUserId();
        e.email = c.email();
        e.anonymized = c.anonymized();
        e.createdAt = c.createdAt();
        e.updatedAt = c.updatedAt();
        return e;
    }

    Customer toDomain() {
        return new Customer(id, externalUserId, email, anonymized, createdAt, updatedAt);
    }
}
