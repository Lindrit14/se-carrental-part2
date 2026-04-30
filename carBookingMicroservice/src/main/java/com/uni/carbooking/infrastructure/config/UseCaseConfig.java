package com.uni.carbooking.infrastructure.config;

import com.uni.carbooking.application.booking.CancelBookingUseCase;
import com.uni.carbooking.application.booking.CreateBookingUseCase;
import com.uni.carbooking.application.booking.ListAllBookingsUseCase;
import com.uni.carbooking.application.booking.ListMyBookingsUseCase;
import com.uni.carbooking.application.customer.AnonymizeCustomerUseCase;
import com.uni.carbooking.application.customer.RegisterCustomerUseCase;
import com.uni.carbooking.application.customer.UpdateCustomerEmailUseCase;
import com.uni.carbooking.application.port.out.CarServicePort;
import com.uni.carbooking.application.port.out.Clock;
import com.uni.carbooking.application.port.out.CurrencyConverterPort;
import com.uni.carbooking.application.port.out.EventPublisher;
import com.uni.carbooking.application.port.out.IdGenerator;
import com.uni.carbooking.domain.booking.BookingRepository;
import com.uni.carbooking.domain.customer.CustomerRepository;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
class UseCaseConfig {

    @Bean RegisterCustomerUseCase registerCustomerUseCase(CustomerRepository customers, IdGenerator ids, Clock clock) {
        return new RegisterCustomerUseCase(customers, ids, clock);
    }

    @Bean UpdateCustomerEmailUseCase updateCustomerEmailUseCase(CustomerRepository customers, Clock clock) {
        return new UpdateCustomerEmailUseCase(customers, clock);
    }

    @Bean AnonymizeCustomerUseCase anonymizeCustomerUseCase(CustomerRepository customers, Clock clock) {
        return new AnonymizeCustomerUseCase(customers, clock);
    }

    @Bean CreateBookingUseCase createBookingUseCase(
            BookingRepository bookings, CarServicePort carService, CustomerRepository customers,
            CurrencyConverterPort converter, EventPublisher publisher, Clock clock, IdGenerator ids) {
        return new CreateBookingUseCase(bookings, carService, customers, converter, publisher, clock, ids);
    }

    @Bean CancelBookingUseCase cancelBookingUseCase(
            BookingRepository bookings, CustomerRepository customers,
            EventPublisher publisher, Clock clock) {
        return new CancelBookingUseCase(bookings, customers, publisher, clock);
    }

    @Bean ListMyBookingsUseCase listMyBookingsUseCase(BookingRepository bookings, CustomerRepository customers) {
        return new ListMyBookingsUseCase(bookings, customers);
    }

    @Bean ListAllBookingsUseCase listAllBookingsUseCase(BookingRepository bookings) {
        return new ListAllBookingsUseCase(bookings);
    }
}
