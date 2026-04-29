package com.uni.carbooking.interfaces.rest;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
class HealthController {

    @GetMapping("/healthz")
    Map<String, String> liveness() {
        return Map.of("status", "ok");
    }

    @GetMapping("/readyz")
    Map<String, String> readiness() {
        // Spring Actuator at /actuator/health gives full DB/AMQP readiness;
        // /readyz is a stable shorthand the platform compose can probe.
        return Map.of("status", "ok");
    }
}
