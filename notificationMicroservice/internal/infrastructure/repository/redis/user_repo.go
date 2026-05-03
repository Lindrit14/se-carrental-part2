package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lindritprekaj/notification-service/internal/domain/user"
)

// UserRepo is a Redis-backed implementation of user.Repository. It stores
// the user_id→email read model populated from user.registered events and
// read by the booking handlers, which receive only customerExternalUserId
// in their event payloads.
type UserRepo struct {
	client *redis.Client
	prefix string
}

// New constructs a UserRepo and pings the server to fail fast on a
// misconfigured address — matching the eager-failure semantics the SQLite
// implementation had via its embedded migration on Open.
func New(addr, password string, db int, keyPrefix string) (*UserRepo, error) {
	if addr == "" {
		return nil, errors.New("redis: empty addr")
	}
	if keyPrefix == "" {
		keyPrefix = "notif"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis: ping %s: %w", addr, err)
	}

	return &UserRepo{client: client, prefix: keyPrefix}, nil
}

func (r *UserRepo) Close() error { return r.client.Close() }

func (r *UserRepo) Upsert(ctx context.Context, u user.User) error {
	if u.ID == "" || u.Email == "" {
		return errors.New("upsert: empty id or email")
	}
	if err := r.client.Set(ctx, r.userKey(u.ID), u.Email, 0).Err(); err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetEmail(ctx context.Context, userID string) (string, error) {
	email, err := r.client.Get(ctx, r.userKey(userID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", user.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get email: %w", err)
	}
	return email, nil
}

func (r *UserRepo) userKey(userID string) string {
	return r.prefix + ":user:" + userID
}
