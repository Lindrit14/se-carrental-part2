package com.uni.car.application.car;

import com.uni.car.domain.car.Car;
import com.uni.car.domain.car.CarRepository;

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
