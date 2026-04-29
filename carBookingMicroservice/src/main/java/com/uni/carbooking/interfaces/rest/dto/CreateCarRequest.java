package com.uni.carbooking.interfaces.rest.dto;

import jakarta.validation.constraints.DecimalMin;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;

import java.math.BigDecimal;

public record CreateCarRequest(
        @NotBlank @Size(max = 80) String brand,
        @NotBlank @Size(max = 120) String model,
        @NotBlank @Size(max = 40) String licensePlate,
        @DecimalMin("0.0") BigDecimal dailyRateAmount,
        @NotBlank @Pattern(regexp = "^[A-Za-z]{3}$") String dailyRateCurrency
) {}
