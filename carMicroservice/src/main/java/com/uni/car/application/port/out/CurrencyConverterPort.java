package com.uni.car.application.port.out;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.util.Map;

public interface CurrencyConverterPort {

    RateSet getRates();

    record RateSet(String base, LocalDate rateDate, Map<String, BigDecimal> rates) {}
}
