package com.uni.carbooking.infrastructure.inbox;

import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

import java.util.Optional;

interface EventLogJpaRepository extends JpaRepository<EventLogJpaEntity, UUID> {
    boolean existsByEventId(String eventId);
    Optional<EventLogJpaEntity> findByEventId(String eventId);
}
