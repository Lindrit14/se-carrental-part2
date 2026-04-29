package com.uni.carbooking.application.car;

import com.uni.carbooking.domain.car.Car;
import com.uni.carbooking.domain.car.CarCategory;
import com.uni.carbooking.domain.car.CarRepository;

import java.time.LocalDate;
import java.util.List;

public class SearchCarsUseCase {

    private final CarRepository cars;

    public SearchCarsUseCase(CarRepository cars) {
        this.cars = cars;
    }

    public record Input(String location, CarCategory category, LocalDate from, LocalDate to) {}

    public List<Car> execute(Input in) {
        if (in.from() != null && in.to() != null && !in.to().isAfter(in.from())) {
            throw new IllegalArgumentException("'to' must be after 'from'");
        }
        var loc = (in.location() == null || in.location().isBlank()) ? null : in.location().trim();
        return cars.search(loc, in.category(), in.from(), in.to());
    }
}
