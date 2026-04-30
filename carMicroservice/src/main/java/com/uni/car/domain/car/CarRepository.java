package com.uni.car.domain.car;

import java.util.List;
import java.util.Optional;

public interface CarRepository {
    void save(Car car);
    Optional<Car> findById(String id);
    List<Car> findAll();
    void deleteById(String id);

    /**
     * Filter by location (case-insensitive substring) and/or category.
     * Null arguments mean "any". Date-range availability is not checked here —
     * that is the booking service's responsibility.
     */
    List<Car> search(String location, CarCategory category);
}
