package com.uni.car.interfaces.rest;

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
        return Map.of("status", "ok");
    }
}
