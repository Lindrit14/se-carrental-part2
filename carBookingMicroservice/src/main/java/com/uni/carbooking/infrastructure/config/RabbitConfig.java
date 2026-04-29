package com.uni.carbooking.infrastructure.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.amqp.core.AcknowledgeMode;
import org.springframework.amqp.core.Binding;
import org.springframework.amqp.core.BindingBuilder;
import org.springframework.amqp.core.Queue;
import org.springframework.amqp.core.QueueBuilder;
import org.springframework.amqp.core.TopicExchange;
import org.springframework.amqp.rabbit.config.SimpleRabbitListenerContainerFactory;
import org.springframework.amqp.rabbit.connection.ConnectionFactory;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.amqp.support.converter.Jackson2JsonMessageConverter;
import org.springframework.amqp.support.converter.MessageConverter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * Defines exchanges, queues and bindings used by booking, plus the
 * RabbitTemplate (for publishing + RPC) and a JSON message converter.
 *
 * Topology
 *   user.events     (topic, durable)  — owned by user-auth, we declare it again idempotently
 *   booking.events  (topic, durable)  — owned by booking
 *   booking.user.events       (queue) — bound to user.events with "user.*"
 *   currency.requests         (queue) — RPC target, currency-converter consumes
 */
@Configuration
class RabbitConfig {

    static final String USER_EVENTS_EXCHANGE = "user.events";
    static final String BOOKING_EVENTS_EXCHANGE = "booking.events";
    static final String USER_EVENTS_QUEUE = "booking.user.events";

    @Value("${currency.rpc.queue:currency.requests}")
    private String currencyRpcQueue;

    @Value("${currency.rpc.timeout-ms:5000}")
    private long currencyRpcTimeoutMs;

    @Bean TopicExchange userEventsExchange() {
        return new TopicExchange(USER_EVENTS_EXCHANGE, true, false);
    }

    @Bean TopicExchange bookingEventsExchange() {
        return new TopicExchange(BOOKING_EVENTS_EXCHANGE, true, false);
    }

    @Bean Queue userEventsQueue() {
        return QueueBuilder.durable(USER_EVENTS_QUEUE).build();
    }

    @Bean Binding userEventsBinding(Queue userEventsQueue, TopicExchange userEventsExchange) {
        return BindingBuilder.bind(userEventsQueue).to(userEventsExchange).with("user.*");
    }

    @Bean Queue currencyRequestsQueue() {
        return QueueBuilder.durable(currencyRpcQueue).build();
    }

    /**
     * Reuses Spring Boot's auto-configured ObjectMapper — important because
     * the auto-configured one has JavaTimeModule registered (so LocalDate
     * etc. serialize as ISO strings, not "unsupported"). Defining our own
     * @Bean ObjectMapper would replace it and silently drop JSR-310 support.
     * Jackson defaults are tweaked via spring.jackson.* in application.yml.
     */
    @Bean MessageConverter jsonMessageConverter(ObjectMapper objectMapper) {
        return new Jackson2JsonMessageConverter(objectMapper);
    }

    @Bean
    RabbitTemplate rabbitTemplate(ConnectionFactory cf, MessageConverter converter) {
        var template = new RabbitTemplate(cf);
        template.setMessageConverter(converter);
        template.setReplyTimeout(currencyRpcTimeoutMs);
        // Use Direct-Reply-To: no temp queues, lower latency.
        template.setUseDirectReplyToContainer(true);
        return template;
    }

    @Bean
    SimpleRabbitListenerContainerFactory rabbitListenerContainerFactory(
            ConnectionFactory cf, MessageConverter converter) {
        var factory = new SimpleRabbitListenerContainerFactory();
        factory.setConnectionFactory(cf);
        factory.setMessageConverter(converter);
        factory.setAcknowledgeMode(AcknowledgeMode.AUTO);
        factory.setDefaultRequeueRejected(false);
        return factory;
    }
}
