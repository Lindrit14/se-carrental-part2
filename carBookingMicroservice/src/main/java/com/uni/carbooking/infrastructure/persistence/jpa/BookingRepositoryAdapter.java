package com.uni.carbooking.infrastructure.persistence.jpa;

import com.uni.carbooking.domain.booking.Booking;
import com.uni.carbooking.domain.booking.BookingRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
class BookingRepositoryAdapter implements BookingRepository {

    private final BookingSpringRepo repo;

    BookingRepositoryAdapter(BookingSpringRepo repo) {
        this.repo = repo;
    }

    @Override
    public void save(Booking booking) {
        repo.save(BookingJpaEntity.fromDomain(booking));
    }

    @Override
    public Optional<Booking> findById(String id) {
        return repo.findById(id).map(BookingJpaEntity::toDomain);
    }

    @Override
    public List<Booking> findByCustomerId(String customerId) {
        return repo.findByCustomerIdOrderByCreatedAtDesc(customerId)
                .stream()
                .map(BookingJpaEntity::toDomain)
                .toList();
    }

    @Override
    public List<Booking> findAll() {
        return repo.findAllByOrderByCreatedAtDesc()
                .stream()
                .map(BookingJpaEntity::toDomain)
                .toList();
    }
}
