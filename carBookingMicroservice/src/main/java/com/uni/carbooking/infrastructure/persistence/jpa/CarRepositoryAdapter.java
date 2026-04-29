package com.uni.carbooking.infrastructure.persistence.jpa;

import com.uni.carbooking.domain.car.Car;
import com.uni.carbooking.domain.car.CarRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
class CarRepositoryAdapter implements CarRepository {

    private final CarSpringRepo repo;

    CarRepositoryAdapter(CarSpringRepo repo) {
        this.repo = repo;
    }

    @Override
    public void save(Car car) {
        repo.save(CarJpaEntity.fromDomain(car));
    }

    @Override
    public Optional<Car> findById(String id) {
        return repo.findById(id).map(CarJpaEntity::toDomain);
    }

    @Override
    public List<Car> findAll() {
        return repo.findAll().stream().map(CarJpaEntity::toDomain).toList();
    }
}
