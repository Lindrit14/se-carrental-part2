package com.uni.carbooking.domain.car;

import java.time.LocalDate;
import java.util.List;
import java.util.Optional;

public interface CarRepository {
    void save(Car car);
    Optional<Car> findById(String id);
    List<Car> findAll();

    /**
     * Search cars matching all provided filters. {@code null} arguments are
     * treated as "any". When both {@code from} and {@code to} are non-null,
     * cars with an overlapping {@code CONFIRMED} booking are excluded.
     */
    List<Car> search(String location, CarCategory category, LocalDate from, LocalDate to);
}
