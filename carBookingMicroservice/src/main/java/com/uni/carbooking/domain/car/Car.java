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
    private final String location;
    private final CarCategory category;
    private final Instant createdAt;

    public Car(String id, String brand, String model, String licensePlate,
               Money dailyRate, String location, CarCategory category, Instant createdAt) {
        this.id = Objects.requireNonNull(id);
        this.brand = Objects.requireNonNull(brand);
        this.model = Objects.requireNonNull(model);
        this.licensePlate = Objects.requireNonNull(licensePlate);
        this.dailyRate = Objects.requireNonNull(dailyRate);
        this.location = Objects.requireNonNull(location);
        this.category = Objects.requireNonNull(category);
        this.createdAt = Objects.requireNonNull(createdAt);
    }

    public String id() { return id; }
    public String brand() { return brand; }
    public String model() { return model; }
    public String licensePlate() { return licensePlate; }
    public Money dailyRate() { return dailyRate; }
    public String location() { return location; }
    public CarCategory category() { return category; }
    public Instant createdAt() { return createdAt; }
}
