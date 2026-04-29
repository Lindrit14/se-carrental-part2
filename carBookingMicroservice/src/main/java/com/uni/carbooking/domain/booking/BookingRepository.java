package com.uni.carbooking.domain.booking;

import java.util.List;
import java.util.Optional;

public interface BookingRepository {
    void save(Booking booking);
    Optional<Booking> findById(String id);
    List<Booking> findByCustomerId(String customerId);
    List<Booking> findAll();
}
