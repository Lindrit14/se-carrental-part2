CREATE TABLE cars (
    id                  VARCHAR(36) PRIMARY KEY,
    brand               VARCHAR(80)  NOT NULL,
    model               VARCHAR(120) NOT NULL,
    license_plate       VARCHAR(40)  NOT NULL UNIQUE,
    daily_rate_amount   NUMERIC(19,4) NOT NULL CHECK (daily_rate_amount >= 0),
    daily_rate_currency VARCHAR(3)    NOT NULL,
    created_at          TIMESTAMPTZ   NOT NULL
);
