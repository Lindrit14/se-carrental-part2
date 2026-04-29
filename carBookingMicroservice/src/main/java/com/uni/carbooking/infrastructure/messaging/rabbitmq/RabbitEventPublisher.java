package com.uni.carbooking.infrastructure.messaging.rabbitmq;

import com.uni.carbooking.application.event.BookingEvents;
import com.uni.carbooking.application.port.out.EventPublisher;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;

import java.time.Instant;
import java.util.UUID;

@Component
class RabbitEventPublisher implements EventPublisher {

    private static final String EXCHANGE = "booking.events";

    private final RabbitTemplate template;

    RabbitEventPublisher(RabbitTemplate template) {
        this.template = template;
    }

    @Override
    public void publish(String routingKey, String eventType, Object payload) {
        var envelope = new EventEnvelope(
                UUID.randomUUID().toString(),
                eventType,
                BookingEvents.VERSION,
                Instant.now().toString(),
                payload
        );
        template.convertAndSend(EXCHANGE, routingKey, envelope);
    }
}
