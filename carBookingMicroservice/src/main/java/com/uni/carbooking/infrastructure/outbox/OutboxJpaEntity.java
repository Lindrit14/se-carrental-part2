package com.uni.carbooking.infrastructure.outbox;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import java.time.Instant;
import java.util.UUID;

@Entity
@Table(name = "outbox")
class OutboxJpaEntity {

    @Id
    @Column(name = "id")
    private UUID id;

    @Column(name = "event_id", nullable = false, unique = true)
    private String eventId;

    @Column(name = "aggregate_id", nullable = false)
    private String aggregateId;

    @Column(name = "event_type", nullable = false)
    private String eventType;

    @Column(name = "routing_key", nullable = false)
    private String routingKey;

    @JdbcTypeCode(SqlTypes.JSON)
    @Column(name = "payload", nullable = false, columnDefinition = "jsonb")
    private String payload;

    @Column(name = "occurred_at", nullable = false)
    private Instant occurredAt;

    @Column(name = "published_at")
    private Instant publishedAt;

    @Column(name = "attempts", nullable = false)
    private int attempts;

    @Column(name = "last_error")
    private String lastError;

    protected OutboxJpaEntity() {}

    static OutboxJpaEntity create(String eventId, String aggregateId, String eventType,
                                  String routingKey, String payload, Instant occurredAt) {
        var e = new OutboxJpaEntity();
        e.id = UUID.randomUUID();
        e.eventId = eventId;
        e.aggregateId = aggregateId;
        e.eventType = eventType;
        e.routingKey = routingKey;
        e.payload = payload;
        e.occurredAt = occurredAt;
        e.attempts = 0;
        return e;
    }

    UUID getId() { return id; }
    String getEventId() { return eventId; }
    String getEventType() { return eventType; }
    String getRoutingKey() { return routingKey; }
    String getPayload() { return payload; }
    Instant getOccurredAt() { return occurredAt; }
    int getAttempts() { return attempts; }

    void markPublished(Instant when) {
        this.publishedAt = when;
        this.lastError = null;
    }

    void markFailed(String error) {
        this.attempts += 1;
        this.lastError = truncate(error);
    }

    private static String truncate(String s) {
        if (s == null) return null;
        return s.length() > 1000 ? s.substring(0, 1000) : s;
    }
}
