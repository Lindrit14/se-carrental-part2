package com.uni.carbooking.infrastructure.persistence.jpa;

import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;

interface CustomerSpringRepo extends JpaRepository<CustomerJpaEntity, String> {
    Optional<CustomerJpaEntity> findByExternalUserId(String externalUserId);
}
