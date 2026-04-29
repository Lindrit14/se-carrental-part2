package token

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("refresh token not found")
	ErrRevoked  = errors.New("refresh token revoked")
	ErrExpired  = errors.New("refresh token expired")
)

type Repository interface {
	Create(ctx context.Context, t *RefreshToken) error
	FindByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Update(ctx context.Context, t *RefreshToken) error
	RevokeAllForUser(ctx context.Context, userID string) error
}
