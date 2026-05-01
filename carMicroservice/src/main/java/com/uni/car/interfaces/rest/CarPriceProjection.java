package com.uni.car.interfaces.rest;

import com.uni.car.application.port.out.CurrencyConverterPort;
import com.uni.car.application.port.out.CurrencyConverterPort.RateSet;
import com.uni.car.domain.car.Car;
import com.uni.car.domain.error.CurrencyServiceUnavailable;
import com.uni.car.interfaces.rest.dto.CarResponse;
import io.grpc.StatusRuntimeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.math.BigDecimal;
import java.math.RoundingMode;

/**
 * Decorates a {@link Car} into a {@link CarResponse}, applying an optional
 * target-currency conversion using rates fetched from currency-converter.
 *
 * <p>Math mirrors the Python {@code ConvertUseCase}: EUR-pivot, 10-digit
 * intermediate, 4-digit HALF_UP result. If the rates RPC fails, the projection
 * still returns cars (without converted fields) so the listing degrades
 * gracefully rather than 5xx-ing.
 */
@Component
public class CarPriceProjection {

    private static final Logger log = LoggerFactory.getLogger(CarPriceProjection.class);
    private static final String EUR = "EUR";

    private final CurrencyConverterPort currencyConverter;

    CarPriceProjection(CurrencyConverterPort currencyConverter) {
        this.currencyConverter = currencyConverter;
    }

    Renderer rendererFor(String targetCurrency) {
        var normalized = normalize(targetCurrency);
        if (normalized == null) {
            return car -> CarResponse.from(car);
        }
        RateSet rateSet;
        try {
            rateSet = currencyConverter.getRates();
        } catch (CurrencyServiceUnavailable | StatusRuntimeException e) {
            log.warn("currency-converter unavailable, returning source-currency prices only: {}",
                    e.getMessage());
            return car -> CarResponse.from(car);
        }
        return car -> CarResponse.fromWithConversion(
                car,
                convert(car.dailyRate().amount(), car.dailyRate().currency(), normalized, rateSet),
                normalized);
    }

    @FunctionalInterface
    interface Renderer {
        CarResponse render(Car car);
    }

    private static String normalize(String raw) {
        if (raw == null) return null;
        var trimmed = raw.trim();
        if (trimmed.isEmpty() || trimmed.length() != 3) return null;
        return trimmed.toUpperCase();
    }

    private static BigDecimal convert(BigDecimal amount, String from, String to, RateSet rateSet) {
        if (from.equalsIgnoreCase(to)) return null;
        var rateFrom = rateFor(from, rateSet);
        var rateTo = rateFor(to, rateSet);
        if (rateFrom == null || rateTo == null) return null;
        return amount
                .divide(rateFrom, 10, RoundingMode.HALF_UP)
                .multiply(rateTo)
                .setScale(4, RoundingMode.HALF_UP);
    }

    private static BigDecimal rateFor(String currency, RateSet rateSet) {
        if (EUR.equalsIgnoreCase(currency)) return BigDecimal.ONE;
        return rateSet.rates().get(currency.toUpperCase());
    }
}
