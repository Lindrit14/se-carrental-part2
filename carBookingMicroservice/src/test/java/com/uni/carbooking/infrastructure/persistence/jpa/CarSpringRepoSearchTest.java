package com.uni.carbooking.infrastructure.persistence.jpa;

import com.uni.carbooking.domain.car.CarCategory;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jdbc.AutoConfigureTestDatabase;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.sql.Timestamp;
import java.time.Instant;
import java.time.LocalDate;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;

/**
 * Verifies the JPQL search query: location substring match (case-insensitive),
 * category filter, and the {@code NOT EXISTS} overlap exclusion.
 *
 * Uses {@link JdbcTemplate} for seeding so the test does not depend on the
 * private repository adapters.
 */
@DataJpaTest
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
@Testcontainers
class CarSpringRepoSearchTest {

    @Container
    @SuppressWarnings("resource") // Testcontainers lifecycle managed by JUnit ext.
    static final PostgreSQLContainer<?> POSTGRES = new PostgreSQLContainer<>("postgres:16-alpine")
            .withDatabaseName("booking")
            .withUsername("booking")
            .withPassword("booking");

    @DynamicPropertySource
    static void datasourceProps(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", POSTGRES::getJdbcUrl);
        registry.add("spring.datasource.username", POSTGRES::getUsername);
        registry.add("spring.datasource.password", POSTGRES::getPassword);
        registry.add("spring.flyway.enabled", () -> "true");
        registry.add("spring.jpa.hibernate.ddl-auto", () -> "validate");
    }

    @Autowired CarSpringRepo cars;
    @Autowired JdbcTemplate jdbc;

    private String viennaCarId;
    private String munichCarId;

    @BeforeEach
    void seed() {
        jdbc.update("DELETE FROM bookings");
        jdbc.update("DELETE FROM customers");
        jdbc.update("DELETE FROM cars");

        viennaCarId = insertCar("Volkswagen", "Golf", "VIE-1",
                "Vienna International Airport, Vienna, Austria", CarCategory.MEDIUM);
        munichCarId = insertCar("BMW", "X5", "MUC-1",
                "Munich Airport, Munich, Germany", CarCategory.SUV);

        var customerId = insertCustomer("ext-user-1");
        insertConfirmedBooking(customerId, viennaCarId,
                LocalDate.of(2026, 5, 1), LocalDate.of(2026, 5, 8));
    }

    @Test
    void filters_by_location_substring_case_insensitive() {
        var results = cars.search("vienna", null, null, null);
        assertThat(results).hasSize(1);
        assertThat(results.get(0).toDomain().id()).isEqualTo(viennaCarId);
    }

    @Test
    void filters_by_category() {
        var results = cars.search(null, CarCategory.SUV, null, null);
        assertThat(results).hasSize(1);
        assertThat(results.get(0).toDomain().id()).isEqualTo(munichCarId);
    }

    @Test
    void excludes_cars_with_overlapping_confirmed_bookings() {
        // Vienna car has a CONFIRMED 2026-05-01 -> 2026-05-08 booking.
        // Overlapping window (May 5 -> May 10) -> Vienna car excluded.
        var overlapping = cars.search(null, null,
                LocalDate.of(2026, 5, 5), LocalDate.of(2026, 5, 10));
        assertThat(overlapping).hasSize(1);
        assertThat(overlapping.get(0).toDomain().id()).isEqualTo(munichCarId);

        // Non-overlapping window (May 15 -> May 20) -> both visible.
        var both = cars.search(null, null,
                LocalDate.of(2026, 5, 15), LocalDate.of(2026, 5, 20));
        assertThat(both).hasSize(2);
    }

    @Test
    void empty_filters_returns_all() {
        var results = cars.search(null, null, null, null);
        assertThat(results).hasSize(2);
    }

    private String insertCar(String brand, String model, String plate, String location, CarCategory cat) {
        var id = UUID.randomUUID().toString();
        jdbc.update(
                "INSERT INTO cars (id, brand, model, license_plate, daily_rate_amount, daily_rate_currency, location, category, created_at) "
                        + "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
                id, brand, model, plate, "45.00", "EUR", location, cat.name(),
                Timestamp.from(Instant.now()));
        return id;
    }

    private String insertCustomer(String externalUserId) {
        var id = UUID.randomUUID().toString();
        var now = Timestamp.from(Instant.now());
        jdbc.update(
                "INSERT INTO customers (id, external_user_id, email, anonymized, created_at, updated_at) "
                        + "VALUES (?, ?, ?, ?, ?, ?)",
                id, externalUserId, externalUserId + "@example.com", false, now, now);
        return id;
    }

    private void insertConfirmedBooking(String customerId, String carId, LocalDate start, LocalDate end) {
        var now = Timestamp.from(Instant.now());
        jdbc.update(
                "INSERT INTO bookings (id, customer_id, car_id, start_date, end_date, status, "
                        + "total_source_amount, total_source_currency, total_target_amount, total_target_currency, created_at, updated_at) "
                        + "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
                UUID.randomUUID().toString(), customerId, carId,
                java.sql.Date.valueOf(start), java.sql.Date.valueOf(end), "CONFIRMED",
                "315.00", "EUR", "315.00", "EUR", now, now);
    }
}
