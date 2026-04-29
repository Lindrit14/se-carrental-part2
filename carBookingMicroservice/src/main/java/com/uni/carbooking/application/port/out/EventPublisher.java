package com.uni.carbooking.application.port.out;

/**
 * Outbound port for publishing domain events. Adapters serialize and route.
 *
 * @param routingKey topic-exchange routing key (e.g. "booking.created")
 * @param payload    JSON-serializable payload (records work well)
 */
public interface EventPublisher {
    void publish(String routingKey, String eventType, Object payload);
}
