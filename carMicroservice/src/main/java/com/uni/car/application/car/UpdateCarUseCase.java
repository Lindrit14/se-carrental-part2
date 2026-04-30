package com.uni.car.application.car;

import com.uni.car.domain.car.Car;
import com.uni.car.domain.car.CarCategory;
import com.uni.car.domain.car.CarRepository;
import com.uni.car.domain.error.CarNotFound;
import com.uni.car.domain.money.Money;

public class UpdateCarUseCase {

    private final CarRepository cars;

    public UpdateCarUseCase(CarRepository cars) {
        this.cars = cars;
    }

    public record Input(
            String id,
            String brand,
            String model,
            String licensePlate,
            Money dailyRate,
            String location,
            CarCategory category
    ) {}

    public Car execute(Input in) {
        Car car = cars.findById(in.id()).orElseThrow(() -> new CarNotFound(in.id()));
        car.update(in.brand(), in.model(), in.licensePlate(), in.dailyRate(), in.location(), in.category());
        cars.save(car);
        return car;
    }
}
