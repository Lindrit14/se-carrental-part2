package com.uni.carbooking.infrastructure.outbox;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.uni.carbooking.application.event.BookingEvents;
import com.uni.carbooking.application.port.out.EventPublisher;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.UUID;

/**
 * EventPublisher implementation that writes events to the outbox table inside
 * the caller's transaction. A separate scheduled component ({@link OutboxRelay})
 * drains the table to RabbitMQ asynchronously.
 *
 * Uses MANDATORY propagation: every event must be written from within an
 * existing transaction. If a use case forgets @Transactional, this fails fast
 * rather than silently writing the event without the domain change.
 */
@Component
public class OutboxEventPublisher implements EventPublisher {

    private final OutboxJpaRepository outbox;
    private final ObjectMapper json;

    OutboxEventPublisher(OutboxJpaRepository outbox, ObjectMapper json) {
        this.outbox = outbox;
        this.json = json;
    }

    @Override
    @Transactional(propagation = Propagation.MANDATORY)
    public void publish(String routingKey, String eventType, Object payload) {
        String aggregateId = extractAggregateId(payload);
        String eventId = UUID.randomUUID().toString();
        String body = serialize(payload);

        outbox.save(OutboxJpaEntity.create(eventId, aggregateId, eventType, routingKey, body, Instant.now()));
    }

    private String serialize(Object payload) {
        try {
            return json.writeValueAsString(payload);
        } catch (JsonProcessingException e) {
            throw new IllegalStateException("failed to serialize outbox payload", e);
        }
    }

    private static String extractAggregateId(Object payload) {
        if (payload instanceof BookingEvents.BookingCreated bc) return bc.bookingId();
        if (payload instanceof BookingEvents.BookingCancelled bc) return bc.bookingId();
        return "unknown";
    }
}
