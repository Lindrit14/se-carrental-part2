// Package mongodb contains the Mongo driver bootstrap and repository
// implementations. Domain entities are mapped to/from BSON documents
// here — no BSON tags leak into the domain package.
package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	CollectionUsers          = "users"
	CollectionRefreshTokens  = "refresh_tokens"
	CollectionPasswordResets = "password_resets"
)

// Connect establishes a connection and pings to verify reachability.
func Connect(ctx context.Context, uri string, connectTimeout time.Duration) (*mongo.Client, error) {
	cctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := client.Ping(cctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo ping: %w", err)
	}
	return client, nil
}

// EnsureIndexes creates idempotent indexes on first boot.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	users := db.Collection(CollectionUsers)
	if _, err := users.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_email"),
		},
	}); err != nil {
		return fmt.Errorf("users indexes: %w", err)
	}

	rt := db.Collection(CollectionRefreshTokens)
	if _, err := rt.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "token_hash", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_token_hash")},
		{Keys: bson.D{{Key: "user_id", Value: 1}}, Options: options.Index().SetName("user_id")},
		{Keys: bson.D{{Key: "expires_at", Value: 1}}, Options: options.Index().SetName("ttl_expires_at").SetExpireAfterSeconds(0)},
	}); err != nil {
		return fmt.Errorf("refresh_tokens indexes: %w", err)
	}

	pr := db.Collection(CollectionPasswordResets)
	if _, err := pr.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "token_hash", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_reset_hash")},
		{Keys: bson.D{{Key: "expires_at", Value: 1}}, Options: options.Index().SetName("ttl_reset_expires_at").SetExpireAfterSeconds(0)},
	}); err != nil {
		return fmt.Errorf("password_resets indexes: %w", err)
	}
	return nil
}
