package com.uni.car.application.car;

import com.uni.car.domain.car.Car;
import com.uni.car.domain.car.CarCategory;
import com.uni.car.domain.car.CarRepository;

import java.util.List;

public class SearchCarsUseCase {

    private final CarRepository cars;

    public SearchCarsUseCase(CarRepository cars) {
        this.cars = cars;
    }

    public record Input(String location, CarCategory category) {}

    public List<Car> execute(Input in) {
        var loc = (in.location() == null || in.location().isBlank()) ? null : in.location().trim();
        return cars.search(loc, in.category());
    }
}
