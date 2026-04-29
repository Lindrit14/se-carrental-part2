package com.uni.carbooking.interfaces.rest;

import com.uni.carbooking.application.booking.CancelBookingUseCase;
import com.uni.carbooking.application.booking.CreateBookingUseCase;
import com.uni.carbooking.application.booking.ListMyBookingsUseCase;
import com.uni.carbooking.interfaces.rest.dto.BookingResponse;
import com.uni.carbooking.interfaces.rest.dto.CreateBookingRequest;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/api/v1/bookings")
class BookingController {

    private final CreateBookingUseCase create;
    private final ListMyBookingsUseCase listMine;
    private final CancelBookingUseCase cancel;

    BookingController(CreateBookingUseCase create, ListMyBookingsUseCase listMine, CancelBookingUseCase cancel) {
        this.create = create;
        this.listMine = listMine;
        this.cancel = cancel;
    }

    @PostMapping
    ResponseEntity<BookingResponse> create(@Valid @RequestBody CreateBookingRequest req,
                                           @AuthenticationPrincipal Jwt jwt) {
        var booking = create.execute(new CreateBookingUseCase.Input(
                jwt.getSubject(),
                req.carId(),
                req.startDate(),
                req.endDate(),
                req.targetCurrency()
        ));
        return ResponseEntity.status(HttpStatus.CREATED).body(BookingResponse.from(booking));
    }

    @GetMapping("/me")
    List<BookingResponse> mine(@AuthenticationPrincipal Jwt jwt) {
        return listMine.execute(jwt.getSubject()).stream().map(BookingResponse::from).toList();
    }

    @DeleteMapping("/{id}")
    ResponseEntity<Void> cancel(@PathVariable String id, @AuthenticationPrincipal Jwt jwt) {
        cancel.execute(id, jwt.getSubject());
        return ResponseEntity.noContent().build();
    }
}
