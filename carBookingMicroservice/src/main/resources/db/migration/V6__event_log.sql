CREATE TABLE event_log (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id      VARCHAR(36) NOT NULL UNIQUE,
    event_type    VARCHAR(100) NOT NULL,
    source        VARCHAR(100) NOT NULL,
    payload       JSONB NOT NULL,
    received_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at  TIMESTAMPTZ
);

-- The event_id unique constraint already gives us a hash index; this index
-- supports the replay query (find unprocessed events in arrival order).
CREATE INDEX idx_event_log_unprocessed ON event_log(received_at) WHERE processed_at IS NULL;
