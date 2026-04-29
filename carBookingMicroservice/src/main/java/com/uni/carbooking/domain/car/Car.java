package com.uni.carbooking.domain.car;

import com.uni.carbooking.domain.money.Money;

import java.time.Instant;
import java.util.Objects;

/**
 * Aggregate root for a car offered for rental.
 */
public class Car {
    private final String id;
    private String brand;
    private String model;
    private final String licensePlate;
    private Money dailyRate;
    private final Instant createdAt;

    public Car(String id, String brand, String model, String licensePlate,
               Money dailyRate, Instant createdAt) {
        this.id = Objects.requireNonNull(id);
        this.brand = Objects.requireNonNull(brand);
        this.model = Objects.requireNonNull(model);
        this.licensePlate = Objects.requireNonNull(licensePlate);
        this.dailyRate = Objects.requireNonNull(dailyRate);
        this.createdAt = Objects.requireNonNull(createdAt);
    }

    public String id() { return id; }
    public String brand() { return brand; }
    public String model() { return model; }
    public String licensePlate() { return licensePlate; }
    public Money dailyRate() { return dailyRate; }
    public Instant createdAt() { return createdAt; }
}
