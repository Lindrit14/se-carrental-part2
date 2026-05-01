package com.uni.carbooking.interfaces.rest;

import com.uni.carbooking.application.port.out.CurrencyConverterPort;
import com.uni.carbooking.application.port.out.CurrencyConverterPort.RateSet;
import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.booking.BookingStatus;
import com.uni.carbooking.domain.car.CarSnapshot;
import com.uni.carbooking.domain.error.CurrencyConversionFailed;
import com.uni.carbooking.domain.money.Money;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.time.Instant;
import java.time.LocalDate;
import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;

class BookingPriceProjectionTest {

    private static Booking booking(Money source, Money lockedTarget) {
        var car = new CarSnapshot("car-1", "VW", "Golf", "AB-1234",
                new Money(new BigDecimal("45.00"), "EUR"));
        return new Booking(
                "b-1", "cust-1", car,
                LocalDate.parse("2026-05-10"), LocalDate.parse("2026-05-13"),
                BookingStatus.CONFIRMED, source, lockedTarget,
                Instant.parse("2026-04-30T10:00:00Z"), Instant.parse("2026-04-30T10:00:00Z"));
    }

    private static RateSet rates() {
        return new RateSet("EUR", LocalDate.parse("2026-04-28"),
                Map.of("USD", new BigDecimal("1.0823"),
                       "GBP", new BigDecimal("0.8634")));
    }

    @Test
    void no_target_currency_uses_locked_target() {
        var b = booking(
                new Money(new BigDecimal("135.00"), "EUR"),
                new Money(new BigDecimal("146.11"), "USD"));
        var projection = new BookingPriceProjection(stubRates(rates()));

        var total = projection.displayTotalFor(null).totalFor(b);
        assertThat(total.currency()).isEqualTo("USD");
        assertThat(total.amount()).isEqualByComparingTo("146.11");
    }

    @Test
    void target_matches_locked_target_skips_rpc() {
        var b = booking(
                new Money(new BigDecimal("135.00"), "EUR"),
                new Money(new BigDecimal("146.11"), "USD"));
        var port = new ThrowingPort();  // would fail if getRates was called
        var projection = new BookingPriceProjection(port);

        var total = projection.displayTotalFor("USD").totalFor(b);
        assertThat(total.currency()).isEqualTo("USD");
        assertThat(total.amount()).isEqualByComparingTo("146.11");
        assertThat(port.getRatesCalls).isZero();
    }

    @Test
    void target_matches_source_returns_source_without_rpc() {
        var b = booking(
                new Money(new BigDecimal("135.00"), "EUR"),
                new Money(new BigDecimal("146.11"), "USD"));
        var port = new ThrowingPort();
        var projection = new BookingPriceProjection(port);

        var total = projection.displayTotalFor("EUR").totalFor(b);
        assertThat(total.currency()).isEqualTo("EUR");
        assertThat(total.amount()).isEqualByComparingTo("135.00");
        assertThat(port.getRatesCalls).isZero();
    }

    @Test
    void reconverts_through_eur_pivot_when_target_differs() {
        var b = booking(
                new Money(new BigDecimal("135.00"), "EUR"),
                new Money(new BigDecimal("146.11"), "USD"));
        var projection = new BookingPriceProjection(stubRates(rates()));

        // 135 EUR / 1 * 0.8634 = 116.5590
        var total = projection.displayTotalFor("GBP").totalFor(b);
        assertThat(total.currency()).isEqualTo("GBP");
        assertThat(total.amount()).isEqualByComparingTo("116.5590");
    }

    @Test
    void unavailable_currency_service_falls_back_to_locked_target() {
        var b = booking(
                new Money(new BigDecimal("135.00"), "EUR"),
                new Money(new BigDecimal("146.11"), "USD"));
        var projection = new BookingPriceProjection(new ThrowingPort());

        var total = projection.displayTotalFor("GBP").totalFor(b);
        assertThat(total.currency()).isEqualTo("USD");
        assertThat(total.amount()).isEqualByComparingTo("146.11");
    }

    @Test
    void caches_rates_across_bookings_in_same_projection() {
        var b1 = booking(
                new Money(new BigDecimal("100.00"), "EUR"),
                new Money(new BigDecimal("108.23"), "USD"));
        var b2 = booking(
                new Money(new BigDecimal("200.00"), "EUR"),
                new Money(new BigDecimal("216.46"), "USD"));
        var port = new CountingPort(rates());
        var projection = new BookingPriceProjection(port);

        var displayTotal = projection.displayTotalFor("GBP");
        displayTotal.totalFor(b1);
        displayTotal.totalFor(b2);
        assertThat(port.getRatesCalls).isEqualTo(1);
    }

    private static CurrencyConverterPort stubRates(RateSet rs) {
        return new CurrencyConverterPort() {
            @Override
            public Money convert(Money source, String targetCurrency) {
                throw new UnsupportedOperationException();
            }

            @Override
            public RateSet getRates() {
                return rs;
            }
        };
    }

    private static final class ThrowingPort implements CurrencyConverterPort {
        int getRatesCalls;

        @Override
        public Money convert(Money source, String targetCurrency) {
            throw new UnsupportedOperationException();
        }

        @Override
        public RateSet getRates() {
            getRatesCalls++;
            throw new CurrencyConversionFailed("unavailable");
        }
    }

    private static final class CountingPort implements CurrencyConverterPort {
        int getRatesCalls;
        private final RateSet rs;

        CountingPort(RateSet rs) {
            this.rs = rs;
        }

        @Override
        public Money convert(Money source, String targetCurrency) {
            throw new UnsupportedOperationException();
        }

        @Override
        public RateSet getRates() {
            getRatesCalls++;
            return rs;
        }
    }
}
