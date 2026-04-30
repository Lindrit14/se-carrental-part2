package com.uni.car.infrastructure.config;

import com.uni.car.application.car.AddCarUseCase;
import com.uni.car.application.car.DeleteCarUseCase;
import com.uni.car.application.car.GetCarUseCase;
import com.uni.car.application.car.ListCarsUseCase;
import com.uni.car.application.car.SearchCarsUseCase;
import com.uni.car.application.car.UpdateCarUseCase;
import com.uni.car.application.port.out.Clock;
import com.uni.car.application.port.out.IdGenerator;
import com.uni.car.domain.car.CarRepository;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
class UseCaseConfig {

    @Bean ListCarsUseCase listCarsUseCase(CarRepository cars) {
        return new ListCarsUseCase(cars);
    }

    @Bean SearchCarsUseCase searchCarsUseCase(CarRepository cars) {
        return new SearchCarsUseCase(cars);
    }

    @Bean GetCarUseCase getCarUseCase(CarRepository cars) {
        return new GetCarUseCase(cars);
    }

    @Bean AddCarUseCase addCarUseCase(CarRepository cars, IdGenerator ids, Clock clock) {
        return new AddCarUseCase(cars, ids, clock);
    }

    @Bean UpdateCarUseCase updateCarUseCase(CarRepository cars) {
        return new UpdateCarUseCase(cars);
    }

    @Bean DeleteCarUseCase deleteCarUseCase(CarRepository cars) {
        return new DeleteCarUseCase(cars);
    }
}
