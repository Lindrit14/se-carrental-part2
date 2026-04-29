package com.uni.carbooking.application.car;

import com.uni.carbooking.application.port.out.Clock;
import com.uni.carbooking.application.port.out.IdGenerator;
import com.uni.carbooking.domain.car.Car;
import com.uni.carbooking.domain.car.CarCategory;
import com.uni.carbooking.domain.car.CarRepository;
import com.uni.carbooking.domain.money.Money;

public class AddCarUseCase {

    private final CarRepository cars;
    private final IdGenerator ids;
    private final Clock clock;

    public AddCarUseCase(CarRepository cars, IdGenerator ids, Clock clock) {
        this.cars = cars;
        this.ids = ids;
        this.clock = clock;
    }

    public record Input(
            String brand,
            String model,
            String licensePlate,
            Money dailyRate,
            String location,
            CarCategory category
    ) {}

    public Car execute(Input in) {
        Car car = new Car(
                ids.newId(),
                in.brand(), in.model(), in.licensePlate(),
                in.dailyRate(),
                in.location(), in.category(),
                clock.now()
        );
        cars.save(car);
        return car;
    }
}
