package com.uni.carbooking.infrastructure.messaging.rabbitmq;

import com.fasterxml.jackson.databind.JsonNode;
import com.uni.carbooking.application.customer.AnonymizeCustomerUseCase;
import com.uni.carbooking.application.customer.RegisterCustomerUseCase;
import com.uni.carbooking.application.customer.UpdateCustomerEmailUseCase;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;

/**
 * Consumes user.* events from the user-auth service.
 * Queue + binding declared in {@link com.uni.carbooking.infrastructure.config.RabbitConfig}.
 */
@Component
class UserEventListener {

    private static final Logger log = LoggerFactory.getLogger(UserEventListener.class);

    private final RegisterCustomerUseCase register;
    private final UpdateCustomerEmailUseCase updateEmail;
    private final AnonymizeCustomerUseCase anonymize;

    UserEventListener(RegisterCustomerUseCase register,
                      UpdateCustomerEmailUseCase updateEmail,
                      AnonymizeCustomerUseCase anonymize) {
        this.register = register;
        this.updateEmail = updateEmail;
        this.anonymize = anonymize;
    }

    @RabbitListener(queues = "booking.user.events")
    public void handle(JsonNode envelope) {
        String type = textField(envelope, "event_type");
        JsonNode data = envelope.path("data");
        if (type == null || data.isMissingNode()) {
            log.warn("user_event_malformed: {}", envelope);
            return;
        }

        switch (type) {
            case "user.registered" -> register.execute(
                    textField(data, "user_id"),
                    textField(data, "email")
            );
            case "user.updated" -> {
                JsonNode changed = data.path("changed_fields");
                boolean emailChanged = changed.isArray() &&
                        java.util.stream.StreamSupport.stream(changed.spliterator(), false)
                                .anyMatch(n -> "email".equals(n.asText()));
                if (emailChanged) {
                    updateEmail.execute(textField(data, "user_id"), textField(data, "email"));
                }
            }
            case "user.deleted" -> anonymize.execute(textField(data, "user_id"));
            default -> log.debug("user_event_ignored: {}", type);
        }
    }

    private static String textField(JsonNode n, String key) {
        JsonNode v = n.path(key);
        return v.isMissingNode() || v.isNull() ? null : v.asText();
    }
}
