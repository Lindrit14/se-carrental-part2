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
        String location,
        CarCategory category
) {
    public static CarResponse from(Car c) {
        return new CarResponse(
                c.id(), c.brand(), c.model(), c.licensePlate(),
                c.dailyRate().amount(), c.dailyRate().currency(),
                c.location(), c.category()
        );
    }
}
