package com.uni.car.infrastructure.persistence.jpa;

import com.uni.car.domain.car.Car;
import com.uni.car.domain.car.CarCategory;
import com.uni.car.domain.money.Money;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
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

    @Column(nullable = false, length = 160)
    private String location;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 32)
    private CarCategory category;

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
        e.location = c.location();
        e.category = c.category();
        e.createdAt = c.createdAt();
        return e;
    }

    Car toDomain() {
        return new Car(id, brand, model, licensePlate,
                new Money(dailyRateAmount, dailyRateCurrency), location, category, createdAt);
    }
}
