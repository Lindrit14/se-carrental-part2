package notification

import (
	"encoding/json"
	"time"
)

// EventEnvelope is the wire-format shared by all platform services.
// Matches the envelope used by user-auth (Go) and booking (Java).
type EventEnvelope struct {
	EventID    string          `json:"event_id"`
	EventType  string          `json:"event_type"`
	Version    string          `json:"version"`
	OccurredAt time.Time       `json:"occurred_at"`
	Data       json.RawMessage `json:"data"`
}
