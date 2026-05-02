package user

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("user not found")

type Repository interface {
	Upsert(ctx context.Context, u User) error
	GetEmail(ctx context.Context, userID string) (string, error)
}
