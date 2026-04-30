-- Phase 1: Add car snapshot columns to bookings.
-- The bookings.car_id column is kept as a historical reference but the FK
-- constraint is dropped — the cars table now lives in a separate service/DB.

ALTER TABLE bookings ADD COLUMN car_brand               VARCHAR(80)    NOT NULL DEFAULT '';
ALTER TABLE bookings ADD COLUMN car_model               VARCHAR(120)   NOT NULL DEFAULT '';
ALTER TABLE bookings ADD COLUMN car_license_plate       VARCHAR(40)    NOT NULL DEFAULT '';
ALTER TABLE bookings ADD COLUMN daily_rate_snapshot_amount   NUMERIC(19,4) NOT NULL DEFAULT 0;
ALTER TABLE bookings ADD COLUMN daily_rate_snapshot_currency VARCHAR(3)    NOT NULL DEFAULT 'EUR';

-- Backfill snapshot from existing cars table (dev / demo data only).
UPDATE bookings b
SET
    car_brand                    = c.brand,
    car_model                    = c.model,
    car_license_plate            = c.license_plate,
    daily_rate_snapshot_amount   = c.daily_rate_amount,
    daily_rate_snapshot_currency = c.daily_rate_currency
FROM cars c
WHERE b.car_id = c.id;

-- Drop the FK so bookings can exist independently of the cars table.
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_car_id_fkey;

-- Remove defaults now that backfill is done.
ALTER TABLE bookings ALTER COLUMN car_brand DROP DEFAULT;
ALTER TABLE bookings ALTER COLUMN car_model DROP DEFAULT;
ALTER TABLE bookings ALTER COLUMN car_license_plate DROP DEFAULT;
ALTER TABLE bookings ALTER COLUMN daily_rate_snapshot_amount DROP DEFAULT;
ALTER TABLE bookings ALTER COLUMN daily_rate_snapshot_currency DROP DEFAULT;
