package com.uni.carbooking.interfaces.rest;

import com.uni.carbooking.application.port.out.CurrencyConverterPort;
import com.uni.carbooking.application.port.out.CurrencyConverterPort.RateSet;
import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.error.CurrencyConversionFailed;
import com.uni.carbooking.domain.money.Money;
import io.grpc.StatusRuntimeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.math.BigDecimal;
import java.math.RoundingMode;

/**
 * Picks the right "display" total for a booking on listing endpoints.
 *
 * <p>The booking's {@code totalSource} is the authoritative charge in the
 * car's source currency, and {@code totalTarget} is the conversion locked at
 * booking creation. For listing, we want the total rendered in the user's
 * <em>currently selected</em> display currency (sent as {@code targetCurrency}
 * on the request) — which usually differs from the locked target if the user
 * has switched currency since booking.
 *
 * <p>Strategy:
 * <ul>
 *   <li>No {@code targetCurrency} on the request: keep the stored locked target.</li>
 *   <li>Target equals locked target: keep the stored locked target (no RPC).</li>
 *   <li>Target equals source currency: render the source itself, no conversion.</li>
 *   <li>Otherwise: fetch ECB rates once, EUR-pivot math locally per booking.</li>
 *   <li>Currency-converter unavailable: fall back to the stored locked target so
 *       the bookings page still renders.</li>
 * </ul>
 */
@Component
public class BookingPriceProjection {

    private static final Logger log = LoggerFactory.getLogger(BookingPriceProjection.class);
    private static final String EUR = "EUR";

    private final CurrencyConverterPort currencyConverter;

    BookingPriceProjection(CurrencyConverterPort currencyConverter) {
        this.currencyConverter = currencyConverter;
    }

    public DisplayTotal displayTotalFor(String requestedTargetCurrency) {
        var normalized = normalize(requestedTargetCurrency);
        if (normalized == null) {
            return Booking::totalTarget;
        }

        // Lazy-load: only fetch rates when at least one booking actually needs
        // a reconversion (e.g. all stored locked targets already match).
        return new DisplayTotal() {
            private RateSet rates;
            private boolean ratesUnavailable;

            @Override
            public Money totalFor(Booking booking) {
                if (normalized.equalsIgnoreCase(booking.totalTarget().currency())) {
                    return booking.totalTarget();
                }
                if (normalized.equalsIgnoreCase(booking.totalSource().currency())) {
                    return booking.totalSource();
                }
                if (ratesUnavailable) {
                    return booking.totalTarget();
                }
                if (rates == null) {
                    try {
                        rates = currencyConverter.getRates();
                    } catch (CurrencyConversionFailed | StatusRuntimeException e) {
                        log.warn("currency-converter unavailable, falling back to stored locked target: {}",
                                e.getMessage());
                        ratesUnavailable = true;
                        return booking.totalTarget();
                    }
                }
                var converted = convert(booking.totalSource(), normalized, rates);
                return converted != null ? converted : booking.totalTarget();
            }
        };
    }

    @FunctionalInterface
    public interface DisplayTotal {
        Money totalFor(Booking booking);
    }

    private static String normalize(String raw) {
        if (raw == null) return null;
        var trimmed = raw.trim();
        if (trimmed.length() != 3) return null;
        return trimmed.toUpperCase();
    }

    private static Money convert(Money source, String target, RateSet rateSet) {
        var rateFrom = rateFor(source.currency(), rateSet);
        var rateTo = rateFor(target, rateSet);
        if (rateFrom == null || rateTo == null) return null;
        var amount = source.amount()
                .divide(rateFrom, 10, RoundingMode.HALF_UP)
                .multiply(rateTo)
                .setScale(4, RoundingMode.HALF_UP);
        return new Money(amount, target);
    }

    private static BigDecimal rateFor(String currency, RateSet rateSet) {
        if (EUR.equalsIgnoreCase(currency)) return BigDecimal.ONE;
        return rateSet.rates().get(currency.toUpperCase());
    }
}
