package com.uni.carbooking.domain.customer;

import java.time.Instant;
import java.util.Objects;

/**
 * Read model of a user from user-auth — populated and updated by consuming
 * user.* events from RabbitMQ. Booking never writes back to user-auth.
 */
public class Customer {
    private final String id;                 // own UUID
    private final String externalUserId;     // user.id from user-auth
    private String email;
    private boolean anonymized;
    private final Instant createdAt;
    private Instant updatedAt;

    public Customer(String id, String externalUserId, String email,
                    boolean anonymized, Instant createdAt, Instant updatedAt) {
        this.id = Objects.requireNonNull(id);
        this.externalUserId = Objects.requireNonNull(externalUserId);
        this.email = email;
        this.anonymized = anonymized;
        this.createdAt = Objects.requireNonNull(createdAt);
        this.updatedAt = Objects.requireNonNull(updatedAt);
    }

    public String id() { return id; }
    public String externalUserId() { return externalUserId; }
    public String email() { return email; }
    public boolean anonymized() { return anonymized; }
    public Instant createdAt() { return createdAt; }
    public Instant updatedAt() { return updatedAt; }

    public void updateEmail(String newEmail, Instant now) {
        this.email = newEmail;
        this.updatedAt = now;
    }

    public void anonymize(Instant now) {
        this.email = null;
        this.anonymized = true;
        this.updatedAt = now;
    }
}
