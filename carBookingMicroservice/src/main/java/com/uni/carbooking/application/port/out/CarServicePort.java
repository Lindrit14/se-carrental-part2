package com.uni.carbooking.application.port.out;

import com.uni.carbooking.domain.car.CarSnapshot;
import com.uni.carbooking.domain.error.CarNotFound;
import com.uni.carbooking.domain.error.CarServiceUnavailable;

/**
 * Port: the booking service's view of the car catalog.
 * Implemented by the HTTP client adapter in infrastructure.
 */
public interface CarServicePort {
    /**
     * Fetch a car snapshot by ID from the car service.
     *
     * @throws CarNotFound           if the car does not exist (404)
     * @throws CarServiceUnavailable if the car service is unreachable
     */
    CarSnapshot fetchCar(String carId);
}
