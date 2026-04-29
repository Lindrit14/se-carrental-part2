package com.uni.carbooking.domain.car;

import java.util.List;
import java.util.Optional;

public interface CarRepository {
    void save(Car car);
    Optional<Car> findById(String id);
    List<Car> findAll();
}
