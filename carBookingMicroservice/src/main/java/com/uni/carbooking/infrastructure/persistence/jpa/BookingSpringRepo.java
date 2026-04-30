package com.uni.carbooking.infrastructure.persistence.jpa;

import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

interface BookingSpringRepo extends JpaRepository<BookingJpaEntity, String> {
    List<BookingJpaEntity> findByCustomerIdOrderByCreatedAtDesc(String customerId);
    List<BookingJpaEntity> findAllByOrderByCreatedAtDesc();
}
