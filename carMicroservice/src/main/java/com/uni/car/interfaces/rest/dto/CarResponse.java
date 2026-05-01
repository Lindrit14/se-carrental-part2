package com.uni.car.interfaces.rest.dto;

import com.uni.car.domain.car.Car;
import com.uni.car.domain.car.CarCategory;

import java.math.BigDecimal;

public record CarResponse(
        String id,
        String brand,
        String model,
        String licensePlate,
        BigDecimal dailyRateAmount,
        String dailyRateCurrency,
        BigDecimal dailyRateConvertedAmount,
        String dailyRateConvertedCurrency,
        String location,
        CarCategory category
) {
    public static CarResponse from(Car c) {
        return new CarResponse(
                c.id(), c.brand(), c.model(), c.licensePlate(),
                c.dailyRate().amount(), c.dailyRate().currency(),
                null, null,
                c.location(), c.category()
        );
    }

    public static CarResponse fromWithConversion(Car c,
                                                 BigDecimal convertedAmount,
                                                 String convertedCurrency) {
        return new CarResponse(
                c.id(), c.brand(), c.model(), c.licensePlate(),
                c.dailyRate().amount(), c.dailyRate().currency(),
                convertedAmount, convertedAmount == null ? null : convertedCurrency,
                c.location(), c.category()
        );
    }
}
