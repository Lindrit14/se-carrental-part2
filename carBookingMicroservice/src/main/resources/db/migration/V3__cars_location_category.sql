-- Add location + category columns to cars and backfill the demo rows.
ALTER TABLE cars ADD COLUMN location VARCHAR(160) NOT NULL DEFAULT '';
ALTER TABLE cars ADD COLUMN category VARCHAR(32)  NOT NULL DEFAULT 'MEDIUM';

CREATE INDEX idx_cars_location_lower ON cars (LOWER(location));
CREATE INDEX idx_cars_category       ON cars (category);

UPDATE cars SET location = 'Vienna International Airport, Vienna, Austria',  category = 'MEDIUM'  WHERE id = '11111111-1111-1111-1111-111111111111';
UPDATE cars SET location = 'Vienna International Airport, Vienna, Austria',  category = 'PREMIUM' WHERE id = '22222222-2222-2222-2222-222222222222';
UPDATE cars SET location = 'Munich Airport, Munich, Germany',                 category = 'PREMIUM' WHERE id = '33333333-3333-3333-3333-333333333333';
UPDATE cars SET location = 'Munich Airport, Munich, Germany',                 category = 'SMALL'   WHERE id = '44444444-4444-4444-4444-444444444444';

ALTER TABLE cars ALTER COLUMN location DROP DEFAULT;
ALTER TABLE cars ALTER COLUMN category DROP DEFAULT;
