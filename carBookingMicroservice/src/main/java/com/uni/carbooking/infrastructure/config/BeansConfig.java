package com.uni.carbooking.infrastructure.config;

import com.uni.carbooking.application.port.out.Clock;
import com.uni.carbooking.application.port.out.IdGenerator;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.time.Instant;
import java.util.UUID;

@Configuration
class BeansConfig {

    @Bean
    Clock systemClock() {
        return () -> Instant.now();
    }

    @Bean
    IdGenerator uuidGenerator() {
        return () -> UUID.randomUUID().toString();
    }
}
