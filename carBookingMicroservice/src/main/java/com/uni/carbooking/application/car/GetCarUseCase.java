package com.uni.carbooking.application.car;

import com.uni.carbooking.domain.car.Car;
import com.uni.carbooking.domain.car.CarRepository;
import com.uni.carbooking.domain.error.CarNotFound;

public class GetCarUseCase {

    private final CarRepository cars;

    public GetCarUseCase(CarRepository cars) {
        this.cars = cars;
    }

    public Car execute(String id) {
        return cars.findById(id).orElseThrow(() -> new CarNotFound(id));
    }
}
