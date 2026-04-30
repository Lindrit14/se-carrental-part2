package com.uni.carbooking.infrastructure.inbox;

import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;

/**
 * Inbox-pattern guard around incoming events. Call {@link #recordOrSkip} as
 * the first step of any event handler — if it returns {@code false}, the
 * event_id has already been processed (or is being processed) so the handler
 * should ack and return.
 *
 * The unique constraint on event_log.event_id provides the at-most-once
 * guarantee even when the broker redelivers a message.
 */
@Component
public class EventLog {

    private final EventLogJpaRepository repo;

    EventLog(EventLogJpaRepository repo) {
        this.repo = repo;
    }

    /**
     * Tries to insert a row for {@code eventId}. Returns true if this is the
     * first time we see the event (caller should proceed to process it),
     * false if a row already existed (caller should skip).
     *
     * Runs in a NEW transaction so the row stays committed even if the
     * outer handler's transaction rolls back. That guarantees idempotency:
     * a redelivered event is still recognized as a duplicate even if the
     * original processing crashed mid-way. Re-processing is handled by the
     * separate replay path that targets unprocessed event_log rows.
     */
    @Transactional(propagation = Propagation.REQUIRES_NEW)
    public boolean recordOrSkip(String eventId, String eventType, String source, String payload) {
        if (eventId == null || eventId.isBlank()) {
            // Legacy events without an event_id can't be deduped; let them through.
            return true;
        }
        if (repo.existsByEventId(eventId)) {
            return false;
        }
        try {
            repo.save(EventLogJpaEntity.create(eventId, eventType, source, payload, Instant.now()));
            return true;
        } catch (DataIntegrityViolationException e) {
            // Race: another consumer inserted the same event_id concurrently.
            return false;
        }
    }

    @Transactional(propagation = Propagation.REQUIRES_NEW)
    public void markProcessed(String eventId) {
        repo.findByEventId(eventId).ifPresent(e -> e.markProcessed(Instant.now()));
    }
}
