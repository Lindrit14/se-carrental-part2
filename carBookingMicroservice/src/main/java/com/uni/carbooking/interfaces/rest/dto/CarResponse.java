package com.uni.carbooking.interfaces.rest.dto;

import com.uni.carbooking.domain.car.Car;

import java.math.BigDecimal;

public record CarResponse(
        String id,
        String brand,
        String model,
        String licensePlate,
        BigDecimal dailyRateAmount,
        String dailyRateCurrency
) {
    public static CarResponse from(Car c) {
        return new CarResponse(c.id(), c.brand(), c.model(), c.licensePlate(),
                c.dailyRate().amount(), c.dailyRate().currency());
    }
}
