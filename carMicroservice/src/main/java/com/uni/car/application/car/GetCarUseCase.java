package com.uni.car.application.car;

import com.uni.car.domain.car.Car;
import com.uni.car.domain.car.CarRepository;
import com.uni.car.domain.error.CarNotFound;

public class GetCarUseCase {

    private final CarRepository cars;

    public GetCarUseCase(CarRepository cars) {
        this.cars = cars;
    }

    public Car execute(String id) {
        return cars.findById(id).orElseThrow(() -> new CarNotFound(id));
    }
}
