package com.uni.carbooking.infrastructure.persistence.jpa;

import org.springframework.data.jpa.repository.JpaRepository;

interface CarSpringRepo extends JpaRepository<CarJpaEntity, String> {
}
