-- Seed a few demo cars so the platform demo has something to list / book.
INSERT INTO cars (id, brand, model, license_plate, daily_rate_amount, daily_rate_currency, created_at) VALUES
  ('11111111-1111-1111-1111-111111111111', 'Volkswagen', 'Golf',     'M-VW-1234',  45.00, 'EUR', NOW()),
  ('22222222-2222-2222-2222-222222222222', 'BMW',        '3 Series', 'M-BM-5678',  85.00, 'EUR', NOW()),
  ('33333333-3333-3333-3333-333333333333', 'Tesla',      'Model 3',  'M-TS-9012', 110.00, 'USD', NOW()),
  ('44444444-4444-4444-4444-444444444444', 'Toyota',     'Yaris',    'M-TY-3456',  35.00, 'EUR', NOW());
