package com.uni.carbooking.infrastructure.messaging.rabbitmq;

/**
 * Wire-format envelope shared by all events the platform publishes/consumes.
 * Mirrors the shape used by user-auth (Go) so consumers across services agree.
 */
record EventEnvelope(
        String event_id,
        String event_type,
        String version,
        String occurred_at,
        Object data
) {
}
