package com.uni.carbooking.interfaces.rest;

import com.uni.carbooking.application.booking.ListAllBookingsUseCase;
import com.uni.carbooking.interfaces.rest.dto.AdminBookingResponse;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/api/v1/admin")
@PreAuthorize("hasRole('admin')")
class AdminBookingController {

    private final ListAllBookingsUseCase listAll;

    AdminBookingController(ListAllBookingsUseCase listAll) {
        this.listAll = listAll;
    }

    @GetMapping("/bookings")
    List<AdminBookingResponse> all() {
        return listAll.execute().stream().map(AdminBookingResponse::from).toList();
    }
}
