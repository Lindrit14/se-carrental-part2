package com.uni.carbooking.interfaces.rest;

import com.uni.carbooking.application.car.AddCarUseCase;
import com.uni.carbooking.application.car.GetCarUseCase;
import com.uni.carbooking.application.car.ListCarsUseCase;
import com.uni.carbooking.application.car.SearchCarsUseCase;
import com.uni.carbooking.domain.car.CarCategory;
import com.uni.carbooking.domain.money.Money;
import com.uni.carbooking.interfaces.rest.dto.CarResponse;
import com.uni.carbooking.interfaces.rest.dto.CreateCarRequest;
import jakarta.validation.Valid;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.time.LocalDate;
import java.util.List;

@RestController
@RequestMapping("/api/v1/cars")
class CarController {

    private final ListCarsUseCase list;
    private final SearchCarsUseCase search;
    private final GetCarUseCase get;
    private final AddCarUseCase add;

    CarController(ListCarsUseCase list, SearchCarsUseCase search, GetCarUseCase get, AddCarUseCase add) {
        this.list = list;
        this.search = search;
        this.get = get;
        this.add = add;
    }

    @GetMapping
    List<CarResponse> list() {
        return list.execute().stream().map(CarResponse::from).toList();
    }

    @GetMapping("/search")
    List<CarResponse> search(
            @RequestParam(required = false) String location,
            @RequestParam(required = false) CarCategory category,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate from,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate to) {
        return search.execute(new SearchCarsUseCase.Input(location, category, from, to))
                .stream().map(CarResponse::from).toList();
    }

    @GetMapping("/{id}")
    CarResponse one(@PathVariable String id) {
        return CarResponse.from(get.execute(id));
    }

    @PostMapping
    @PreAuthorize("hasAuthority('SCOPE_admin') or hasRole('admin')")
    ResponseEntity<CarResponse> add(@Valid @RequestBody CreateCarRequest req) {
        var car = add.execute(new AddCarUseCase.Input(
                req.brand(), req.model(), req.licensePlate(),
                new Money(req.dailyRateAmount(), req.dailyRateCurrency()),
                req.location(), req.category()
        ));
        return ResponseEntity.status(HttpStatus.CREATED).body(CarResponse.from(car));
    }
}
