package com.uni.car.infrastructure.persistence.jpa;

import com.uni.car.domain.car.Car;
import com.uni.car.domain.car.CarCategory;
import com.uni.car.domain.car.CarRepository;
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

    @Override
    public void deleteById(String id) {
        repo.deleteById(id);
    }

    @Override
    public List<Car> search(String location, CarCategory category) {
        return repo.search(location, category).stream()
                .map(CarJpaEntity::toDomain)
                .toList();
    }
}
