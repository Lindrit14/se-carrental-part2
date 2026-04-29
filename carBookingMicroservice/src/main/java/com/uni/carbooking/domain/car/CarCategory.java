package com.uni.carbooking.domain.car;

/**
 * Vehicle category — used for the search/filter chips on the public listings.
 * Persisted as the enum name in the {@code cars.category} column.
 */
public enum CarCategory {
    SMALL,
    MEDIUM,
    LARGE,
    SUV,
    PEOPLE_CARRIER,
    PREMIUM
}
