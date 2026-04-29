package com.uni.carbooking;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity;

@SpringBootApplication
@EnableMethodSecurity
public class CarBookingApplication {
    public static void main(String[] args) {
        SpringApplication.run(CarBookingApplication.class, args);
    }
}
