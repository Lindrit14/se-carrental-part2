package com.uni.carbooking.application.car;

import com.uni.carbooking.domain.booking.BookingRepository;
import com.uni.carbooking.domain.car.CarRepository;
import com.uni.carbooking.domain.error.CarHasBookings;
import com.uni.carbooking.domain.error.CarNotFound;

public class DeleteCarUseCase {

    private final CarRepository cars;
    private final BookingRepository bookings;

    public DeleteCarUseCase(CarRepository cars, BookingRepository bookings) {
        this.cars = cars;
        this.bookings = bookings;
    }

    public void execute(String id) {
        if (cars.findById(id).isEmpty()) {
            throw new CarNotFound(id);
        }
        if (bookings.existsForCar(id)) {
            throw new CarHasBookings(id);
        }
        cars.deleteById(id);
    }
}
