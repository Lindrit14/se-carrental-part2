package com.uni.carbooking.application.booking;

import com.uni.carbooking.application.event.BookingEvents;
import com.uni.carbooking.application.port.out.Clock;
import com.uni.carbooking.application.port.out.CurrencyConverterPort;
import com.uni.carbooking.application.port.out.EventPublisher;
import com.uni.carbooking.application.port.out.IdGenerator;
import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.booking.BookingRepository;
import com.uni.carbooking.domain.booking.BookingStatus;
import com.uni.carbooking.domain.car.Car;
import com.uni.carbooking.domain.car.CarRepository;
import com.uni.carbooking.domain.customer.Customer;
import com.uni.carbooking.domain.customer.CustomerRepository;
import com.uni.carbooking.domain.error.CarNotFound;
import com.uni.carbooking.domain.error.CustomerNotFound;
import com.uni.carbooking.domain.money.Money;

import java.time.LocalDate;

public class CreateBookingUseCase {

    private final BookingRepository bookings;
    private final CarRepository cars;
    private final CustomerRepository customers;
    private final CurrencyConverterPort converter;
    private final EventPublisher publisher;
    private final Clock clock;
    private final IdGenerator ids;

    public CreateBookingUseCase(BookingRepository bookings, CarRepository cars,
                                CustomerRepository customers, CurrencyConverterPort converter,
                                EventPublisher publisher, Clock clock, IdGenerator ids) {
        this.bookings = bookings;
        this.cars = cars;
        this.customers = customers;
        this.converter = converter;
        this.publisher = publisher;
        this.clock = clock;
        this.ids = ids;
    }

    public record Input(
            String externalUserId,   // from JWT subject
            String carId,
            LocalDate startDate,
            LocalDate endDate,
            String targetCurrency
    ) {}

    public Booking execute(Input in) {
        Customer customer = customers.findByExternalUserId(in.externalUserId())
                .orElseThrow(() -> new CustomerNotFound(in.externalUserId()));
        Car car = cars.findById(in.carId())
                .orElseThrow(() -> new CarNotFound(in.carId()));

        var now = clock.now();
        var totalSource = car.dailyRate().multiply(daysBetween(in.startDate(), in.endDate()));
        var totalTarget = converter.convert(totalSource, in.targetCurrency());

        var booking = new Booking(
                ids.newId(),
                customer.id(),
                car.id(),
                in.startDate(),
                in.endDate(),
                BookingStatus.CONFIRMED,
                totalSource,
                totalTarget,
                now,
                now
        );
        bookings.save(booking);

        publisher.publish(
                BookingEvents.BOOKING_CREATED,
                BookingEvents.BOOKING_CREATED,
                new BookingEvents.BookingCreated(
                        booking.id(),
                        booking.customerId(),
                        booking.carId(),
                        booking.startDate().toString(),
                        booking.endDate().toString(),
                        booking.totalSource().amount().toPlainString(),
                        booking.totalSource().currency(),
                        booking.totalTarget().amount().toPlainString(),
                        booking.totalTarget().currency()
                )
        );
        return booking;
    }

    private static long daysBetween(LocalDate start, LocalDate end) {
        return java.time.temporal.ChronoUnit.DAYS.between(start, end);
    }
}
