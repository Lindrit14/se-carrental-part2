package com.uni.carbooking.application.port.out;

import com.uni.carbooking.domain.money.Money;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.util.Map;

/** Outbound port for the currency-converter service (gRPC adapter). */
public interface CurrencyConverterPort {

    /** Single conversion. Used at booking creation to lock in a target total. */
    Money convert(Money source, String targetCurrency);

    /** Snapshot of EUR-based reference rates. Used at list time to convert
     *  many bookings into the user's currently selected display currency
     *  without making one RPC per booking. */
    RateSet getRates();

    record RateSet(String base, LocalDate rateDate, Map<String, BigDecimal> rates) {}
}
