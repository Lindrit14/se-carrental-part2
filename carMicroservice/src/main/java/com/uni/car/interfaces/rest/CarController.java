package com.uni.car.interfaces.rest;

import com.uni.car.application.car.AddCarUseCase;
import com.uni.car.application.car.DeleteCarUseCase;
import com.uni.car.application.car.GetCarUseCase;
import com.uni.car.application.car.ListCarsUseCase;
import com.uni.car.application.car.SearchCarsUseCase;
import com.uni.car.application.car.UpdateCarUseCase;
import com.uni.car.domain.car.CarCategory;
import com.uni.car.domain.money.Money;
import com.uni.car.interfaces.rest.dto.CarResponse;
import com.uni.car.interfaces.rest.dto.CreateCarRequest;
import com.uni.car.interfaces.rest.dto.UpdateCarRequest;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/api/v1/cars")
class CarController {

    private final ListCarsUseCase list;
    private final SearchCarsUseCase search;
    private final GetCarUseCase get;
    private final AddCarUseCase add;
    private final UpdateCarUseCase update;
    private final DeleteCarUseCase delete;
    private final CarPriceProjection priceProjection;

    CarController(ListCarsUseCase list, SearchCarsUseCase search, GetCarUseCase get,
                  AddCarUseCase add, UpdateCarUseCase update, DeleteCarUseCase delete,
                  CarPriceProjection priceProjection) {
        this.list = list;
        this.search = search;
        this.get = get;
        this.add = add;
        this.update = update;
        this.delete = delete;
        this.priceProjection = priceProjection;
    }

    @GetMapping
    List<CarResponse> list(@RequestParam(required = false) String targetCurrency) {
        var renderer = priceProjection.rendererFor(targetCurrency);
        return list.execute().stream().map(renderer::render).toList();
    }

    @GetMapping("/search")
    List<CarResponse> search(
            @RequestParam(required = false) String location,
            @RequestParam(required = false) CarCategory category,
            @RequestParam(required = false) String targetCurrency) {
        var renderer = priceProjection.rendererFor(targetCurrency);
        return search.execute(new SearchCarsUseCase.Input(location, category))
                .stream().map(renderer::render).toList();
    }

    @GetMapping("/{id}")
    CarResponse one(@PathVariable String id,
                    @RequestParam(required = false) String targetCurrency) {
        return priceProjection.rendererFor(targetCurrency).render(get.execute(id));
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

    @PutMapping("/{id}")
    @PreAuthorize("hasAuthority('SCOPE_admin') or hasRole('admin')")
    CarResponse update(@PathVariable String id, @Valid @RequestBody UpdateCarRequest req) {
        var car = update.execute(new UpdateCarUseCase.Input(
                id,
                req.brand(), req.model(), req.licensePlate(),
                new Money(req.dailyRateAmount(), req.dailyRateCurrency()),
                req.location(), req.category()
        ));
        return CarResponse.from(car);
    }

    @DeleteMapping("/{id}")
    @PreAuthorize("hasAuthority('SCOPE_admin') or hasRole('admin')")
    ResponseEntity<Void> delete(@PathVariable String id) {
        delete.execute(id);
        return ResponseEntity.noContent().build();
    }
}
