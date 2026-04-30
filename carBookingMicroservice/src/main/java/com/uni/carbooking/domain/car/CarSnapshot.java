package com.uni.carbooking.domain.car;

import com.uni.carbooking.domain.money.Money;

import java.util.Objects;

/**
 * Immutable snapshot of a car's details captured at booking creation time.
 * Stored denormalized in the bookings table so booking history is preserved
 * even if the car is later modified or deleted in the car service.
 */
public record CarSnapshot(
        String carId,
        String brand,
        String model,
        String licensePlate,
        Money dailyRate
) {
    public CarSnapshot {
        Objects.requireNonNull(carId, "carId");
        Objects.requireNonNull(brand, "brand");
        Objects.requireNonNull(model, "model");
        Objects.requireNonNull(licensePlate, "licensePlate");
        Objects.requireNonNull(dailyRate, "dailyRate");
    }
}
