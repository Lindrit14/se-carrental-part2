package com.uni.carbooking.infrastructure.inbox;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import java.time.Instant;
import java.util.UUID;

@Entity
@Table(name = "event_log")
class EventLogJpaEntity {

    @Id
    @Column(name = "id")
    private UUID id;

    @Column(name = "event_id", nullable = false, unique = true)
    private String eventId;

    @Column(name = "event_type", nullable = false)
    private String eventType;

    @Column(name = "source", nullable = false)
    private String source;

    @JdbcTypeCode(SqlTypes.JSON)
    @Column(name = "payload", nullable = false, columnDefinition = "jsonb")
    private String payload;

    @Column(name = "received_at", nullable = false)
    private Instant receivedAt;

    @Column(name = "processed_at")
    private Instant processedAt;

    protected EventLogJpaEntity() {}

    static EventLogJpaEntity create(String eventId, String eventType, String source,
                                    String payload, Instant receivedAt) {
        var e = new EventLogJpaEntity();
        e.id = UUID.randomUUID();
        e.eventId = eventId;
        e.eventType = eventType;
        e.source = source;
        e.payload = payload;
        e.receivedAt = receivedAt;
        return e;
    }

    void markProcessed(Instant when) { this.processedAt = when; }
    String getEventId() { return eventId; }
}
