package com.uni.car.interfaces.rest;

import com.uni.car.application.port.out.CurrencyConverterPort;
import com.uni.car.application.port.out.CurrencyConverterPort.RateSet;
import com.uni.car.domain.car.Car;
import com.uni.car.domain.car.CarCategory;
import com.uni.car.domain.error.CurrencyServiceUnavailable;
import com.uni.car.domain.money.Money;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.time.Instant;
import java.time.LocalDate;
import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;

class CarPriceProjectionTest {

    private static Car eurCar() {
        return new Car(
                "c-1", "VW", "Golf", "AB-123",
                new Money(new BigDecimal("100.00"), "EUR"),
                "Berlin", CarCategory.MEDIUM, Instant.parse("2026-04-01T00:00:00Z"));
    }

    @Test
    void no_target_currency_returns_source_only() {
        var projection = new CarPriceProjection(stub(staticRateSet()));
        var renderer = projection.rendererFor(null);
        var resp = renderer.render(eurCar());

        assertThat(resp.dailyRateAmount()).isEqualByComparingTo("100.00");
        assertThat(resp.dailyRateCurrency()).isEqualTo("EUR");
        assertThat(resp.dailyRateConvertedAmount()).isNull();
        assertThat(resp.dailyRateConvertedCurrency()).isNull();
    }

    @Test
    void converts_eur_source_to_usd_target() {
        var projection = new CarPriceProjection(stub(staticRateSet()));
        var resp = projection.rendererFor("USD").render(eurCar());

        // 100 / 1 * 1.0823 = 108.23, scaled to 4dp HALF_UP
        assertThat(resp.dailyRateConvertedAmount()).isEqualByComparingTo("108.2300");
        assertThat(resp.dailyRateConvertedCurrency()).isEqualTo("USD");
    }

    @Test
    void converts_through_eur_pivot_for_non_eur_source() {
        var car = new Car(
                "c-2", "Ford", "Focus", "ZZ-1",
                new Money(new BigDecimal("100.00"), "USD"),
                "London", CarCategory.MEDIUM, Instant.parse("2026-04-01T00:00:00Z"));
        var projection = new CarPriceProjection(stub(staticRateSet()));
        var resp = projection.rendererFor("GBP").render(car);

        // 100 USD -> EUR: 100 / 1.0823 (10dp HALF_UP) = 92.3957774184
        // 92.3957774184 * 0.8634 (4dp HALF_UP) = 79.7746
        assertThat(resp.dailyRateConvertedAmount()).isEqualByComparingTo("79.7746");
        assertThat(resp.dailyRateConvertedCurrency()).isEqualTo("GBP");
    }

    @Test
    void source_equals_target_skips_conversion() {
        var projection = new CarPriceProjection(stub(staticRateSet()));
        var resp = projection.rendererFor("EUR").render(eurCar());

        assertThat(resp.dailyRateConvertedAmount()).isNull();
        assertThat(resp.dailyRateConvertedCurrency()).isNull();
    }

    @Test
    void unknown_target_currency_returns_null_converted() {
        var projection = new CarPriceProjection(stub(staticRateSet()));
        var resp = projection.rendererFor("XYZ").render(eurCar());

        assertThat(resp.dailyRateConvertedAmount()).isNull();
        assertThat(resp.dailyRateConvertedCurrency()).isNull();
    }

    @Test
    void currency_service_unavailable_degrades_gracefully() {
        CurrencyConverterPort failing = () -> {
            throw new CurrencyServiceUnavailable("boom");
        };
        var projection = new CarPriceProjection(failing);
        var resp = projection.rendererFor("USD").render(eurCar());

        assertThat(resp.dailyRateAmount()).isEqualByComparingTo("100.00");
        assertThat(resp.dailyRateConvertedAmount()).isNull();
        assertThat(resp.dailyRateConvertedCurrency()).isNull();
    }

    private static RateSet staticRateSet() {
        return new RateSet(
                "EUR", LocalDate.of(2026, 4, 28),
                Map.of("USD", new BigDecimal("1.0823"), "GBP", new BigDecimal("0.8634")));
    }

    private static CurrencyConverterPort stub(RateSet rs) {
        return () -> rs;
    }
}
