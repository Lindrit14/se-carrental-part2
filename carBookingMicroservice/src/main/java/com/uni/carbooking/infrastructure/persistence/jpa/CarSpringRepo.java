package com.uni.carbooking.infrastructure.persistence.jpa;

import com.uni.carbooking.domain.car.CarCategory;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.time.LocalDate;
import java.util.List;

interface CarSpringRepo extends JpaRepository<CarJpaEntity, String> {

    /**
     * Search with optional filters. {@code loc} is matched as a case-insensitive
     * substring; {@code cat} is exact; the date range, when provided, excludes
     * cars with any overlapping {@code CONFIRMED} booking via {@code NOT EXISTS}.
     */
    @Query("""
        SELECT c FROM CarJpaEntity c
        WHERE (:loc IS NULL OR LOWER(c.location) LIKE LOWER(CONCAT('%', :loc, '%')))
          AND (:cat IS NULL OR c.category = :cat)
          AND (CAST(:from AS date) IS NULL OR CAST(:to AS date) IS NULL OR NOT EXISTS (
              SELECT 1 FROM BookingJpaEntity b
               WHERE b.carId = c.id
                 AND b.status = com.uni.carbooking.domain.booking.BookingStatus.CONFIRMED
                 AND b.startDate < :to
                 AND b.endDate   > :from
          ))
        ORDER BY c.createdAt
        """)
    List<CarJpaEntity> search(@Param("loc") String loc,
                              @Param("cat") CarCategory cat,
                              @Param("from") LocalDate from,
                              @Param("to") LocalDate to);
}
