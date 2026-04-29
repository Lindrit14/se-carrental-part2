package com.uni.carbooking.domain.booking;

import com.uni.carbooking.domain.error.InvalidDateRange;
import com.uni.carbooking.domain.money.Money;

import java.time.Instant;
import java.time.LocalDate;
import java.time.temporal.ChronoUnit;
import java.util.Objects;

public class Booking {
    private final String id;
    private final String customerId;
    private final String carId;
    private final LocalDate startDate;
    private final LocalDate endDate;
    private BookingStatus status;
    private final Money totalSource;
    private final Money totalTarget;
    private final Instant createdAt;
    private Instant updatedAt;

    public Booking(String id, String customerId, String carId,
                   LocalDate startDate, LocalDate endDate, BookingStatus status,
                   Money totalSource, Money totalTarget,
                   Instant createdAt, Instant updatedAt) {
        if (!endDate.isAfter(startDate)) {
            throw new InvalidDateRange("end_date must be after start_date");
        }
        this.id = Objects.requireNonNull(id);
        this.customerId = Objects.requireNonNull(customerId);
        this.carId = Objects.requireNonNull(carId);
        this.startDate = startDate;
        this.endDate = endDate;
        this.status = Objects.requireNonNull(status);
        this.totalSource = Objects.requireNonNull(totalSource);
        this.totalTarget = Objects.requireNonNull(totalTarget);
        this.createdAt = Objects.requireNonNull(createdAt);
        this.updatedAt = Objects.requireNonNull(updatedAt);
    }

    public String id() { return id; }
    public String customerId() { return customerId; }
    public String carId() { return carId; }
    public LocalDate startDate() { return startDate; }
    public LocalDate endDate() { return endDate; }
    public BookingStatus status() { return status; }
    public Money totalSource() { return totalSource; }
    public Money totalTarget() { return totalTarget; }
    public Instant createdAt() { return createdAt; }
    public Instant updatedAt() { return updatedAt; }

    public long days() {
        return ChronoUnit.DAYS.between(startDate, endDate);
    }

    public void cancel(Instant now) {
        this.status = BookingStatus.CANCELLED;
        this.updatedAt = now;
    }
}
