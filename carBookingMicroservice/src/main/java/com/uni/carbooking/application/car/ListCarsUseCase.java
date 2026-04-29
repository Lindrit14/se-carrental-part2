package com.uni.carbooking.application.car;

import com.uni.carbooking.domain.car.Car;
import com.uni.carbooking.domain.car.CarRepository;

import java.util.List;

public class ListCarsUseCase {

    private final CarRepository cars;

    public ListCarsUseCase(CarRepository cars) {
        this.cars = cars;
    }

    public List<Car> execute() {
        return cars.findAll();
    }
}
