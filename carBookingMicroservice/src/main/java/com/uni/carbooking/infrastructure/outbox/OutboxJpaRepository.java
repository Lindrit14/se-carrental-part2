package com.uni.carbooking.infrastructure.outbox;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.UUID;

interface OutboxJpaRepository extends JpaRepository<OutboxJpaEntity, UUID> {

    /**
     * Reads the next batch of unpublished events, locking them so a parallel
     * publisher (multi-instance deploy) cannot grab the same rows. Native
     * query is used because JPA's pagination + SKIP LOCKED interaction is
     * brittle across versions.
     */
    @Query(value = """
            SELECT * FROM outbox
            WHERE published_at IS NULL
            ORDER BY occurred_at ASC
            LIMIT :limit
            FOR UPDATE SKIP LOCKED
            """, nativeQuery = true)
    List<OutboxJpaEntity> findUnpublishedForUpdate(@Param("limit") int limit);
}
