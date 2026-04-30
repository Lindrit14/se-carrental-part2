package com.uni.carbooking.infrastructure.http;

import com.uni.carbooking.application.port.out.CarServicePort;
import com.uni.carbooking.domain.car.CarSnapshot;
import com.uni.carbooking.domain.error.CarNotFound;
import com.uni.carbooking.domain.error.CarServiceUnavailable;
import com.uni.carbooking.domain.money.Money;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatusCode;
import org.springframework.stereotype.Component;
import org.springframework.web.client.ResourceAccessException;
import org.springframework.web.client.RestClient;

import java.math.BigDecimal;

@Component
public class CarServiceClient implements CarServicePort {

    private static final Logger log = LoggerFactory.getLogger(CarServiceClient.class);

    private final RestClient restClient;

    public CarServiceClient(RestClient.Builder builder,
                            @org.springframework.beans.factory.annotation.Value("${car.service.url:http://localhost:8082}") String carServiceUrl) {
        this.restClient = builder.baseUrl(carServiceUrl).build();
    }

    @Override
    public CarSnapshot fetchCar(String carId) {
        try {
            CarDto dto = restClient.get()
                    .uri("/api/v1/cars/{id}", carId)
                    .retrieve()
                    .onStatus(HttpStatusCode::is4xxClientError, (req, resp) -> {
                        if (resp.getStatusCode().value() == 404) {
                            throw new CarNotFound(carId);
                        }
                        throw new CarServiceUnavailable("car service returned " + resp.getStatusCode());
                    })
                    .onStatus(HttpStatusCode::is5xxServerError, (req, resp) -> {
                        throw new CarServiceUnavailable("car service error: " + resp.getStatusCode());
                    })
                    .body(CarDto.class);

            if (dto == null) {
                throw new CarNotFound(carId);
            }
            return new CarSnapshot(
                    dto.id(),
                    dto.brand(),
                    dto.model(),
                    dto.licensePlate(),
                    new Money(dto.dailyRateAmount(), dto.dailyRateCurrency())
            );
        } catch (CarNotFound | CarServiceUnavailable e) {
            throw e;
        } catch (ResourceAccessException e) {
            log.warn("car service unreachable: {}", e.getMessage());
            throw new CarServiceUnavailable("car service unreachable: " + e.getMessage());
        }
    }

    record CarDto(
            String id,
            String brand,
            String model,
            String licensePlate,
            BigDecimal dailyRateAmount,
            String dailyRateCurrency,
            String location,
            String category
    ) {}
}
