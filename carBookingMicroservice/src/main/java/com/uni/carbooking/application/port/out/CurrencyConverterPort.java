package com.uni.carbooking.application.port.out;

import com.uni.carbooking.domain.money.Money;

/**
 * Outbound port for currency conversion. Adapter implements this via
 * RabbitMQ-RPC against the currency-converter service.
 */
public interface CurrencyConverterPort {
    Money convert(Money source, String targetCurrency);
}
