package com.uni.car.application.car;

import com.uni.car.domain.car.CarRepository;
import com.uni.car.domain.error.CarNotFound;

public class DeleteCarUseCase {

    private final CarRepository cars;

    public DeleteCarUseCase(CarRepository cars) {
        this.cars = cars;
    }

    public void execute(String id) {
        if (cars.findById(id).isEmpty()) {
            throw new CarNotFound(id);
        }
        cars.deleteById(id);
    }
}
