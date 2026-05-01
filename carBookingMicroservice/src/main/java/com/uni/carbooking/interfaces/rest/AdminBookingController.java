package com.uni.carbooking.interfaces.rest;

import com.uni.carbooking.application.booking.ListAllBookingsUseCase;
import com.uni.carbooking.interfaces.rest.dto.AdminBookingResponse;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/api/v1/admin")
@PreAuthorize("hasRole('admin')")
class AdminBookingController {

    private final ListAllBookingsUseCase listAll;
    private final BookingPriceProjection priceProjection;

    AdminBookingController(ListAllBookingsUseCase listAll, BookingPriceProjection priceProjection) {
        this.listAll = listAll;
        this.priceProjection = priceProjection;
    }

    @GetMapping("/bookings")
    List<AdminBookingResponse> all(@RequestParam(required = false) String targetCurrency) {
        var displayTotal = priceProjection.displayTotalFor(targetCurrency);
        return listAll.execute().stream()
                .map(b -> AdminBookingResponse.fromWithDisplayTotal(b, displayTotal.totalFor(b)))
                .toList();
    }
}
