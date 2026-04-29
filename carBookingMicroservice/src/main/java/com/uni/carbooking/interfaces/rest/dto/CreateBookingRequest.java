package com.uni.carbooking.interfaces.rest.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;

import java.time.LocalDate;

public record CreateBookingRequest(
        @NotBlank String carId,
        @NotNull LocalDate startDate,
        @NotNull LocalDate endDate,
        @NotBlank @Pattern(regexp = "^[A-Za-z]{3}$") String targetCurrency
) {}
