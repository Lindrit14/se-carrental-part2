package com.uni.car.infrastructure.persistence.jpa;

import com.uni.car.domain.car.CarCategory;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;

interface CarSpringRepo extends JpaRepository<CarJpaEntity, String> {

    @Query("""
        SELECT c FROM CarJpaEntity c
        WHERE (:loc IS NULL OR LOWER(c.location) LIKE :loc)
          AND (:cat IS NULL OR c.category = :cat)
        ORDER BY c.createdAt
        """)
    List<CarJpaEntity> search(@Param("loc") String loc, @Param("cat") CarCategory cat);
}
