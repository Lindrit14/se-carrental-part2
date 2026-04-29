package token

import (
	"context"
	"time"
)

// PasswordReset stores the hash of a one-time password-reset token.
type PasswordReset struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
}

func (p *PasswordReset) IsActive(now time.Time) bool {
	return p.UsedAt == nil && now.Before(p.ExpiresAt)
}

type PasswordResetRepository interface {
	Create(ctx context.Context, r *PasswordReset) error
	FindByHash(ctx context.Context, tokenHash string) (*PasswordReset, error)
	MarkUsed(ctx context.Context, id string, usedAt time.Time) error
}
