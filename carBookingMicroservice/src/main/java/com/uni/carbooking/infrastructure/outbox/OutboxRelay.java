package com.uni.carbooking.infrastructure.outbox;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.uni.carbooking.application.event.BookingEvents;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.AmqpException;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Drains the outbox table to RabbitMQ. Runs every {@code outbox.poll-ms} and
 * publishes up to {@code outbox.batch-size} events per tick. Failures are
 * recorded on the row (attempts++, last_error) and retried on the next tick;
 * after {@code outbox.max-attempts} we keep retrying but log at ERROR so an
 * operator can investigate.
 */
@Component
class OutboxRelay {

    private static final Logger log = LoggerFactory.getLogger(OutboxRelay.class);
    private static final String EXCHANGE = "booking.events";

    private final OutboxJpaRepository outbox;
    private final RabbitTemplate rabbit;
    private final ObjectMapper json;

    @Value("${outbox.batch-size:50}")
    private int batchSize;

    @Value("${outbox.max-attempts:5}")
    private int maxAttempts;

    OutboxRelay(OutboxJpaRepository outbox, RabbitTemplate rabbit, ObjectMapper json) {
        this.outbox = outbox;
        this.rabbit = rabbit;
        this.json = json;
    }

    @Scheduled(fixedDelayString = "${outbox.poll-ms:1000}")
    @Transactional
    public void drain() {
        List<OutboxJpaEntity> batch = outbox.findUnpublishedForUpdate(batchSize);
        if (batch.isEmpty()) return;

        for (OutboxJpaEntity event : batch) {
            try {
                rabbit.convertAndSend(EXCHANGE, event.getRoutingKey(), buildEnvelope(event));
                event.markPublished(Instant.now());
            } catch (AmqpException e) {
                event.markFailed(e.getMessage());
                if (event.getAttempts() >= maxAttempts) {
                    log.error("outbox event {} exceeded {} attempts: {}",
                            event.getEventId(), maxAttempts, e.getMessage());
                } else {
                    log.warn("outbox event {} publish failed (attempt {}): {}",
                            event.getEventId(), event.getAttempts(), e.getMessage());
                }
            }
        }
    }

    /**
     * Builds the envelope as a Map (rather than a record) so the relay can
     * include event_id from the outbox row — the same id that consumers use
     * for idempotency, even after retries.
     */
    private Map<String, Object> buildEnvelope(OutboxJpaEntity event) {
        Map<String, Object> env = new LinkedHashMap<>();
        env.put("event_id", event.getEventId());
        env.put("event_type", event.getEventType());
        env.put("version", BookingEvents.VERSION);
        env.put("occurred_at", event.getOccurredAt().toString());
        env.put("data", parseData(event.getPayload()));
        return env;
    }

    private Object parseData(String payload) {
        try {
            return json.readTree(payload);
        } catch (JsonProcessingException e) {
            throw new IllegalStateException("corrupted outbox payload", e);
        }
    }
}
