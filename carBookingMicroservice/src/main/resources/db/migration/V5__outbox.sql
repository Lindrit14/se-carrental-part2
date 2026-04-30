CREATE TABLE outbox (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id      VARCHAR(36) NOT NULL UNIQUE,
    aggregate_id  VARCHAR(36) NOT NULL,
    event_type    VARCHAR(100) NOT NULL,
    routing_key   VARCHAR(100) NOT NULL,
    payload       JSONB NOT NULL,
    occurred_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at  TIMESTAMPTZ,
    attempts      INT NOT NULL DEFAULT 0,
    last_error    TEXT
);

-- Hot path: scheduled publisher locks unpublished rows in occurred_at order.
CREATE INDEX idx_outbox_unpublished ON outbox(occurred_at) WHERE published_at IS NULL;
