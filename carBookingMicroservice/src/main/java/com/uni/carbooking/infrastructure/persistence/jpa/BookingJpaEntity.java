package com.uni.carbooking.infrastructure.persistence.jpa;

import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.booking.BookingStatus;
import com.uni.carbooking.domain.car.CarSnapshot;
import com.uni.carbooking.domain.money.Money;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

import java.math.BigDecimal;
import java.time.Instant;
import java.time.LocalDate;

@Entity
@Table(name = "bookings")
class BookingJpaEntity {

    @Id
    @Column(length = 36)
    private String id;

    @Column(name = "customer_id", nullable = false, length = 36)
    private String customerId;

    // Historical reference only — no FK constraint (car lives in a separate service/DB).
    @Column(name = "car_id", nullable = false, length = 36)
    private String carId;

    @Column(name = "car_brand", nullable = false, length = 80)
    private String carBrand;

    @Column(name = "car_model", nullable = false, length = 120)
    private String carModel;

    @Column(name = "car_license_plate", nullable = false, length = 40)
    private String carLicensePlate;

    @Column(name = "daily_rate_snapshot_amount", nullable = false, precision = 19, scale = 4)
    private BigDecimal dailyRateSnapshotAmount;

    @Column(name = "daily_rate_snapshot_currency", nullable = false, length = 3)
    private String dailyRateSnapshotCurrency;

    @Column(name = "start_date", nullable = false)
    private LocalDate startDate;

    @Column(name = "end_date", nullable = false)
    private LocalDate endDate;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 16)
    private BookingStatus status;

    @Column(name = "total_source_amount", nullable = false, precision = 19, scale = 4)
    private BigDecimal totalSourceAmount;

    @Column(name = "total_source_currency", nullable = false, length = 3)
    private String totalSourceCurrency;

    @Column(name = "total_target_amount", nullable = false, precision = 19, scale = 4)
    private BigDecimal totalTargetAmount;

    @Column(name = "total_target_currency", nullable = false, length = 3)
    private String totalTargetCurrency;

    @Column(name = "created_at", nullable = false)
    private Instant createdAt;

    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt;

    protected BookingJpaEntity() {}

    static BookingJpaEntity fromDomain(Booking b) {
        var e = new BookingJpaEntity();
        e.id = b.id();
        e.customerId = b.customerId();
        e.carId = b.carSnapshot().carId();
        e.carBrand = b.carSnapshot().brand();
        e.carModel = b.carSnapshot().model();
        e.carLicensePlate = b.carSnapshot().licensePlate();
        e.dailyRateSnapshotAmount = b.carSnapshot().dailyRate().amount();
        e.dailyRateSnapshotCurrency = b.carSnapshot().dailyRate().currency();
        e.startDate = b.startDate();
        e.endDate = b.endDate();
        e.status = b.status();
        e.totalSourceAmount = b.totalSource().amount();
        e.totalSourceCurrency = b.totalSource().currency();
        e.totalTargetAmount = b.totalTarget().amount();
        e.totalTargetCurrency = b.totalTarget().currency();
        e.createdAt = b.createdAt();
        e.updatedAt = b.updatedAt();
        return e;
    }

    Booking toDomain() {
        var snapshot = new CarSnapshot(
                carId, carBrand, carModel, carLicensePlate,
                new Money(dailyRateSnapshotAmount, dailyRateSnapshotCurrency)
        );
        return new Booking(id, customerId, snapshot, startDate, endDate, status,
                new Money(totalSourceAmount, totalSourceCurrency),
                new Money(totalTargetAmount, totalTargetCurrency),
                createdAt, updatedAt);
    }
}
