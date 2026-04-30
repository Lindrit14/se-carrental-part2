package com.uni.car.application.car;

import com.uni.car.application.port.out.Clock;
import com.uni.car.application.port.out.IdGenerator;
import com.uni.car.domain.car.Car;
import com.uni.car.domain.car.CarCategory;
import com.uni.car.domain.car.CarRepository;
import com.uni.car.domain.money.Money;

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
