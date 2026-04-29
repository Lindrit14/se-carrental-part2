package com.uni.carbooking.infrastructure.persistence.jpa;

import com.uni.carbooking.domain.car.Car;
import com.uni.carbooking.domain.money.Money;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

import java.math.BigDecimal;
import java.time.Instant;

@Entity
@Table(name = "cars")
class CarJpaEntity {

    @Id
    @Column(length = 36)
    private String id;

    @Column(nullable = false)
    private String brand;

    @Column(nullable = false)
    private String model;

    @Column(name = "license_plate", nullable = false, unique = true)
    private String licensePlate;

    @Column(name = "daily_rate_amount", nullable = false, precision = 19, scale = 4)
    private BigDecimal dailyRateAmount;

    @Column(name = "daily_rate_currency", nullable = false, length = 3)
    private String dailyRateCurrency;

    @Column(name = "created_at", nullable = false)
    private Instant createdAt;

    protected CarJpaEntity() {}

    static CarJpaEntity fromDomain(Car c) {
        var e = new CarJpaEntity();
        e.id = c.id();
        e.brand = c.brand();
        e.model = c.model();
        e.licensePlate = c.licensePlate();
        e.dailyRateAmount = c.dailyRate().amount();
        e.dailyRateCurrency = c.dailyRate().currency();
        e.createdAt = c.createdAt();
        return e;
    }

    Car toDomain() {
        return new Car(id, brand, model, licensePlate,
                new Money(dailyRateAmount, dailyRateCurrency), createdAt);
    }
}
