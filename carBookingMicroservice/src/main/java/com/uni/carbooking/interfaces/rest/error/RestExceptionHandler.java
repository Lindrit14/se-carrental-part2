package com.uni.carbooking.interfaces.rest.error;

import com.uni.carbooking.domain.error.BookingNotFound;
import com.uni.carbooking.domain.error.BookingNotOwned;
import com.uni.carbooking.domain.error.CarNotFound;
import com.uni.carbooking.domain.error.CurrencyConversionFailed;
import com.uni.carbooking.domain.error.CustomerNotFound;
import com.uni.carbooking.domain.error.InvalidDateRange;
import com.uni.carbooking.interfaces.rest.dto.ErrorResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
class RestExceptionHandler {

    private static final Logger log = LoggerFactory.getLogger(RestExceptionHandler.class);

    @ExceptionHandler({CarNotFound.class, BookingNotFound.class, CustomerNotFound.class})
    ResponseEntity<ErrorResponse> notFound(RuntimeException e) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND)
                .body(new ErrorResponse("not_found", e.getMessage()));
    }

    @ExceptionHandler(BookingNotOwned.class)
    ResponseEntity<ErrorResponse> notOwned(BookingNotOwned e) {
        return ResponseEntity.status(HttpStatus.FORBIDDEN)
                .body(new ErrorResponse("forbidden", e.getMessage()));
    }

    @ExceptionHandler(InvalidDateRange.class)
    ResponseEntity<ErrorResponse> badRange(InvalidDateRange e) {
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
                .body(new ErrorResponse("invalid_date_range", e.getMessage()));
    }

    @ExceptionHandler(CurrencyConversionFailed.class)
    ResponseEntity<ErrorResponse> currency(CurrencyConversionFailed e) {
        log.warn("currency conversion failed: {}", e.getMessage());
        return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE)
                .body(new ErrorResponse("currency_conversion_failed", e.getMessage()));
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    ResponseEntity<ErrorResponse> validation(MethodArgumentNotValidException e) {
        var first = e.getBindingResult().getFieldErrors().stream().findFirst();
        var message = first.map(fe -> fe.getField() + ": " + fe.getDefaultMessage())
                .orElse("validation failed");
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
                .body(new ErrorResponse("validation_failed", message));
    }

    @ExceptionHandler(IllegalArgumentException.class)
    ResponseEntity<ErrorResponse> illegal(IllegalArgumentException e) {
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
                .body(new ErrorResponse("invalid_request", e.getMessage()));
    }
}
