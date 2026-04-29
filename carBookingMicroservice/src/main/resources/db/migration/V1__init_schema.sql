CREATE TABLE cars (
    id                  VARCHAR(36) PRIMARY KEY,
    brand               VARCHAR(80)  NOT NULL,
    model               VARCHAR(120) NOT NULL,
    license_plate       VARCHAR(40)  NOT NULL UNIQUE,
    daily_rate_amount   NUMERIC(19,4) NOT NULL CHECK (daily_rate_amount >= 0),
    daily_rate_currency VARCHAR(3)    NOT NULL,
    created_at          TIMESTAMPTZ   NOT NULL
);

CREATE TABLE customers (
    id               VARCHAR(36) PRIMARY KEY,
    external_user_id VARCHAR(36) NOT NULL UNIQUE,
    email            VARCHAR(255),
    anonymized       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

CREATE TABLE bookings (
    id                    VARCHAR(36) PRIMARY KEY,
    customer_id           VARCHAR(36) NOT NULL REFERENCES customers(id),
    car_id                VARCHAR(36) NOT NULL REFERENCES cars(id),
    start_date            DATE        NOT NULL,
    end_date              DATE        NOT NULL,
    status                VARCHAR(16) NOT NULL,
    total_source_amount   NUMERIC(19,4) NOT NULL,
    total_source_currency VARCHAR(3)    NOT NULL,
    total_target_amount   NUMERIC(19,4) NOT NULL,
    total_target_currency VARCHAR(3)    NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL,
    updated_at            TIMESTAMPTZ NOT NULL,
    CHECK (end_date > start_date)
);

CREATE INDEX idx_bookings_customer_created ON bookings(customer_id, created_at DESC);
CREATE INDEX idx_bookings_car              ON bookings(car_id);
