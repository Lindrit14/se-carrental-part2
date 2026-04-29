package com.uni.carbooking.infrastructure.messaging.rabbitmq;

import com.fasterxml.jackson.databind.JsonNode;
import com.uni.carbooking.application.port.out.CurrencyConverterPort;
import com.uni.carbooking.domain.error.CurrencyConversionFailed;
import com.uni.carbooking.domain.money.Money;
import org.springframework.amqp.AmqpException;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.math.BigDecimal;
import java.util.Map;

/**
 * RPC client for the currency-converter service. Sends a JSON request to
 * {@code currency.requests} via the default exchange and waits (bounded by
 * RabbitTemplate.replyTimeout) for a JSON reply on the Direct-Reply-To channel.
 */
@Component
class RabbitCurrencyConverterClient implements CurrencyConverterPort {

    private final RabbitTemplate template;
    private final String queue;

    RabbitCurrencyConverterClient(
            RabbitTemplate template,
            @Value("${currency.rpc.queue:currency.requests}") String queue) {
        this.template = template;
        this.queue = queue;
    }

    @Override
    public Money convert(Money source, String targetCurrency) {
        var request = Map.of(
                "amount", source.amount().toPlainString(),
                "from", source.currency(),
                "to", targetCurrency
        );
        JsonNode reply;
        try {
            reply = (JsonNode) template.convertSendAndReceiveAsType(
                    "", queue, request, new org.springframework.core.ParameterizedTypeReference<JsonNode>() {}
            );
        } catch (AmqpException e) {
            throw new CurrencyConversionFailed("currency-converter unavailable: " + e.getMessage());
        }
        if (reply == null) {
            throw new CurrencyConversionFailed("currency-converter timed out");
        }
        if (reply.has("error")) {
            throw new CurrencyConversionFailed("currency-converter error: " + reply.get("error").asText());
        }
        return new Money(
                new BigDecimal(reply.get("amount").asText()),
                reply.get("currency").asText()
        );
    }
}
