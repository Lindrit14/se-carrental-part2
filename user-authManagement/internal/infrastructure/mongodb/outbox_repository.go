package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const outboxCollection = "outbox"

// OutboxEvent is the BSON document shape for an outbox row.
type OutboxEvent struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	EventID     string        `bson:"event_id"`
	AggregateID string        `bson:"aggregate_id"`
	EventType   string        `bson:"event_type"`
	RoutingKey  string        `bson:"routing_key"`
	Version     string        `bson:"version"`
	Payload     bson.Raw      `bson:"payload"`
	OccurredAt  time.Time     `bson:"occurred_at"`
	PublishedAt *time.Time    `bson:"published_at,omitempty"`
	Attempts    int           `bson:"attempts"`
	LastError   string        `bson:"last_error,omitempty"`
}

type OutboxRepository struct {
	col *mongo.Collection
}

func NewOutboxRepository(db *mongo.Database) *OutboxRepository {
	return &OutboxRepository{col: db.Collection(outboxCollection)}
}

// EnsureOutboxIndexes creates the indexes used by the publisher hot path.
// Safe to call repeatedly; Mongo is idempotent on identical index specs.
func EnsureOutboxIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection(outboxCollection)
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "event_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_event_id"),
		},
		{
			// Partial index — only unpublished rows. Keeps the index small and the
			// pending-events scan cheap even after millions of historical rows.
			Keys: bson.D{{Key: "occurred_at", Value: 1}},
			Options: options.Index().
				SetName("pending_by_occurred_at").
				SetPartialFilterExpression(bson.M{"published_at": bson.M{"$exists": false}}),
		},
	})
	return err
}

func (r *OutboxRepository) Insert(ctx context.Context, e *OutboxEvent) error {
	_, err := r.col.InsertOne(ctx, e)
	if err != nil {
		return fmt.Errorf("outbox insert: %w", err)
	}
	return nil
}

// FindPending returns up to limit unpublished events in occurrence order.
// No locking — works for single-instance deploys; the relay should be
// singleton in production. (Distributed locking is a Phase 7 concern.)
func (r *OutboxRepository) FindPending(ctx context.Context, limit int) ([]OutboxEvent, error) {
	cur, err := r.col.Find(ctx,
		bson.M{"published_at": bson.M{"$exists": false}},
		options.Find().SetSort(bson.D{{Key: "occurred_at", Value: 1}}).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, fmt.Errorf("outbox find: %w", err)
	}
	defer cur.Close(ctx)

	var out []OutboxEvent
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("outbox decode: %w", err)
	}
	return out, nil
}

func (r *OutboxRepository) MarkPublished(ctx context.Context, id bson.ObjectID, when time.Time) error {
	_, err := r.col.UpdateByID(ctx, id, bson.M{
		"$set":   bson.M{"published_at": when, "last_error": ""},
		"$inc":   bson.M{"attempts": 1},
	})
	return err
}

func (r *OutboxRepository) MarkFailed(ctx context.Context, id bson.ObjectID, errMsg string) error {
	if len(errMsg) > 1000 {
		errMsg = errMsg[:1000]
	}
	_, err := r.col.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{"last_error": errMsg},
		"$inc": bson.M{"attempts": 1},
	})
	return err
}

// MustMarshalPayload converts any JSON-serializable payload to bson.Raw.
// Used by the OutboxPublisher when building rows.
func MustMarshalPayload(v any) (bson.Raw, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal payload to json: %w", err)
	}
	var doc bson.M
	if err := json.Unmarshal(jsonBytes, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal json into bson.M: %w", err)
	}
	raw, err := bson.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("marshal bson: %w", err)
	}
	return raw, nil
}
